package components

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// Component is any struct that holds some kind of data.
type PositionData struct {
	x, y int
}

type VelocityData struct {
	x, y int
}

// ComponentType represents kind of component which is used to create or query entities.
var Position = donburi.NewComponentType[PositionData]()
var Velocity = donburi.NewComponentType[VelocityData]()

type Renderer interface {
	Draw(screen *ebiten.Image, entry *donburi.Entry)
	GetRect(entry *donburi.Entry) image.Rectangle
}

type RenderData struct {
	renderers []Renderer
}

var Render = donburi.NewComponentType[RenderData]()

func NewRenderer(renderers ...Renderer) *RenderData {
	return &RenderData{renderers: renderers}
}

func (r *RenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	for _, render := range r.renderers {
		render.Draw(screen, entry)
	}
}

func (r *RenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	return r.GetPrimaryRenderer().GetRect(entry)
}

func (r *RenderData) GetPrimaryRenderer() Renderer {
	return r.renderers[0]
}
