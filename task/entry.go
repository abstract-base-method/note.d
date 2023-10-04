package task

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	config "noted/config"
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
	Detail       string
	Status       Status
}

type Task struct {
	File         string
	CreatedAt    time.Time
	DueAt        *time.Time
	ScheduledFor *time.Time
	Task         string
	Detail       string
	Status       Status
}

func (t Task) Title() string {
	return t.Task
}

func (t Task) Description() string {
	return t.Detail
}

func (t Task) FilterValue() string {
	return fmt.Sprintf("%s %s %d", t.Task, t.Detail, t.Status)
}

func (t Entry) toTask(file string) Task {
	return Task{
		File:         file,
		CreatedAt:    t.CreatedAt,
		DueAt:        t.DueAt,
		ScheduledFor: t.ScheduledFor,
		Task:         t.Task,
		Detail:       t.Detail,
		Status:       t.Status,
	}
}

type Status byte

const (
	ToDo = iota
	Scheduled
	Paused
	Cancelled
	Done
)

func StatusAsString(status Status) string {
	switch status {
	case ToDo:
		return "TODO"
	case Scheduled:
		return "SCHEDULED"
	case Paused:
		return "PAUSE"
	case Cancelled:
		return "CANCELLED"
	case Done:
		return "DONE"
	default:
		return "UNK"
	}
}

func CreateTask(task string, detail string, due *time.Time) error {
	createdTime := time.Now()
	prefix := viper.GetString(config.ConfigTaskPrefix)
	storageDir := viper.GetString(config.ConfigStorageDir)
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
		Detail:       detail,
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

func ListTasks(includeCompleted bool) []Task {
	storageDir := viper.GetString(config.ConfigStorageDir)
	prefixDir := viper.GetString(config.ConfigTaskPrefix)
	taskPath := path.Join(storageDir, prefixDir)
	files, err := os.ReadDir(taskPath)

	if err != nil {
		logging.Logger.Fatal("failed to read from task directory", zap.Error(err), zap.String("directory", taskPath))
	}

	tasks := make([]Task, 0)

	for _, file := range files {
		if data, err := os.ReadFile(path.Join(taskPath, file.Name())); err != nil {
			logging.Logger.Error("failed to read task file", zap.String("file", file.Name()), zap.Error(err))
		} else {
			var contents EntryFile
			if err = yaml.Unmarshal(data, &contents); err != nil {
				logging.Logger.Error("failed to unmarshal yaml for task file", zap.Error(err), zap.String("file", file.Name()))
			} else {
				for _, entry := range contents.Entries {
					if entry.Status == Done && includeCompleted {
						tasks = append(tasks, entry.toTask(path.Join(taskPath, file.Name())))
					} else {
						tasks = append(tasks, entry.toTask(path.Join(taskPath, file.Name())))
					}
				}
			}
		}
	}

	return tasks
}
