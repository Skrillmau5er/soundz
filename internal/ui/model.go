package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep/v2"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	PlaySong  key.Binding
	PlayPause key.Binding
	NextSong  key.Binding
	PrevSong  key.Binding
	Help      key.Binding
	Quit      key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.PlaySong, k.PlayPause, k.NextSong, k.PrevSong},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "rewind 5 seconds"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "skip forward 5 seconds"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	PlaySong: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "play song"),
	),
	PlayPause: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "toggle pause/play"),
	),
	NextSong: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next song"),
	),
	PrevSong: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "previous song"),
	),
}

func (m model) Init() tea.Cmd {
	return nil
}

type model struct {
	keys             keyMap
	help             help.Model
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
