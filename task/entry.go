package task

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	noted "noted/config"
	"noted/logging"
	"os"
	"path"
	"time"
)

type EntryFile struct {
	Entries []Entry
}

type Entry struct {
	CreatedAt    time.Time
	DueAt        *time.Time `yaml:"due_at,omitempty"`
	ScheduledFor *time.Time `yaml:"scheduled_for,omitempty"`
	Task         string
	Status       Status
}

type Task struct {
	File         string
	CreatedAt    time.Time
	DueAt        *time.Time
	ScheduledFor *time.Time
	Task         string
	Status       Status
}

type Status byte

const (
	ToDo = iota
	Scheduled
	Paused
	Cancelled
	Done
)

func CreateTask(task string, due *time.Time) error {
	createdTime := time.Now()
	prefix := viper.GetString(noted.ConfigTaskPrefix)
	storageDir := viper.GetString(noted.ConfigStorageDir)
	taskPath := path.Join(storageDir, prefix)
	taskFilePath := path.Join(storageDir, prefix, fmt.Sprintf("%d-%s.yaml", createdTime.Year(), createdTime.Month()))
	var taskFile *os.File
	var err error
	var taskEntries EntryFile

	if _, err := os.Stat(taskPath); errors.Is(err, os.ErrNotExist) {
		if err = os.MkdirAll(taskPath, 0755); err != nil {
			logging.Logger.Error("failed to create journal path", zap.Error(err))
			return err
		}
		taskFile, err = os.Create(taskFilePath)
		taskEntries = EntryFile{
			Entries: make([]Entry, 0),
		}
	} else {
		taskFile, err = os.OpenFile(taskFilePath, os.O_WRONLY, os.ModeAppend)
		data, err := os.ReadFile(taskFilePath)
		if err != nil {
			logging.Logger.Error("failed to read existing task file", zap.Error(err), zap.String("file", taskFilePath))
			return err
		}
		if err = yaml.Unmarshal(data, &taskEntries); err != nil {
			logging.Logger.Error("failed to unmarshal existing task file", zap.Error(err), zap.String("file", taskFilePath))
			return err
		}
	}

	if err != nil {
		logging.Logger.Error("failed to create/open task file", zap.Error(err), zap.String("file", taskFilePath))
		return err
	}

	defer taskFile.Close()

	taskEntries.Entries = append(taskEntries.Entries, Entry{
		CreatedAt:    createdTime,
		DueAt:        due,
		ScheduledFor: nil,
		Task:         task,
		Status:       ToDo,
	})

	if output, err := yaml.Marshal(taskEntries); err != nil {
		logging.Logger.Error("failed to marshal task file data", zap.Error(err))
		return err
	} else {
		_, err = taskFile.Write(output)
	}

	return err
}

func ListTasks() []Task {
	return nil
}
