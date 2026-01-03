package ui

import (
	"dbear/internal/config"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Bold(true).
			MarginBottom(1)
	nameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)
	typeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99"))
	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))
	itemStyle = lipgloss.NewStyle().
			MarginLeft(2).
			MarginBottom(1)
)

func DisplayConnectionsList(connections []config.Connection) error {
	if len(connections) == 0 {
		fmt.Println("No connections found.")
		return nil
	}

	fmt.Println(headerStyle.Render("Database Connections"))
	fmt.Println()

	for _, conn := range connections {
		portStr := ""
		if conn.Port > 0 {
			portStr = fmt.Sprintf(":%d", conn.Port)
		}
		hostDisplay := conn.Host + portStr

		output := itemStyle.Render(
			nameStyle.Render(conn.Name) + " " +
				typeStyle.Render("("+conn.Type+")") + "\n" +
				infoStyle.Render("  Host: "+hostDisplay+" | Database: "+conn.Database),
		)
		fmt.Println(output)
	}

	return nil
}
