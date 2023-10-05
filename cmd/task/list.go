package task

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	config "noted/config"
	"noted/logging"
	"noted/task"
)

type ListModel struct {
	list list.Model
}

func newListModel(tasks []task.Task) ListModel {
	items := make([]list.Item, 0)

	for _, t := range tasks {
		items = append(items, t)
	}

	taskList := list.New(items, newTaskItemDelegate(), 0, 0)
	taskList.Title = config.TitleStyle.Render("Tasks")
	taskList.Styles.Title = config.TitleStyle

	return ListModel{
		list: taskList,
	}
}

func newTaskItemDelegate() list.DefaultDelegate {
	selectionKeyBinding := key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "change status"),
	)

	rotateStatusKeyBinding := key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "iterate status"),
	)

	cancelKeyBinding := key.NewBinding(
		key.WithKeys("x", "backspace"),
		key.WithHelp("x", "cancel"),
	)
	delegate := list.NewDefaultDelegate()

	delegate.UpdateFunc = func(msg tea.Msg, model *list.Model) tea.Cmd {
		if taskItem, ok := model.SelectedItem().(task.Task); ok {
			title := taskItem.Title()
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch {
				case key.Matches(msg, selectionKeyBinding):
					return model.NewStatusMessage(title)
				case key.Matches(msg, rotateStatusKeyBinding):
					if taskItem.Status == task.Done {
						taskItem.Status = task.ToDo
					} else {
						taskItem.Status++
					}
					if err := task.UpdateTask(taskItem); err != nil {
						logging.Logger.Error("failed to update task item", zap.Error(err))
						return nil
					} else {
						model.Items()[model.Index()] = taskItem
						return model.NewStatusMessage(fmt.Sprintf("task updated to %s", taskItem.Status.AsString()))
					}
				case key.Matches(msg, cancelKeyBinding):
					taskItem.Status = task.Cancelled
					if err := task.UpdateTask(taskItem); err != nil {
						logging.Logger.Error("failed to update task item", zap.Error(err))
						tea.Quit()
					} else {
						model.Items()[model.Index()] = taskItem
						return model.NewStatusMessage("task cancelled")
					}
				}
			}
			return nil
		} else {
			return nil
		}
	}

	help := []key.Binding{selectionKeyBinding, rotateStatusKeyBinding, cancelKeyBinding}
	delegate.ShortHelpFunc = func() []key.Binding {
		return help
	}
	delegate.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return delegate
}

var ListTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "list and alter todo list items",
	Long:  "manage your task list",
	Run: func(cmd *cobra.Command, args []string) {
		tasks := task.ListTasks(false)
		program := tea.NewProgram(newListModel(tasks))
		if _, err := program.Run(); err != nil {
			logging.Logger.Fatal("failed to execute program", zap.Error(err))
		}
	},
}

func (l ListModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (l ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := config.DocStyle.GetFrameSize()
		l.list.SetSize(msg.Width-h, msg.Height-v)
	}

	newModel, cmd := l.list.Update(msg)
	l.list = newModel
	commands = append(commands, cmd)

	return l, tea.Batch(commands...)
}

func (l ListModel) View() string {
	return config.DocStyle.Render(l.list.View())
}
