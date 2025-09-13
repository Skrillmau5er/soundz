package ui

import (
	"fmt"
	"math"

	"example.com/soundz/internal/player"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep/v2/speaker"
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

	ctrl, streamer, format, songSampleRate := player.PlaySong(row[3], row[5])
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

	m.progress.Width = termWidth - 10

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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateUIWidths(msg.Width)
		m.table.SetWidth(msg.Width - 2) // Set the table width to the terminal width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m.startSong(m.table.Cursor())
		case " ":
			return m.togglePauseState()
		case "left", "h":
			return m.seek(cmd, "left")
		case "right", "l":
			return m.seek(cmd, "right")
		case "n":
			return m.nextPrevSong(true)
		case "p":
			return m.nextPrevSong(false)

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
	return m, cmd
}
