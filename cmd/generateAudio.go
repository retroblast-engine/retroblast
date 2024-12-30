package cmd

import (
	"fmt"
)

func generateAudioGo(importPath string) string {
	return fmt.Sprintf(`package audio
	
import (
	"fmt"
	"io"

	"%s/settings"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	sampleRate = 48000 // hardcoded value
)

type MusicType int

const (
	MP3 MusicType = iota
	OGG
	WAV
)

type AudioStream interface {
	io.ReadSeeker
	Length() int64
}

type AudioHandler struct {
	context   *audio.Context
	player    *audio.Player
	musicType MusicType
	volume    int
}

func InitAudioHandler(source string, musicType MusicType) (*AudioHandler, error) {
	audioContext := audio.NewContext(sampleRate)
	var stream AudioStream

	audioFile, err := settings.Assets.Open(settings.AssetsAudioPath + source)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file", err)
	}
	defer audioFile.Close()

	switch musicType {
	case MP3:
		var err error
		stream, err = mp3.DecodeF32(audioFile)
		if err != nil {
			return nil, err
		}

	case OGG:
		var err error
		stream, err = vorbis.DecodeF32(audioFile)
		if err != nil {
			return nil, err
		}

	case WAV:
		var err error
		stream, err = wav.DecodeF32(audioFile)
		if err != nil {
			return nil, err
		}
	default:
		panic("not reached")
	}

	audioPlayer, err := audioContext.NewPlayerF32(stream)
	if err != nil {
		return nil, err
	}

	audioHandler := &AudioHandler{
		context:   audioContext,
		player:    audioPlayer,
		musicType: musicType,
		volume:    128,
	}

	audioHandler.player.Play()

	return audioHandler, nil
}

func (audioHandler *AudioHandler) VolumeUp() {
	audioHandler.volume++

	if audioHandler.volume >= 100 {
		audioHandler.volume = 100
	}
	audioHandler.updateVolume()
}

func (audioHandler *AudioHandler) VolumeDown() {
	audioHandler.volume--

	if audioHandler.volume < 0 {
		audioHandler.volume = 0
	}
	audioHandler.updateVolume()
}

func (audioHandler *AudioHandler) updateVolume() {
	audioHandler.player.SetVolume(float64(audioHandler.volume) / 128)
}

func (audioHandler *AudioHandler) Close() error {
	return audioHandler.player.Close()
}
`, importPath)
}
