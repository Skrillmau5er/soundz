package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) renderArt() string {
	return ` ▗▄▄▖ ▄▄▄  █  ▐▌▄▄▄▄     ▐▌▄▄▄▄▄ 
▐▌   █   █ ▀▄▄▞▘█   █    ▐▌ ▄▄▄▀ 
 ▝▀▚▖▀▄▄▄▀      █   █ ▗▞▀▜▌█▄▄▄▄ 
▗▄▄▞▘                 ▝▚▄▟▌      `
}

func (m model) renderTable() string {
	const debugMessagesHeight = 10
	if m.selectDir {
		m.table.SetHeight(m.termHeight - 30 - debugMessagesHeight)
	} else {
		m.table.SetHeight(m.termHeight - 14 - debugMessagesHeight)
	}
	return baseStyle.Render(m.table.View())
}

func (m model) renderProgress() string {
	if m.streamer == nil {
		return "\n"
	}

	songPos := m.songPos
	if songPos == "" {
		songPos = strings.Repeat(" ", 4)
	}
	songLength := m.songLength
	if songLength == "" {
		songLength = strings.Repeat(" ", 4)
	}
	style := lipgloss.NewStyle().PaddingBottom(1)
	return style.Render(lipgloss.JoinHorizontal(lipgloss.Top, songPos, " ", m.progress.View(), " ", songLength))
}

func (m model) renderVizualizer() string {
	if m.streamer == nil {
		return ""
	}
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1).
		Width(m.termWidth - 4)

	return style.Render(m.visualizer.View())
}

func (m model) renderHelp() string {
	style := lipgloss.NewStyle().PaddingLeft(4).PaddingBottom(1)

	helpView := m.help.View(m.keys)
	return style.Render(helpView)
}

func (m model) renderHeader() string {
	style := lipgloss.NewStyle().Padding(1)
	return style.Render(lipgloss.JoinHorizontal(lipgloss.Bottom, m.renderArt(), m.renderHelp()))
}

func (m model) renderDirPicker() string {
	var s strings.Builder
	s.WriteString("\n  ")
	if m.selectDir {
		if m.err != nil {
			s.WriteString(m.dirpicker.Styles.DisabledFile.Render(m.err.Error()))
		} else if m.currentDir == "" {
			s.WriteString("Select a directory:")
		} else {
			s.WriteString("Currently selected directory: " + m.dirpicker.Styles.Selected.Render(m.currentDir))
		}

		// Add navigation instructions
		instructions := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			PaddingLeft(2)

		navText := "• (↑/j)/(↓/k): navigate \n• (←/h)/(→/l): exit/enter directory \n• g: first \n• G: last"
		s.WriteString("\n" + instructions.Render(navText))

		// Use dynamic height based on terminal size, but cap it at a reasonable maximum
		// Account for: header (~8 lines) + progress (~2 lines) + instructions (~1 line) + padding (~3 lines)
		availableHeight := m.termHeight - 14 // Leave space for header, progress, instructions, and padding
		maxHeight := 12                      // Maximum number of items to show
		if availableHeight > maxHeight {
			availableHeight = maxHeight
		}
		// Ensure we have at least 3 items visible
		if availableHeight < 3 {
			availableHeight = 3
		}
		m.dirpicker.SetHeight(availableHeight)

		s.WriteString("\n\n" + m.dirpicker.View() + "\n")
	} else {
		s.WriteString(fmt.Sprintf("Viewing audio files from directory: %v", m.currentDir))
	}
	return s.String()
}

func (m model) renderDebugMessages() string {
	if len(m.debugMessages) == 0 {
		return ""
	}
	return strings.Join(m.debugMessages, "\n")
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.renderHeader(),
		m.renderProgress(),
		m.renderVizualizer(),
		m.renderDirPicker(),
		m.renderTable(),
		m.renderDebugMessages(),
	)
}
