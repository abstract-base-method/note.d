package main

import (
	"noted/cmd"
	"noted/logging"
	"os"
)

func main() {
	logging.Logger.Debug("logger construction succeeded")
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(127)
	}
}
