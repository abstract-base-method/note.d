package journal

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	config "noted/config"
	"noted/journal"
	"noted/logging"
)

type entryList struct {
	list list.Model
}

var ListJournalCmd = &cobra.Command{
	Use:   "list",
	Short: "list recent journal entries",
	Long:  "list recent journal entries",
	Run: func(cmd *cobra.Command, args []string) {
		// first we need to read all entries
		entries := journal.GetEntries(true)
		program := tea.NewProgram(newEntryList(entries), tea.WithAltScreen())
		if _, err := program.Run(); err != nil {
			logging.Logger.Fatal("program failure", zap.Error(err))
		}
	},
}

func newEntryList(entries []journal.Entry) entryList {
	items := make([]list.Item, 0)
	for _, journalEntry := range entries {
		items = append(items, journalEntry)
		//items = append(items, entry(fmt.Sprintf("%d/%s/%d: %s", journalEntry.Year, journalEntry.Month, journalEntry.Day, journalEntry.Message)))
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = config.TitleStyle.Render("Recent Journal Entries")
	l.Styles.Title = config.TitleStyle
	//l.Styles.PaginationStyle = paginationStyle
	//l.Styles.HelpStyle = helpStyle

	return entryList{
		list: l,
	}
}

func (e entryList) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (e entryList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := config.DocStyle.GetFrameSize()
		e.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return e, tea.Quit

		case tea.KeyEnter:
			_ = e.list.SelectedItem().(journal.Entry)
			return e, tea.Quit
		}
	}

	var cmd tea.Cmd
	e.list, cmd = e.list.Update(msg)
	return e, cmd
}

func (e entryList) View() string {
	return config.DocStyle.Render(e.list.View())
}
