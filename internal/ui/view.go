package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) renderArt() string {
	return `
    ███████╗ ██████╗ ██╗   ██╗███╗   ██╗██████╗ ███████╗
    ██╔════╝██╔═══██╗██║   ██║████╗  ██║██╔══██╗    ███║
    ███████╗██║   ██║██║   ██║██╔██╗ ██║██║  ██║  ███╔═╝
    ╚════██║██║   ██║██║   ██║██║╚██╗██║██║  ██║███╔═╝
    ███████║╚██████╔╝╚██████╔╝██║ ╚████║██████╔╝███████╗
    ╚══════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═══╝╚═════╝ ╚══════╝
`
}

func (m model) renderTable() string {
	return baseStyle.Render(m.table.View())
}

func (m model) renderBottom() string {
	songPos := m.songPos
	if songPos == "" {
		songPos = "    "
	}
	songLength := m.songLength
	if songLength == "" {
		songLength = "    "
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, songPos, m.progress.View(), songLength)
}

func (m model) renderHelp() string {
	style := lipgloss.NewStyle().Padding(2)

	helpView := m.help.View(m.keys)
	return style.Render(helpView)
}

func (m model) View() string {

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, m.renderArt(), m.renderHelp()),
		m.renderTable(),
		m.renderBottom(),
	)
}
