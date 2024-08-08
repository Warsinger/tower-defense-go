package components

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"tower-defense/assets"
	"tower-defense/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/debug"
	"github.com/yohamta/donburi/filter"
)

type InfoRenderData struct {
}

var InfoRender = donburi.NewComponentType[InfoRenderData]()

func (t *InfoRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	rect := GetRect(entry)

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
		var clr color.Color
		if percentHealth >= 0.5 {
			clr = color.RGBA{0, 255, 0, 255}
		} else if percentHealth >= 0.25 {
			clr = color.RGBA{255, 255, 0, 255}
		} else {
			clr = color.RGBA{255, 0, 0, 255}
		}
		vector.StrokeRect(screen, float32(rect.Min.X), float32(rect.Max.Y), float32(rect.Dx()), barHeight, 1, clr, true)
		vector.DrawFilledRect(screen, float32(rect.Min.X), float32(rect.Max.Y), float32(rect.Dx())*percentHealth, barHeight, clr, true)

		op.GeoM.Translate(float64(rect.Min.X)+(float64(rect.Dx())-textWidth)/2, float64(rect.Max.Y+barHeight))
		text.Draw(screen, str, assets.InfoFace, op)
	}

	config := config.GetConfig(entry.World)
	if config.IsDebug() {
		if entry.HasComponent(Attack) {
			// draw power & cooldown info centered below the health
			attack := Attack.Get(entry)
			var cd int = 0
			if attack.inCooldown {
				cd = attack.Cooldown - attack.GetTicker()
			}
			str := fmt.Sprintf("%d/CD %d", attack.Power, cd)
			op := &text.DrawOptions{}
			textWidth, _ = text.Measure(str, assets.InfoFace, op.LineSpacing)
			op.GeoM.Translate(float64(rect.Min.X)+(float64(rect.Dx())-textWidth)/2, float64(rect.Max.Y)+textHeight)
			text.Draw(screen, str, assets.InfoFace, op)
		}
	}
}

func DrawGridLines(screen *ebiten.Image) {
	size := screen.Bounds().Size()
	cellSize := 10
	for i := 0; i <= size.Y; i += cellSize {
		vector.StrokeLine(screen, 0, float32(i), float32(size.X), float32(i), 1, color.White, true)
	}
	for i := 0; i <= size.X; i += cellSize {
		vector.StrokeLine(screen, float32(i), 0, float32(i), float32(size.Y), 1, color.White, true)
	}
}

func DrawEntry(screen *ebiten.Image, entry *donburi.Entry, debug bool) {
	if entry.HasComponent(SpriteRender) {
		render := SpriteRender.Get(entry)
		render.Draw(screen, entry)
		if debug {
			ebitenutil.DebugPrintAt(screen, render.Name, GetRect(entry).Min.X, GetRect(entry).Min.Y-10)
		}
	}
	if entry.HasComponent(InfoRender) {
		info := InfoRender.Get(entry)
		info.Draw(screen, entry)
	}
	if entry.HasComponent(RangeRender) {
		rangeRender := RangeRender.Get(entry)
		rangeRender.Draw(screen, entry)
	}
	if entry.HasComponent(PlayerRender) {
		playerRender := PlayerRender.Get(entry)
		playerRender.Draw(screen, entry)
	}
	if entry.HasComponent(BulletRender) {
		bulletRender := BulletRender.Get(entry)
		bulletRender.Draw(screen, entry)
	}
}

func GetRect(entry *donburi.Entry) image.Rectangle {
	if entry.HasComponent(SpriteRender) {
		render := SpriteRender.Get(entry)
		return render.GetRect(entry)
	} else if entry.HasComponent(BulletRender) {
		render := BulletRender.Get(entry)
		return render.GetRect(entry)
	}
	panic("GetRect() unimplemented for entry without SpriteRender or BulletRender component")
}

func DrawBoard(image *ebiten.Image, world donburi.World, config *config.ConfigData, drawText func(*ebiten.Image)) {
	image.Clear()

	background := assets.GetImage("backgroundV")
	opts := &ebiten.DrawImageOptions{}
	image.DrawImage(background, opts)

	DebugPrint(image, world, config)

	if config.IsGridLines() {
		DrawGridLines(image)
	}

	query := donburi.NewQuery(filter.Contains(Position))

	query.Each(world, func(entry *donburi.Entry) {
		DrawEntry(image, entry, config.IsDebug())
	})

	if drawText != nil {
		drawText(image)
	}
}

func DebugPrint(image *ebiten.Image, world donburi.World, config *config.ConfigData) {
	if !config.IsDebug() {
		return
	}
	var out bytes.Buffer
	out.WriteString("Entity Counts:\n")
	for _, c := range debug.GetEntityCounts(world) {
		out.WriteString(c.String())
		out.WriteString("\n")
	}
	out.WriteString("\n")
	msg := fmt.Sprint(out.String())

	ebitenutil.DebugPrint(image, msg)
}
