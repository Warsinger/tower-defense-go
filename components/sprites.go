package components

import (
	"image"
	"image/color"
	"tower-defense/assets"
	"tower-defense/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

type SpriteData struct {
	name  string
	image *ebiten.Image
}

func NewSprite(name string) *SpriteData {
	image := assets.GetImage(name)
	return &SpriteData{name: name, image: image}
}

func (s *SpriteData) ensureImage() {
	if s.image == nil {
		s.image = assets.GetImage(s.name)
	}
}
func (s *SpriteData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	s.ensureImage()
	pos := Position.Get(entry)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(pos.X), float64(pos.Y))
	screen.DrawImage(s.image, opts)

	config := config.GetConfig(entry.World)
	if config.IsDebug() {
		// draw bounding rect used for collision detection
		rect := s.GetRect(entry)
		vector.StrokeRect(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), 1, color.White, true)
	}
}

func (s *SpriteData) GetRect(entry *donburi.Entry) image.Rectangle {
	s.ensureImage()
	pos := Position.Get(entry)
	rect := s.image.Bounds()
	return rect.Add(image.Pt(pos.X, pos.Y))
}
