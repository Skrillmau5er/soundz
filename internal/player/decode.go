package player

import (
	"log"
	"os"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

func OpenFileAndDecode(filePath string, extType string) (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open("../../audio_samples/" + filePath)

	var streamer beep.StreamSeekCloser
	var format beep.Format

	if err != nil {
		log.Fatal(err)
	}

	switch extType {
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".flac":
		streamer, format, err = flac.Decode(f)
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
	}

	return streamer, format, err
}
