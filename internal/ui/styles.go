package ui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("24")).
			Padding(0, 1)

	bodyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	rowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	selectedRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("230")).
				Background(lipgloss.Color("236"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("204")).
			Bold(true)

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("248")).
			Italic(true)
)
