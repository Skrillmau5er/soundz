package player

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhowden/tag"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

var supportedFileExt = []string{".wav", ".mp3", ".flac", ".ogg"}

func GetSongsInDir(currentDir string) []table.Row {
	rows := []table.Row{}
	dirEntry, err := os.ReadDir(currentDir)

	if err != nil {
		fmt.Println("Error reading current directory: ", err)
	}

	for _, v := range dirEntry {
		isDir := v.IsDir()
		if isDir {
			continue
		}
		name := v.Name()
		ext := filepath.Ext(name)

		supported := false
		for _, v := range supportedFileExt {
			if v == ext {
				supported = true
				break
			}
		}

		if !supported {
			continue
		}

		filePath := currentDir + "/" + name

		length := GetFileAudioLength(filePath, ext)

		formatedLength := fmt.Sprintf("%v:%02d", math.Floor(float64(length)/60), length%60)

		file, err := os.Open(filePath)

		if err != nil {
			fmt.Printf("Error opening file: %s", name)
		}
		m, err := tag.ReadFrom(file)
		if err != nil {
			rows = append(rows, []string{"", "", "", name, formatedLength, ext})
			continue
		}

		rows = append(rows, []string{"", m.Title(), m.Artist(), name, formatedLength, ext})
	}

	return rows
}

func GetPosAndLen(streamer beep.StreamSeekCloser, sampleRate beep.SampleRate) (int, int) {
	length := streamer.Len() / int(sampleRate)
	pos := streamer.Position() / int(sampleRate)
	return pos, length
}

func PlaySongCmd(filePath string, extType string) tea.Cmd {
	return func() tea.Msg {
		PlaySong(filePath, extType)
		return nil
	}
}

func GetFileAudioLength(filePath string, extType string) int {
	streamer, format, err := OpenFileAndDecode(filePath, extType)

	if err != nil {
		log.Fatal(err)
	}

	length := streamer.Len() / int(format.SampleRate)
	return length
}

func PlaySong(filePath string, extType string) (*beep.Ctrl, beep.StreamSeekCloser, beep.Format, beep.SampleRate) {
	streamer, format, err := OpenFileAndDecode(filePath, extType)
	originalSampleRate := format.SampleRate

	if err != nil {
		log.Fatal(err)
	}

	sampleRate := beep.SampleRate(44100)

	if format.SampleRate != sampleRate {
		resampled := beep.Resample(4, format.SampleRate, sampleRate, streamer)
		format.SampleRate = sampleRate
		ctrl := &beep.Ctrl{Streamer: resampled, Paused: false}
		speaker.Play(ctrl)
		return ctrl, streamer, format, originalSampleRate
	}

	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	speaker.Play(ctrl)

	return ctrl, streamer, format, originalSampleRate
}

// Visualizer interface for audio visualization
type Visualizer interface {
	Update(samples [][2]float64)
}

// PlaySongWithVisualizer plays a song and sends audio samples to the visualizer
func PlaySongWithVisualizer(filePath string, extType string, viz Visualizer) (*beep.Ctrl, beep.StreamSeekCloser, beep.Format, beep.SampleRate) {
	streamer, format, err := OpenFileAndDecode(filePath, extType)
	originalSampleRate := format.SampleRate

	if err != nil {
		log.Fatal(err)
	}

	sampleRate := beep.SampleRate(44100)

	var finalStreamer beep.Streamer

	if format.SampleRate != sampleRate {
		// Resample first
		resampled := beep.Resample(4, format.SampleRate, sampleRate, streamer)
		format.SampleRate = sampleRate
		// Then wrap with visualizer (AFTER resampling)
		finalStreamer = &visualizerStreamer{
			Streamer:   resampled,
			visualizer: viz,
		}
	} else {
		// Wrap with visualizer directly
		finalStreamer = &visualizerStreamer{
			Streamer:   streamer,
			visualizer: viz,
		}
	}

	ctrl := &beep.Ctrl{Streamer: finalStreamer, Paused: false}
	speaker.Play(ctrl)

	return ctrl, streamer, format, originalSampleRate
}

// visualizerStreamer wraps a beep.Streamer to capture audio samples for visualization
type visualizerStreamer struct {
	beep.Streamer
	visualizer Visualizer
}

// Stream streams audio and captures samples for visualization
func (vs *visualizerStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = vs.Streamer.Stream(samples)
	if ok && n > 0 {
		// Send samples to visualizer
		vs.visualizer.Update(samples[:n])
	}
	return n, ok
}
