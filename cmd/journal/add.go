package journal

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"noted/journal"
	"noted/logging"
	"time"
)

var AddToJournalCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new journal command",
	Long:  "Create a new journal entry",
	Run: func(cmd *cobra.Command, args []string) {
		model := newEntry()
		program := tea.NewProgram(model)
		if _, err := program.Run(); err != nil {
			logging.Logger.Fatal("program failure", zap.Error(err))
		}
	},
}

type NewEntry struct {
	textInput textinput.Model
	err       error
}

type errMsg error

func newEntry() NewEntry {
	var ti = textinput.New()
	ti.Placeholder = "a new thing that happened"
	ti.Focus()
	return NewEntry{
		textInput: ti,
		err:       nil,
	}
}

func (n NewEntry) Init() tea.Cmd {
	return textinput.Blink
}

func (n NewEntry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			if err := journal.SaveJournalEntry(time.Now(), n.textInput.Value()); err != nil {
				logging.Logger.Error("failed to save entry", zap.Error(err))
			}
			return n, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		n.err = msg
		return n, nil
	}

	n.textInput, cmd = n.textInput.Update(msg)
	return n, cmd
}

func (n NewEntry) View() string {
	return n.textInput.View()
}
