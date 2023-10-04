package task

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	config "noted/config"
	"noted/logging"
	"noted/task"
	"strings"
)

type ListModel struct {
	entries      []task.Task
	selected     *task.Task
	taskView     []textinput.Model
	viewMode     viewMode
	taskListView list.Model
}

func (l ListModel) enterDetailMode(taskItem task.Task) {
	l.viewMode = viewModeDetail
	l.selected = &taskItem
	for i := range l.taskView {
		switch i {
		case detailViewTask:
			l.taskView[i].SetValue(taskItem.Task)
		case detailViewDescription:
			l.taskView[i].SetValue(taskItem.Detail)
		case detailViewStatus:
			l.taskView[i].SetValue(task.StatusAsString(taskItem.Status))
		case detailViewDueAt:
			l.taskView[i].SetValue("TODO FIX ME")
		case detailViewScheduledFor:
			l.taskView[i].SetValue("TODO FIX ME")
		}
	}
}

func (l ListModel) enterListMode() {
	l.viewMode = viewModeList
}

type viewMode int

const (
	viewModeList = iota
	viewModeDetail
)

const (
	detailViewTask = iota
	detailViewDescription
	detailViewStatus
	detailViewScheduledFor
	detailViewDueAt
)

func newListModel(tasks []task.Task) ListModel {
	items := make([]list.Item, 0)
	for _, taskItem := range tasks {
		items = append(items, taskItem)
	}

	fields := make([]textinput.Model, 5)
	for i := range fields {
		input := textinput.New()
		switch i {
		case detailViewTask:
			input.Placeholder = "task description"
			input.Focus()
			input.PromptStyle = focusedStyle
			input.TextStyle = focusedStyle
		case detailViewDescription:
			input.Placeholder = "task detail"
			input.PromptStyle = noStyle
			input.TextStyle = noStyle
		case detailViewStatus:
			input.Placeholder = "task status"
			input.PromptStyle = noStyle
			input.TextStyle = noStyle
		case detailViewScheduledFor:
			input.Placeholder = "scheduled for"
			input.PromptStyle = noStyle
			input.TextStyle = noStyle
		case detailViewDueAt:
			input.Placeholder = "due at"
			input.PromptStyle = noStyle
			input.TextStyle = noStyle
		}
		fields[i] = input
	}

	taskListModel := list.New(items, list.NewDefaultDelegate(), 0, 0)
	taskListModel.Title = config.TitleStyle.Render("Task List")
	taskListModel.Styles.Title = config.TitleStyle
	taskListModel.Styles.PaginationStyle = config.PaginationStyle
	taskListModel.Styles.HelpStyle = helpStyle

	return ListModel{
		entries:      tasks,
		selected:     nil,
		taskView:     fields,
		viewMode:     viewModeList,
		taskListView: taskListModel,
	}
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := config.DocStyle.GetFrameSize()
		l.taskListView.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return l, tea.Quit
		case tea.KeyEsc:
			if l.viewMode == viewModeList {
				return l, tea.Quit
			} else {
				l.enterListMode()
				return l, tea.Quit
			}
		case tea.KeyEnter:
			commands := make([]tea.Cmd, 0)
			if l.viewMode == viewModeList {
				if selected, err := l.taskListView.SelectedItem().(task.Task); err != true {
					logging.Logger.Fatal("failed to select task")
				} else {
					l.enterDetailMode(selected)
					for i := range l.taskView {
						commands = append(commands, l.taskView[i].Cursor.SetMode(cursor.CursorBlink))
					}
				}
			} else {
				// todo: handle updates in the detail view
				panic("yikes")
			}
			_, cmd := l.taskListView.Update(msg)
			commands = append(commands, cmd)
			return l, tea.ClearScreen
		}
	}

	var cmd tea.Cmd
	l.taskListView, cmd = l.taskListView.Update(msg)
	return l, cmd
}

func (l ListModel) View() string {
	switch l.viewMode {
	case viewModeList:
		return config.DocStyle.Render(l.taskListView.View())
	case viewModeDetail:
		// todo: render the detail view with the ability to change
		var builder strings.Builder

		for i := range l.taskView {
			builder.WriteString(l.taskView[i].View())

			if i < len(l.taskView)-1 {
				builder.WriteRune('\n')
			}
		}

		return builder.String()
	default:
		panic("whoops")
	}
}
