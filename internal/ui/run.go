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
	"github.com/skrillmau5er/soundz/internal/components/visualizer"
)

const sampleRate = beep.SampleRate(44100)

func Run() {
	dir, err := os.Getwd()
	fmt.Println("dir", dir)
	if err != nil {
		os.Exit(1)
	}
	t := GetTable(dir)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	dp := dirpicker.New()
	dp.CurrentDirectory = dir
	dp.SetHeight(15)

	viz := visualizer.New()

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
		visualizer:       &viz,
	}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
