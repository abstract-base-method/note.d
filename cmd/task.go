package cmd

import (
	"github.com/spf13/cobra"
	"noted/cmd/task"
)

func init() {
	TaskCmd.AddCommand(task.AddTaskCmd)
	TaskCmd.AddCommand(task.ListTasksCmd)
}

var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "interact with tasks",
	Long:  "Use tasks to remind yourself to do a thing",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
