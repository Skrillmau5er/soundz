package ui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep/v2"
)

type model struct {
	table            table.Model
	progress         progress.Model
	ctrl             *beep.Ctrl
	streamer         beep.StreamSeekCloser
	format           beep.Format
	songLength       string
	songPos          string
	songSampleRate   beep.SampleRate
	currentSongIndex int
}

func (m model) Init() tea.Cmd { return nil }
