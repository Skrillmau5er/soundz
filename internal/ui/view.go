package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) renderArt() string {
	art := `
███████╗ ██████╗ ██╗   ██╗███╗   ██╗██████╗ ███████╗
██╔════╝██╔═══██╗██║   ██║████╗  ██║██╔══██╗    ███║
███████╗██║   ██║██║   ██║██╔██╗ ██║██║  ██║  ███╔═╝
╚════██║██║   ██║██║   ██║██║╚██╗██║██║  ██║███╔═╝
███████║╚██████╔╝╚██████╔╝██║ ╚████║██████╔╝███████╗
╚══════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═══╝╚═════╝ ╚══════╝`
	artStyle := lipgloss.NewStyle().Width(172).Align(lipgloss.Center)
	return artStyle.Render(art)
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

func (m model) View() string {
	blockB := "...\n...\n..."
	blockA := "...\n...\n...\n...\n..."

	// Join on the top edge
	str := lipgloss.JoinHorizontal(lipgloss.Top, blockA, blockB)
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.renderArt(),
		m.renderTable(),
		m.renderBottom(),
		str,
	)
}
