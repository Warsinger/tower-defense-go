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

type SpriteRenderData struct {
	image *ebiten.Image
	Name  string
}

var SpriteRender = donburi.NewComponentType[SpriteRenderData]()

func (s *SpriteRenderData) GetImage(entry *donburi.Entry) *ebiten.Image {
	if s.image == nil {
		s.image = assets.GetImage(s.Name)
	}
	return s.image
}
func (s *SpriteRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	img := s.GetImage(entry)
	pos := Position.Get(entry)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(pos.X), float64(pos.Y))
	screen.DrawImage(img, opts)

	config := config.GetConfig(entry.World)

	rect := s.GetRect(entry)
	var clr color.Color = nil
	if (entry.HasComponent(Tower) && image.Pt(ebiten.CursorPosition()).In(rect)) || config.IsDebug() {
		clr = color.White
		vector.StrokeRect(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), 1, clr, true)
	}
}

func (s *SpriteRenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	img := s.GetImage(entry)
	pos := Position.Get(entry)
	rect := img.Bounds()
	return rect.Add(image.Pt(pos.X, pos.Y))
}
