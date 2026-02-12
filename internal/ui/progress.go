package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type progressTask func() (interface{}, error)

type progressDoneMsg struct {
	value interface{}
	err   error
}

type progressModel struct {
	spinner  spinner.Model
	message  string
	task     progressTask
	value    interface{}
	err      error
	quitting bool
}

func newProgressModel(message string, task progressTask) progressModel {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return progressModel{
		spinner: spin,
		message: message,
		task:    task,
	}
}

func (m progressModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, runTaskCmd(m.task))
}

func runTaskCmd(task progressTask) tea.Cmd {
	return func() tea.Msg {
		value, err := task()
		return progressDoneMsg{value: value, err: err}
	}
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case progressDoneMsg:
		m.value = msg.value
		m.err = msg.err
		m.quitting = true
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m progressModel) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("\n  %s %s\n\n", m.spinner.View(), m.message)
}

func RunWithSpinner(message string, task progressTask) (interface{}, error) {
	m := newProgressModel(message, task)
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return nil, err
	}
	model := final.(progressModel)
	return model.value, model.err
}
