package game

import (
	"fmt"
	"os"
	"strconv"
	"tower-defense/assets"
	"tower-defense/scenes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}
type GameData struct {
	scene                Scene
	width, height, speed int
	highScore            int
	debug                bool
}

func NewGame(width, height, speed int, debug bool) (*GameData, error) {
	err := assets.LoadAssets()
	if err != nil {
		return nil, err
	}

	ebiten.SetWindowSize(int(width), int(height))
	ebiten.SetWindowTitle("Tower Defense")

	highScore := LoadScores()
	fmt.Printf("hs load %v\n", highScore)

	game := &GameData{width: width, height: height, speed: speed, highScore: highScore, debug: debug}

	err = game.switchToTitle(highScore)
	if err != nil {
		return nil, err
	}
	return game, nil
}
func (g *GameData) switchToBattle() error {
	battle, err := scenes.NewBattleScene(g.width, g.height, g.speed, g.highScore, g.debug, g.switchToTitle)
	if err != nil {
		return err
	}
	battle.Init()
	g.scene = battle
	return nil
}

func (g *GameData) switchToTitle(score int) error {
	if score > g.highScore {
		g.highScore = score
		g.SaveScores()
	}
	title, err := scenes.NewTitleScene(g.width, g.height, g.highScore, g.switchToBattle)
	if err != nil {
		return err
	}

	g.scene = title
	return nil
}

const highScoreFile = "score/highscore.txt"

func LoadScores() int {
	var highScore int = 0
	bytes, err := os.ReadFile(highScoreFile)
	if err == nil {
		highScore, err = strconv.Atoi(string(bytes))
		if err != nil {
			fmt.Printf("WARN high score formatting err %v\n", err)
		}
	}

	return highScore
}
func (g *GameData) SaveScores() error {
	str := strconv.Itoa(g.highScore)

	err := os.WriteFile(highScoreFile, []byte(str), 0644)
	return err
}

func (g *GameData) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if err := g.SaveScores(); err != nil {
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
