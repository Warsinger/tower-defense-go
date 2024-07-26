package game

import (
	"tower-defense/assets"
	"tower-defense/scenes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
	End() error
}
type GameData struct {
	scene Scene
}

func NewGame(width, height, speed int, debug bool) (*GameData, error) {
	err := assets.LoadAssets()
	if err != nil {
		return nil, err
	}

	ebiten.SetWindowSize(int(width), int(height))
	ebiten.SetWindowTitle("Tower Defense")

	scene, err := scenes.NewBattleScene(width, height, speed, debug)
	if err != nil {
		return nil, err
	}
	err = scene.Init()
	if err != nil {
		return nil, err
	}
	game := &GameData{scene: scene}
	return game, nil
}

func (g *GameData) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if err := g.scene.End(); err != nil {
			return err
		}
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	return g.scene.Update()
}

func (g *GameData) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

func (g *GameData) Layout(width, height int) (int, int) {
	return width, height
}
