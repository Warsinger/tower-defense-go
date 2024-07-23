package components

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
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
}

func (s *SpriteData) GetRect(entry *donburi.Entry) image.Rectangle {
	pos := Position.Get(entry)
	rect := s.image.Bounds()
	return rect.Add(image.Pt(pos.x, pos.y))
}
