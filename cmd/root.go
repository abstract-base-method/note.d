package cmd

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"noted/config"
	"noted/logging"
	"os"
	"path"
)

func init() {
	cobra.OnInitialize(initConfiguration)
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "configuration file (default is $HOME/.noted.yaml")
	RootCmd.AddCommand(JournalCmd)
	RootCmd.AddCommand(TaskCmd)
}

var RootCmd = &cobra.Command{
	Use:   "noted",
	Short: "a note taking tool",
	Long:  "Note.d is a note taking tool for kool kids 😎",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var configFile string

func initConfiguration() {
	var home, homeErr = homedir.Dir()
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		if homeErr != nil {
			log.Fatal(homeErr)
		}
		viper.AddConfigPath(home)
		viper.SetConfigFile(".noted")
	}

	// config defaults
	viper.SetDefault(noted.ConfigStorageDir, path.Join(home, ".noted"))
	viper.SetDefault(noted.ConfigJournalPrefix, "journal")
	viper.SetDefault(noted.ConfigTaskPrefix, "task")

	if err := viper.ReadInConfig(); err != nil {
		logging.Logger.Debug("cannot find config file")
	}

	// ensure filesystem is as we expect
	storagePath := viper.GetString(noted.ConfigStorageDir)
	journalPrefix := viper.GetString(noted.ConfigJournalPrefix)
	taskPrefix := viper.GetString(noted.ConfigTaskPrefix)
	journalPath := path.Join(storagePath, journalPrefix)
	taskPath := path.Join(storagePath, taskPrefix)

	directories := []string{storagePath, journalPath, taskPath}

	for _, dir := range directories {
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			if err = os.MkdirAll(storagePath, 0755); err != nil {
				logging.Logger.Fatal("failed to initialize directory", zap.String("directory", dir), zap.Error(err))
			}
		}
	}
}
