package game

import (
	"tower-defense/assets"
	"tower-defense/config"
	"tower-defense/network"
	"tower-defense/scenes"

	comp "tower-defense/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leap-fish/necs/router"
	"github.com/yohamta/donburi"
)

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}
type GameData struct {
	world                donburi.World
	scenes               []Scene
	width, height, speed int
	gameStats            *comp.GameStats
	startingTowerLevel   int
}

func NewGame(width, height, speed int, startingTowerLevel int, debug, computer, nosound bool) (*GameData, error) {
	err := assets.LoadAssets()
	if err != nil {
		return nil, err
	}

	ebiten.SetWindowTitle("Tower Defense")

	gameStats := comp.LoadStats()

	game := &GameData{world: donburi.NewWorld(), width: width, height: height, speed: speed, gameStats: gameStats, startingTowerLevel: startingTowerLevel}

	err = game.switchToTitle(gameStats, config.NewConfig(game.world, debug, computer, !nosound))
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (g *GameData) switchToBattle(broadcast bool, controller *scenes.Controller, gameOptions *config.ConfigData) error {
	if broadcast {
		router.Broadcast(network.StartGameMessage{})
	}
	clientWorld := controller.GetClientWorld()
	multiplayer := clientWorld != nil
	battle, err := scenes.NewBattleScene(g.world, g.width, g.height, g.speed, g.gameStats, multiplayer, gameOptions, g.startingTowerLevel, g.switchToTitle)
	if err != nil {
		return err
	}
	battle.Init()

	g.scenes = []Scene{battle}
	if multiplayer {
		scene, err := scenes.NewViewerScene(clientWorld, g.width, g.height, gameOptions, true)
		if err != nil {
			return err
		}
		ebiten.SetWindowSize(g.width*2, g.height)
		g.adjustWindowPosition()
		g.scenes = append(g.scenes, scene)
	} else {
		ebiten.SetWindowSize(g.width, g.height)
	}
	return nil
}

func (g *GameData) adjustWindowPosition() {
	monWidth, _ := ebiten.Monitor().Size()
	winX, winY := ebiten.WindowPosition()
	winWidth, _ := ebiten.WindowSize()
	if winX+winWidth > monWidth {
		winX = monWidth - winWidth
	}

	ebiten.SetWindowPosition(winX, winY)
}

func (g *GameData) switchToTitle(gameStats *comp.GameStats, gameOptions *config.ConfigData) error {
	if gameStats != g.gameStats {
		if gameStats.HighScore > g.gameStats.HighScore {
			g.gameStats.HighScore = gameStats.HighScore
		}
		if gameStats.HighCreepLevel > g.gameStats.HighCreepLevel {
			g.gameStats.HighCreepLevel = gameStats.HighCreepLevel
		}
		if gameStats.HighTowerLevel > g.gameStats.HighTowerLevel {
			g.gameStats.HighTowerLevel = gameStats.HighTowerLevel
		}
		g.gameStats.Update(gameStats)
		g.gameStats.SaveStats()
	}

	title, err := scenes.NewTitleScene(g.world, g.width, g.height, g.gameStats, gameOptions, g.switchToBattle)
	if err != nil {
		return err
	}

	ebiten.SetWindowSize(g.width, g.height)
	g.scenes = []Scene{title}
	return nil
}

func (g *GameData) Update() error {
	// TODO move this into the title scene so we can determine if the config has focus
	if !scenes.IsModalOpen() && inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if err := g.gameStats.SaveStats(); err != nil {
			return err
		}

		return ebiten.Termination
	}

	for _, scene := range g.scenes {
		if err := scene.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (g *GameData) Draw(screen *ebiten.Image) {
	for _, scene := range g.scenes {
		scene.Draw(screen)
	}
}

func (g *GameData) Layout(width, height int) (int, int) {
	return width, height
}
