package task

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"noted/task"
	"strings"
)

type errMsg error

type newTaskModel struct {
	inputs     []textinput.Model
	focusIndex int
	cursorMode cursor.Mode
	err        error
}

const (
	taskInputId = iota
	dueDateInputId
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

var AddTaskCmd = &cobra.Command{
	Use:   "add",
	Short: "add task",
	Long:  "create a new task",
	Run: func(cmd *cobra.Command, args []string) {
		task := createNewTaskModel()
		program := tea.NewProgram(task)
		if _, err := program.Run(); err != nil {
			log.Fatal("failed to run program", zap.Error(err))
		}
	},
}

func createNewTaskModel() newTaskModel {
	model := newTaskModel{
		inputs:     make([]textinput.Model, 2),
		focusIndex: 0,
		cursorMode: cursor.CursorBlink,
		err:        nil,
	}

	var t textinput.Model
	for i := range model.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle

		switch i {
		case taskInputId:
			t.Placeholder = "task description"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case dueDateInputId:
			t.Placeholder = "due date"
			t.TextStyle = noStyle
			t.PromptStyle = noStyle
		}

		model.inputs[i] = t
	}

	return model
}

func (n newTaskModel) Init() tea.Cmd {
	return textinput.Blink
}

func (n newTaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return n, tea.Quit
		case tea.KeyCtrlR:
			n.cursorMode++

			if n.cursorMode > cursor.CursorHide {
				n.cursorMode = cursor.CursorBlink
			}

			commands := make([]tea.Cmd, len(n.inputs))
			for i := range n.inputs {
				commands[i] = n.inputs[i].Cursor.SetMode(n.cursorMode)
			}
			return n, tea.Batch(commands...)
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			command := msg.String()

			if command == "enter" && n.focusIndex == len(n.inputs) {
				// todo: handle time
				n.err = task.CreateTask(n.inputs[taskInputId].Value(), nil)
				// todo: process form data
				return n, tea.Quit
			}

			if command == "up" || command == "shift+tab" {
				n.focusIndex--
			} else {
				n.focusIndex++
			}

			if n.focusIndex > len(n.inputs) {
				n.focusIndex = 0
			} else if n.focusIndex < 0 {
				n.focusIndex = len(n.inputs)
			}

			commands := make([]tea.Cmd, len(n.inputs))
			for i := 0; i < len(n.inputs); i++ {
				if i == n.focusIndex {
					commands[i] = n.inputs[i].Focus()
					n.inputs[i].PromptStyle = focusedStyle
					n.inputs[i].TextStyle = focusedStyle
					continue
				}
				n.inputs[i].Blur()
				n.inputs[i].PromptStyle = noStyle
				n.inputs[i].TextStyle = noStyle
			}
			return n, tea.Batch(commands...)
		}
	}

	cmd := n.updateInputs(msg)

	return n, cmd
}

func (n newTaskModel) updateInputs(msg tea.Msg) tea.Cmd {
	commands := make([]tea.Cmd, len(n.inputs))

	for i := range n.inputs {
		n.inputs[i], commands[i] = n.inputs[i].Update(msg)
	}

	return tea.Batch(commands...)
}

func (n newTaskModel) View() string {
	var builder strings.Builder

	for i := range n.inputs {
		builder.WriteString(n.inputs[i].View())
		if i < len(n.inputs)-1 {
			builder.WriteRune('\n')
		}
	}

	button := blurredButton
	if n.focusIndex == len(n.inputs) {
		button = focusedButton
	}
	fmt.Fprintf(&builder, "\n\n%s\n", button)

	builder.WriteString(helpStyle.Render("cursor mode is "))
	builder.WriteString(cursorModeHelpStyle.Render(n.cursorMode.String()))
	builder.WriteString(helpStyle.Render(" -- ctrl+r to change style"))

	return builder.String()
}
