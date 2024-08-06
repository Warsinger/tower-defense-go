package components

import (
	"fmt"
	"image"
	"image/color"
	"tower-defense/assets"
	"tower-defense/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

// Component is any struct that holds some kind of data.
type PositionData struct {
	X, Y int
}

type VelocityData struct {
	X, Y    int
	blocked bool
}

type Renderer interface {
	Draw(screen *ebiten.Image, entry *donburi.Entry)
	GetRect(entry *donburi.Entry) image.Rectangle
}

type RenderData struct {
	renderers []Renderer
}

type InfoRenderData struct {
}

var Position = donburi.NewComponentType[PositionData]()
var Velocity = donburi.NewComponentType[VelocityData]()
var Render = donburi.NewComponentType[RenderData]()
var InfoRender = donburi.NewComponentType[InfoRenderData]()

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

		percentHealth := float32(health.Health) / float32(health.MaxHealth)
		// draw a green filled rect with health below entity the height of the text
		const barHeight = 4
		vector.StrokeRect(screen, float32(rect.Min.X), float32(rect.Max.Y), float32(rect.Dx()), barHeight, 1, color.RGBA{0, 255, 0, 255}, true)
		vector.DrawFilledRect(screen, float32(rect.Min.X), float32(rect.Max.Y), float32(rect.Dx())*percentHealth, barHeight, color.RGBA{0, 255, 0, 255}, true)

		op.GeoM.Translate(float64(rect.Min.X)+(float64(rect.Dx())-textWidth)/2, float64(rect.Max.Y+barHeight))
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

func DetectCollisions(world donburi.World, rect image.Rectangle, excludeFilter filter.LayoutFilter) *donburi.Entry {
	var collision *donburi.Entry = nil
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Render, Position),
			filter.Not(excludeFilter),
		),
	)

	query.Each(world, func(testEntry *donburi.Entry) {
		if collision == nil {
			testRect := Render.Get(testEntry).GetRect(testEntry)
			if rect.Overlaps(testRect) {
				collision = testEntry
			}
		}
	})
	return collision
}
