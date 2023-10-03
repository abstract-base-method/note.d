package journal

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"noted/config"
	"noted/logging"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	Year    int
	Month   string
	Day     int
	Message string
}

func (e Entry) Title() string {
	return e.Message
}

func (e Entry) Description() string {
	return fmt.Sprintf("%d/%s/%d", e.Year, e.Month, e.Day)
}

func (e Entry) FilterValue() string {
	return e.Message
}

func SaveJournalEntry(datetime time.Time, entry string) error {
	prefix := viper.GetString(noted.ConfigJournalPrefix)
	storageDir := viper.GetString(noted.ConfigStorageDir)
	journalPath := path.Join(storageDir, prefix)
	journalFilePath := path.Join(storageDir, prefix, fmt.Sprintf("%d-%s.md", datetime.Year(), datetime.Month()))
	var journalFile *os.File
	var err error

	if _, err := os.Stat(journalPath); errors.Is(err, os.ErrNotExist) {
		if err = os.MkdirAll(journalPath, 0755); err != nil {
			logging.Logger.Error("failed to create journal path", zap.Error(err))
			return err
		}
		journalFile, err = os.Create(journalFilePath)
	} else {
		journalFile, err = os.OpenFile(journalFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	}
	if err != nil {
		logging.Logger.Error("failed to create/open journal file", zap.Error(err))
		return err
	}

	defer journalFile.Close()

	_, _, day := datetime.Date()
	_, err = journalFile.WriteString(fmt.Sprintf("- %s %d: %s\n", datetime.Weekday(), day, entry))

	journalFile.Sync()

	if err != nil {
		logging.Logger.Error("failed to append to journal", zap.String("file", journalFilePath), zap.Error(err))
	}

	return err
}

func GetEntries(sortOldestAscending bool) []Entry {
	prefix := viper.GetString(noted.ConfigJournalPrefix)
	storageDir := viper.GetString(noted.ConfigStorageDir)
	journalPath := path.Join(storageDir, prefix)
	files, err := os.ReadDir(journalPath)
	if err != nil {
		logging.Logger.Fatal("failed to list journal files", zap.Error(err), zap.String("directory", journalPath))
	}

	entries := make([]Entry, 0)

	for _, file := range files {
		if fileHandle, err := os.Open(path.Join(journalPath, file.Name())); err != nil {
			logging.Logger.Error("failed to read file", zap.String("file", path.Join(journalPath, file.Name())), zap.Error(err))
		} else {
			defer fileHandle.Close()
			fileElements := strings.Split(file.Name(), "-")
			year, _ := strconv.Atoi(fileElements[0])
			month := strings.Split(fileElements[1], ".")[0]
			scanner := bufio.NewScanner(fileHandle)
			for scanner.Scan() {
				text := scanner.Text()[3:]
				itemElements := strings.Split(text, ":")
				dateElements := strings.Split(itemElements[0], " ")
				day, _ := strconv.Atoi(dateElements[1])
				message := strings.Join(itemElements[1:], ":")
				entries = append(entries, Entry{
					Year:    year,
					Month:   month,
					Day:     day,
					Message: message,
				})
			}
		}
	}

	if sortOldestAscending {
		slices.Reverse(entries)
	}

	return entries
}
