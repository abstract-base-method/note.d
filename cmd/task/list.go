package task

import "github.com/spf13/cobra"

var ListTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "list and alter todo list items",
	Long:  "manage your task list",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
