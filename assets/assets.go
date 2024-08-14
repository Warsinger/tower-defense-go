package assets

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
)

var images map[string]*ebiten.Image = make(map[string]*ebiten.Image)
var sounds map[string]*audio.Player = make(map[string]*audio.Player)

//go:embed images sounds
var fs embed.FS

var ScoreFace *text.GoTextFace
var InfoFace *text.GoTextFace
var GoFace font.Face
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

	// font for the UI
	ttfFont, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		return err
	}

	GoFace = truetype.NewFace(ttfFont, &truetype.Options{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	return nil
}

func loadImages() error {
	dir, err := fs.ReadDir("images")
	if err != nil {
		return nil
	}
	for _, file := range dir {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".png") {
			err := loadImageAsset(file.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func loadImageAsset(name string) error {
	filename := fmt.Sprintf("images/%s", name)
	data, err := fs.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read embedded image %v: %v", name, err)
		return err
	}
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("failed to load image %v: %v", name, err)
		return err
	}

	index := strings.TrimSuffix(name, filepath.Ext(name))
	images[index] = img
	return nil
}

func loadAudio() error {
	dir, err := fs.ReadDir("sounds")
	if err != nil {
		return nil
	}
	audioContext = audio.NewContext(44100)
	for _, file := range dir {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".wav") {
			err := loadAudioAsset(file.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func loadAudioAsset(name string) error {
	filename := fmt.Sprintf("sounds/%s", name)
	data, err := fs.ReadFile(filename)
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
	index := strings.TrimSuffix(name, filepath.Ext(name))
	sounds[index] = player

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
	image := images[name]
	if image == nil {
		log.Fatalf("Image asset not found %v", name)
	}
	return image
}
