package scenes

import (
	"fmt"

	"tower-defense/assets"
	comp "tower-defense/components"
	"tower-defense/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"

	"github.com/ebitenui/ebitenui"
)

type TitleScene struct {
	width           int
	height          int
	gameStats       *GameStats
	gameOptions     *config.ConfigData
	newGameCallback NewGameCallback
	world           donburi.World
	ui              *ebitenui.UI
}

var controller = Controller{}

type NewGameCallback func(broadcast bool, controller *Controller, gameOptions *config.ConfigData) error

func NewTitleScene(world donburi.World, width, height int, gameStats *GameStats, gameOptions *config.ConfigData, newGameCallback NewGameCallback) (*TitleScene, error) {
	title := &TitleScene{world: world, width: width, height: height, gameStats: gameStats, gameOptions: gameOptions, newGameCallback: newGameCallback}
	title.ui = initUI(title.gameOptions, newGameCallback, title.handleOptions)
	return title, nil
}

func (t *TitleScene) handleOptions(gameOptions *config.ConfigData) {
	t.gameOptions = gameOptions
	if len(gameOptions.ServerPort) != 0 {
		controller.StartServer(t.world, gameOptions, t.newGameCallback)
	} else if len(gameOptions.ClientHostPort) != 0 {

		controller.StartClient(t.world, gameOptions, t.newGameCallback)
	}
}
func (t *TitleScene) Update() error {
	t.ui.Update()

	if !IsModalOpen() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		return t.newGameCallback(true, &controller, t.gameOptions)
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

	str := fmt.Sprintf("HIGH %05d", t.gameStats.HighScore)
	nextY := comp.DrawTextLines(screen, assets.ScoreFace, str, width, comp.TextBorder, text.AlignEnd, text.AlignStart)
	str = fmt.Sprintf("High Creep Level %d\nHigh Tower Level %d\n", t.gameStats.HighCreepLevel, t.gameStats.HighTowerLevel)
	_ = comp.DrawTextLines(screen, assets.InfoFace, str, width, nextY, text.AlignEnd, text.AlignStart)

	str = "TOWER DEFENSE"
	_ = comp.DrawTextLines(screen, assets.ScoreFace, str, width, 100, text.AlignCenter, text.AlignStart)

	str = "Click to place towers Cost: $50"
	nextY = comp.DrawTextLines(screen, assets.ScoreFace, str, width, 250, text.AlignCenter, text.AlignStart)

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
	str = "Click 'Start Game' or press space to start"
	_ = comp.DrawTextLines(screen, assets.ScoreFace, str, width, nextY, text.AlignCenter, text.AlignStart)

	// draw UI elements
	t.ui.Draw(screen)
}
