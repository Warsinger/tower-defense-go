package components

import (
	"fmt"
	"image"
	"tower-defense/assets"
	"tower-defense/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
)

// Component is any struct that holds some kind of data.
type PositionData struct {
	x, y int
}

type VelocityData struct {
	x, y    int
	blocked bool
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

type InfoRenderData struct {
}

func (t *InfoRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	render := Render.Get(entry)
	rect := render.GetRect(entry)

	var textWidth, textHeight float64 = 0, 0
	if entry.HasComponent(Health) {
		health := Health.Get(entry)
		// draw health info centered below the entity
		str := fmt.Sprintf("HP %d", health.Health)
		op := &text.DrawOptions{}
		textWidth, textHeight = text.Measure(str, assets.InfoFace, op.LineSpacing)
		op.GeoM.Translate(float64(rect.Min.X)+(float64(rect.Dx())-textWidth)/2, float64(rect.Max.Y))
		text.Draw(screen, str, assets.InfoFace, op)
	}

	config := config.GetConfig(entry.World)
	if config.IsDebug() {
		if entry.HasComponent(Attack) {
			// draw cooldown info centered below the health
			attack := Attack.Get(entry)
			var cd int = 0
			if attack.inCooldown {
				cd = attack.Cooldown - attack.GetTicker()
			}
			str := fmt.Sprintf("CD %d", cd)
			op := &text.DrawOptions{}
			textWidth, _ = text.Measure(str, assets.InfoFace, op.LineSpacing)
			op.GeoM.Translate(float64(rect.Min.X)+(float64(rect.Dx())-textWidth)/2, float64(rect.Max.Y)+textHeight)
			text.Draw(screen, str, assets.InfoFace, op)
		}
	}
}

func (t *InfoRenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	panic("InfoRenderData.GetRect() unimplemented")
}
