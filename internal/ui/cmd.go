package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type tickNowMsg time.Time

func TickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func TickNow() tea.Cmd {
	return func() tea.Msg {
		return tickNowMsg(time.Now())
	}
}
