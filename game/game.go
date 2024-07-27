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
	scenes               []Scene
	width, height, speed int
	highScore            int
	debug                bool
}

func NewGame(width, height, speed int, debug bool) (*GameData, error) {
	err := assets.LoadAssets()
	if err != nil {
		return nil, err
	}

	ebiten.SetWindowTitle("Tower Defense")

	highScore := LoadScores()

	game := &GameData{width: width, height: height, speed: speed, highScore: highScore, debug: debug}

	err = game.switchToTitle(highScore)
	if err != nil {
		return nil, err
	}
	return game, nil
}
func (g *GameData) switchToBattle(viewerMode bool) error {
	battle, err := scenes.NewBattleScene(g.width, g.height, g.speed, g.highScore, g.debug, g.switchToTitle)
	if err != nil {
		return err
	}
	battle.Init()
	g.scenes = []Scene{battle}

	if viewerMode {
		viewer, err := scenes.NewViewerScene(battle.World(), g.width, g.height)
		if err != nil {
			return err
		}
		ebiten.SetWindowSize(g.width*2, g.height)
		g.adjustWindowPosition()
		g.scenes = append(g.scenes, viewer)
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

func (g *GameData) switchToTitle(score int) error {
	if score > g.highScore {
		g.highScore = score
		g.SaveScores()
	}
	title, err := scenes.NewTitleScene(g.width, g.height, g.highScore, g.switchToBattle)
	if err != nil {
		return err
	}

	ebiten.SetWindowSize(g.width, g.height)
	g.scenes = []Scene{title}
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
