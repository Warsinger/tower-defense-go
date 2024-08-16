package scenes

import (
	"image/color"
	"tower-defense/assets"
	comp "tower-defense/components"
	"tower-defense/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

type ViewerScene struct {
	world     donburi.World
	width     int
	height    int
	image     *ebiten.Image
	config    *config.ConfigData
	translate bool
}

func NewViewerScene(world donburi.World, width, height int, gameOptions *config.ConfigData, translate bool) (*ViewerScene, error) {
	return &ViewerScene{
		world:     world,
		width:     width,
		height:    height,
		config:    config.NewConfig(world, gameOptions.Debug),
		translate: translate,
	}, nil
}

func (v *ViewerScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		v.config.GridLines = !v.config.GridLines
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		v.config.Debug = !v.config.Debug
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

	comp.DrawBoard(v.image, v.world, v.config, v.DrawText)

	if v.translate {
		vector.StrokeLine(v.image, 0, 0, 0, float32(v.height), 3, color.White, true)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(v.width), 0)
		screen.DrawImage(v.image, opts)
	}
}

func (v *ViewerScene) DrawText(image *ebiten.Image) {
	str := "Viewer Mode"
	comp.DrawTextLines(image, assets.InfoFace, str, float64(v.width), comp.TextBorder, text.AlignEnd, text.AlignStart)

	bssEntry, ok := comp.BattleState.First(v.world)
	if ok {
		bss := comp.BattleState.Get(bssEntry)
		bss.Draw(image, float64(v.width), float64(v.height))
	}
}
