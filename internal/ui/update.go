package ui

import (
	"fmt"
	"math"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/skrillmau5er/soundz/internal/player"
)

func (m model) Cleanup() {
	if m.streamer != nil {
		m.streamer.Close()
	}
}

func (m *model) markPlayingRow(index int) {
	rows := m.table.Rows()

	for i := range rows {
		if i == index {
			rows[i][0] = "🔊" //
		} else {
			rows[i][0] = ""
		}
	}

	m.table.SetRows(rows)
}

func (m model) startSong(index int) (model, tea.Cmd) {
	if m.ctrl != nil {
		m.ctrl.Streamer = nil
	}

	rows := m.table.Rows()
	if index < 0 || index >= len(rows) {
		index = 0
	}

	row := rows[index]
	m.currentSongIndex = index

	ctrl, streamer, format, songSampleRate := player.PlaySongWithVisualizer(m.currentDir+"/"+row[3], row[5], m.visualizer)
	m.ctrl = ctrl
	m.streamer = streamer
	m.format = format
	m.songSampleRate = songSampleRate
	m.markPlayingRow(index)
	return m, tea.Batch(TickCmd(), TickNow())
}

func (m model) tickAction(now bool) (model, tea.Cmd) {
	pos, length := player.GetPosAndLen(m.streamer, m.songSampleRate)

	m.songPos = fmt.Sprintf("%v:%02d", math.Floor(float64(pos)/60), pos%60)
	m.songLength = fmt.Sprintf("%v:%02d", math.Floor(float64(length)/60), length%60)
	cmd := m.progress.SetPercent(float64(pos) / float64(length))

	if m.streamer.Position() >= m.streamer.Len() {
		return m.startSong(m.currentSongIndex + 1)
	}

	if now {
		return m, cmd
	}
	return m, tea.Batch(TickCmd(), cmd)
}

func (m model) seek(cmd tea.Cmd, dir string) (model, tea.Cmd) {
	if m.ctrl == nil {
		return m, cmd
	}

	fiveSecs := int(m.songSampleRate) * 5

	if dir == "left" {
		fiveSecs = fiveSecs * -1
	}
	newPosition := m.streamer.Position() + fiveSecs

	speaker.Lock()
	if dir == "left" {
		if newPosition <= 0 {
			m.streamer.Seek(0)
		} else {
			m.streamer.Seek(newPosition)
		}
	} else {
		if newPosition >= m.streamer.Len() {
			// Its the end of the song, we should move on to the next one
			speaker.Unlock()
			return m.startSong(m.currentSongIndex + 1)

		} else {
			m.streamer.Seek(newPosition)
		}
	}

	speaker.Unlock()
	return m, tea.Batch(TickNow(), cmd)
}

func (m model) togglePauseState() (model, tea.Cmd) {
	if m.ctrl == nil {
		return m, nil
	}
	m.ctrl.Paused = !m.ctrl.Paused

	return m, nil
}

func (m model) nextPrevSong(goForward bool) (model, tea.Cmd) {
	nextSongIndex := m.currentSongIndex
	if goForward {
		nextSongIndex += 1
	} else {
		nextSongIndex -= 1
	}
	return m.startSong(nextSongIndex)
}

func (m *model) updateUIWidths(termWidth int) {
	colsToUpdate := []string{"Title", "Artist", "File Name"}
	width := int(math.Floor(float64(termWidth-30) / float64(len(colsToUpdate))))

	m.progress.Width = termWidth - 11

	cols := m.table.Columns()
	for i := range cols {
		for _, name := range colsToUpdate {
			if name == cols[i].Title {
				cols[i].Width = width
			}
		}
	}

	m.table.SetColumns(cols)
}

func (m model) toggleDirSelect() (model, tea.Cmd) {
	m.selectDir = !m.selectDir
	if m.selectDir {
		m.table.Blur()
	} else {
		m.table.Focus()
	}

	return m, nil
}

func (m model) toggleVisualizer() (model, tea.Cmd) {
	if m.visualizer != nil {
		enabled := !m.visualizer.IsEnabled()
		m.visualizer.Enable(enabled)
	}
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.debugMessages = append(m.debugMessages, fmt.Sprintf("WindowSizeMsg: %d x %d", msg.Width, msg.Height))
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		m.updateUIWidths(msg.Width)
		m.table.SetWidth(msg.Width - 2) // Set the table width to the terminal width
		if m.visualizer != nil {
			m.visualizer.SetSize(msg.Width-8, 5)
		}
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.PlayPause):
			return m.togglePauseState()
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.ChangeDir):
			return m.toggleDirSelect()
		case key.Matches(msg, m.keys.ToggleVisualizer):
			return m.toggleVisualizer()
		}

		if !m.selectDir {
			switch {
			case key.Matches(msg, m.keys.NextSong):
				if len(m.table.Rows()) == 0 {
					return m, nil
				}
				return m.nextPrevSong(true)
			case key.Matches(msg, m.keys.PrevSong):
				if len(m.table.Rows()) == 0 {
					return m, nil
				}
				return m.nextPrevSong(false)
			case key.Matches(msg, m.keys.Left):
				return m.seek(cmd, "left")
			case key.Matches(msg, m.keys.Right):
				return m.seek(cmd, "right")
			case key.Matches(msg, m.keys.PlaySong):
				if len(m.table.Rows()) == 0 {
					return m, nil
				}
				return m.startSong(m.table.Cursor())
			}
		}
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case tickNowMsg:
		return m.tickAction(true)
	case tickMsg:
		return m.tickAction(false)
	}
	m.table, cmd = m.table.Update(msg)
	m.dirpicker, cmd = m.dirpicker.Update(msg)

	if didSelect, path := m.dirpicker.DidSelectDir(msg); didSelect {
		m.currentDir = path
		m.selectDir = false
		m.debugMessages = append(m.debugMessages, "Selected directory: "+path)
		m.table.SetRows(player.GetSongsInDir(path))
		m.table.Focus()
	}

	return m, cmd
}
