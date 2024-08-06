package scenes

import (
	"bytes"
	"fmt"
	"image/color"
	"tower-defense/assets"
	comp "tower-defense/components"
	"tower-defense/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/debug"
	"github.com/yohamta/donburi/filter"
)

type ViewerScene struct {
	world     donburi.World
	width     int
	height    int
	image     *ebiten.Image
	config    *config.ConfigData
	translate bool
}

func NewViewerScene(world donburi.World, width, height int, debug, translate bool) (*ViewerScene, error) {
	return &ViewerScene{
		world:     world,
		width:     width,
		height:    height,
		config:    config.NewConfig(world, debug),
		translate: translate,
	}, nil
}

func (v *ViewerScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		v.config.SetGridLines(!v.config.IsGridLines())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		v.config.SetDebug(!v.config.IsDebug())
	}
	return nil
}

func (v *ViewerScene) Draw(screen *ebiten.Image) {
	if v.translate {
		if v.image == nil {
			v.image = ebiten.NewImage(v.width, v.height)
		}
	} else {
		v.image = screen
	}
	v.image.Clear()

	img := assets.GetImage("backgroundV")
	opts := &ebiten.DrawImageOptions{}
	v.image.DrawImage(img, opts)

	if v.translate {
		vector.StrokeLine(v.image, 0, 0, 0, float32(v.height), 3, color.White, true)
	}

	DebugPrint(v)

	// query for all entities with Position component
	query := donburi.NewQuery(filter.Contains(comp.Position))
	// ebitenutil.DebugPrint(v.image, fmt.Sprintf("entities in draw %d\n", query.Count(v.world)))
	// draw all entities
	query.Each(v.world, func(entry *donburi.Entry) {
		comp.DrawEntry(v.image, entry)
	})

	v.DrawText(v.image)

	if v.translate {
		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(v.width), 0)
		screen.DrawImage(v.image, opts)
	}
}

func DebugPrint(v *ViewerScene) {
	var out bytes.Buffer
	out.WriteString("Entity Counts:\n")
	for _, c := range debug.GetEntityCounts(v.world) {
		out.WriteString(c.String())
		out.WriteString("\n")
	}
	out.WriteString("\n")
	msg := fmt.Sprint(out.String())

	ebitenutil.DebugPrint(v.image, msg)
}

func (v *ViewerScene) DrawText(image *ebiten.Image) {
	str := "Viewer Mode"
	op := &text.DrawOptions{}
	x, _ := text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(float64(v.width)-x-comp.TextBorder, comp.TextBorder)
	text.Draw(image, str, assets.InfoFace, op)
}
