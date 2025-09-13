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

func GetSongsInDir() []table.Row {
	rows := []table.Row{}
	dirEntry, err := os.ReadDir("../../audio_samples")

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

		length := GetFileAudioLength(name, ext)

		formatedLength := fmt.Sprintf("%v:%02d", math.Floor(float64(length)/60), length%60)

		file, err := os.Open("../../audio_samples/" + name)

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
