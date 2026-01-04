package ui

import (
	"dbear/internal/config"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			MarginBottom(1)
)

type connectionItem struct {
	conn config.Connection
}

func (i connectionItem) FilterValue() string {
	return i.conn.Name
}

func (i connectionItem) Title() string {
	portStr := ""
	if i.conn.Port > 0 {
		portStr = fmt.Sprintf(":%d", i.conn.Port)
	}
	return fmt.Sprintf("%s (%s) - %s%s/%s", i.conn.Name, i.conn.Type, i.conn.Host, portStr, i.conn.Database)
}

func (i connectionItem) Description() string {
	return fmt.Sprintf("%s on %s", i.conn.Type, i.conn.Host)
}

type selectorModel struct {
	list     list.Model
	choice   string
	quitting bool
	title    string
}

func (m selectorModel) Init() tea.Cmd {
	return nil
}

func (m selectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			selectedItem := m.list.SelectedItem()
			if selectedItem != nil {
				m.choice = selectedItem.(connectionItem).conn.Name
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m selectorModel) View() string {
	if m.quitting {
		return ""
	}
	displayTitle := "Select a connection"
	if m.title != "" {
		displayTitle = m.title
	}
	return "\n" + titleStyle.Render(displayTitle) + "\n" + m.list.View()
}

func SelectConnection(connections []config.Connection) (string, error) {
	return SelectConnectionWithTitle(connections, "Select a connection")
}

func SelectConnectionWithTitle(connections []config.Connection, title string) (string, error) {
	if len(connections) == 0 {
		return "", fmt.Errorf("no connections available")
	}

	items := make([]list.Item, len(connections))
	for i, conn := range connections {
		items[i] = connectionItem{conn: conn}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = ""
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle

	m := selectorModel{list: l, title: title}
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	if finalModel.(selectorModel).quitting {
		return "", fmt.Errorf("selection cancelled")
	}

	return finalModel.(selectorModel).choice, nil
}

