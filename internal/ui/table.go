package ui

import (
	"example.com/soundz/internal/player"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func GetTable() table.Model {
	columns := []table.Column{
		{Title: " ", Width: 4},
		{Title: "Title", Width: 30},
		{Title: "Artist", Width: 20},
		{Title: "File Name", Width: 30},
		{Title: "Length", Width: 6},
		{Title: "Format", Width: 6},
	}

	rows := player.GetSongsInDir()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	return t
}
