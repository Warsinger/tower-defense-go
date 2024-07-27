package scenes

import (
	"fmt"
	"tower-defense/assets"
	comp "tower-defense/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TitleScene struct {
	width           int
	height          int
	highScore       int
	newGameCallback func(bool) error
	viewer          bool
}

func NewTitleScene(width, height, highScore int, newGameCallback func(bool) error) (*TitleScene, error) {
	return &TitleScene{width: width, height: height, highScore: highScore, newGameCallback: newGameCallback}, nil
}

func (t *TitleScene) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		return t.newGameCallback(t.viewer)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		t.viewer = !t.viewer
	}
	return nil
}

func (t *TitleScene) Draw(screen *ebiten.Image) {
	screen.Clear()

	backgroundImage := assets.GetImage("backgroundV")
	opts := &ebiten.DrawImageOptions{}
	screen.DrawImage(backgroundImage, opts)

	// draw high score
	str := fmt.Sprintf("HIGH %05d", t.highScore)
	op := &text.DrawOptions{}
	x, _ := text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(float64(t.width)-x-comp.TextBorder, comp.TextBorder)
	text.Draw(screen, str, assets.ScoreFace, op)

	halfWidth := float64(t.width / 2)
	str = "TOWER DEFENSE"
	op = &text.DrawOptions{}
	x, _ = text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(halfWidth-x/2, 100)
	text.Draw(screen, str, assets.ScoreFace, op)

	str = "Click to place towers"
	op = &text.DrawOptions{}
	x, y := text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(halfWidth-x/2, 300)
	text.Draw(screen, str, assets.ScoreFace, op)

	towerImage := assets.GetImage("tower")
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(halfWidth-float64(towerImage.Bounds().Dx()/2), 300+y)
	screen.DrawImage(towerImage, opts)

	str = "Protect your base from aliens"
	op = &text.DrawOptions{}
	x, y = text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(halfWidth-x/2, 400)
	text.Draw(screen, str, assets.ScoreFace, op)

	creepImage := assets.GetImage("creep2")
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(halfWidth-float64(creepImage.Bounds().Dx()/2), 400+y)
	screen.DrawImage(creepImage, opts)

	str = "Click or press space to start"
	op = &text.DrawOptions{}
	x, y = text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(halfWidth-x/2, 600)
	text.Draw(screen, str, assets.ScoreFace, op)

	str = fmt.Sprintf("Viewer mode %v (Press V to toggle)", t.viewer)
	op = &text.DrawOptions{}
	x, _ = text.Measure(str, assets.InfoFace, op.LineSpacing)
	op.GeoM.Translate(halfWidth-x/2, 600+y)
	text.Draw(screen, str, assets.InfoFace, op)
}
