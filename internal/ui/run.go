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
	"github.com/skrillmau5er/soundz/internal/components/dirpicker"
)

const sampleRate = beep.SampleRate(44100)

func Run() {
	dir, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}
	t := GetTable(dir)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	dp := dirpicker.New()
	dp.CurrentDirectory, _ = os.UserHomeDir()
	dp.SetHeight(15)

	m := model{
		table:    t,
		ctrl:     nil,
		streamer: nil,
		progress: progress.New(progress.WithDefaultGradient(),
			progress.WithoutPercentage()),
		currentSongIndex: -1,
		keys:             keys,
		help:             help.New(),
		dirpicker:        dp,
		selectDir:        false,
		currentDir:       dir,
	}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
