package assets

import (
	"bytes"
	"embed"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var images map[string]*ebiten.Image = make(map[string]*ebiten.Image)
var sounds map[string]*audio.Player = make(map[string]*audio.Player)

//go:embed images sounds
var fs embed.FS

var ScoreFace *text.GoTextFace
var InfoFace *text.GoTextFace
var audioContext *audio.Context

func LoadAssets() error {
	err := loadImages()
	if err != nil {
		return err
	}

	err = loadFonts()
	if err != nil {
		return err
	}

	err = loadAudio()
	if err != nil {
		return err
	}

	return nil
}
func loadFonts() error {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		return err
	}

	ScoreFace = &text.GoTextFace{
		Source: s,
		Size:   24,
	}
	InfoFace = &text.GoTextFace{
		Source: s,
		Size:   12,
	}
	return nil
}

func loadImages() error {
	names := []string{"creep1", "creep2", "tower", "base", "backgroundH", "backgroundV"}
	for _, name := range names {
		err := loadImageAsset(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadAudio() error {
	audioContext = audio.NewContext(44100)
	names := []string{"explosion", "killed", "shoot"}
	for _, name := range names {
		err := loadAudioAsset(name)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadAudioAsset(name string) error {
	filepath := fmt.Sprintf("sounds/%s.wav", name)
	data, err := fs.ReadFile(filepath)
	if err != nil {
		return err
	}

	d, err := wav.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		return err
	}
	player, err := audioContext.NewPlayer(d)
	if err != nil {
		return err
	}
	sounds[name] = player

	return nil
}

func loadImageAsset(name string) error {
	filepath := fmt.Sprintf("images/%s.png", name)
	data, err := fs.ReadFile(filepath)
	if err != nil {
		log.Fatalf("failed to read embedded image %v: %v", name, err)
		return err
	}
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("failed to load image %v: %v", name, err)
		return err
	}

	images[name] = img
	return nil
}

func PlaySound(name string) {
	sound := sounds[name]
	if sound != nil {
		sound.Rewind()
		sound.Play()
	}
}

func GetImage(name string) *ebiten.Image {
	return images[name]
}
