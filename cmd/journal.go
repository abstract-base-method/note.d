package cmd

import (
	"github.com/spf13/cobra"
	"noted/cmd/journal"
)

func init() {
	JournalCmd.AddCommand(journal.AddToJournalCmd)
	JournalCmd.AddCommand(journal.ListJournalCmd)
}

var JournalCmd = &cobra.Command{
	Use:   "journal",
	Short: "journaling functions",
	Long:  "Maintain your journal",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
