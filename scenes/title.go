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
	width := float64(t.width)
	halfWidth := width / 2

	str := fmt.Sprintf("HIGH %05d", t.highScore)
	_ = comp.DrawTextLines(screen, assets.ScoreFace, str, width, comp.TextBorder, text.AlignEnd, text.AlignStart)

	str = "TOWER DEFENSE"
	_ = comp.DrawTextLines(screen, assets.ScoreFace, str, width, 100, text.AlignCenter, text.AlignStart)

	str = "Click to place towers Cost: $50"
	nextY := comp.DrawTextLines(screen, assets.ScoreFace, str, width, 250, text.AlignCenter, text.AlignStart)

	str = "Mouse over a tower\nPress H to heal to full Cost: $25\nPress U to upgrade and heal to full Cost: $50\nMax upgrade level is 5 (+1 for every 20 upgrades)"
	nextY = comp.DrawTextLines(screen, assets.InfoFace, str, width, nextY, text.AlignCenter, text.AlignStart)

	towerImage := assets.GetImage("tower")
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(halfWidth-float64(towerImage.Bounds().Dx()/2), nextY)
	screen.DrawImage(towerImage, opts)
	nextY += float64(towerImage.Bounds().Dy()) + 10

	str = "Protect your base from aliens and earn $$"
	nextY = comp.DrawTextLines(screen, assets.ScoreFace, str, width, nextY, text.AlignCenter, text.AlignStart)

	const creepCount = 4
	const creepSize = 48
	x := halfWidth - creepSize*creepCount/2
	for i := 1; i <= creepCount; i++ {
		creepImage := assets.GetImage(fmt.Sprintf("creep%v", i))
		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(x, nextY)
		screen.DrawImage(creepImage, opts)
		x += float64(creepImage.Bounds().Dx())
	}

	nextY = 600
	str = "Click or press space to start"
	nextY = comp.DrawTextLines(screen, assets.ScoreFace, str, width, nextY, text.AlignCenter, text.AlignStart)

	str = fmt.Sprintf("Viewer mode %v (Press V to toggle)", t.viewer)
	_ = comp.DrawTextLines(screen, assets.InfoFace, str, width, nextY, text.AlignCenter, text.AlignStart)
}
