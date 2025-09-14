package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

const sampleRate = beep.SampleRate(44100)

func Run() {

	t := GetTable()
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	m := model{
		table:    t,
		ctrl:     nil,
		streamer: nil,
		progress: progress.New(progress.WithDefaultGradient(),
			progress.WithoutPercentage()),
		currentSongIndex: -1,
		keys:             keys,
		help:             help.New(),
	}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
