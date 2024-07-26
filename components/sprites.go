package components

import (
	"image"
	"image/color"
	"tower-defense/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

type SpriteData struct {
	image *ebiten.Image
}

func (s *SpriteData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	pos := Position.Get(entry)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(pos.x), float64(pos.y))
	screen.DrawImage(s.image, opts)

	config := config.GetConfig(entry.World)
	if config.IsDebug() {
		// draw bounding rect used for collision detection
		rect := s.GetRect(entry)
		vector.StrokeRect(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), 1, color.White, true)
	}
}

func (s *SpriteData) GetRect(entry *donburi.Entry) image.Rectangle {
	pos := Position.Get(entry)
	rect := s.image.Bounds()
	return rect.Add(image.Pt(pos.x, pos.y))
}
