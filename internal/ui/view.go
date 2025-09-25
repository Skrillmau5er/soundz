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
	return `███████╗ ██████╗ ██╗   ██╗███╗   ██╗██████╗ ███████╗
██╔════╝██╔═══██╗██║   ██║████╗  ██║██╔══██╗╚══███╔╝
███████╗██║   ██║██║   ██║██╔██╗ ██║██║  ██║  ███╔╝ 
╚════██║██║   ██║██║   ██║██║╚██╗██║██║  ██║ ███╔╝  
███████║╚██████╔╝╚██████╔╝██║ ╚████║██████╔╝███████╗
╚══════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═══╝╚═════╝ ╚══════╝`
}

func (m model) renderTable() string {
	m.table.SetHeight(m.termHeight - 40)
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
			s.WriteString("Selected file: " + m.dirpicker.Styles.Selected.Render(m.currentDir))
		}
		m.dirpicker.SetHeight(20)
		s.WriteString("\n\n" + m.dirpicker.View() + "\n")
	} else {
		s.WriteString(fmt.Sprintf("Viewing audio files from directory: %v", m.currentDir))
	}
	return s.String()
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.renderHeader(),
		m.renderProgress(),
		m.renderDirPicker(),
		m.renderTable(),
	)
}
