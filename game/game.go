package game

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"tower-defense/assets"
	"tower-defense/network"
	"tower-defense/scenes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	highScore            int
	debug                bool
	server               *network.Server
	client               *network.Client
}

func NewGame(width, height, speed int, debug bool, server, client string) (*GameData, error) {
	err := assets.LoadAssets()
	if err != nil {
		return nil, err
	}

	ebiten.SetWindowTitle("Tower Defense")

	highScore := LoadScores()

	game := &GameData{world: donburi.NewWorld(), width: width, height: height, speed: speed, highScore: highScore, debug: debug}

	if server != "" {
		fmt.Printf("listening on port %v\n", server)
		game.server = network.NewServer(game.world, "", server)
		err = game.server.Start()
		if err != nil {
			return nil, err
		}
	}

	if client != "" {
		fmt.Printf("connect to %v\n", client)
		game.client = network.NewClientNewWorld(client)
		err = game.client.Start()
		if err != nil {
			return nil, err
		}
	}

	if game.client != nil {
		err = game.switchToViewer()
	} else {
		err = game.switchToTitle(highScore)
	}
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (g *GameData) switchToViewer() error {
	scene, err := scenes.NewViewerScene(g.client.World, g.width, g.height, g.debug, g.client == nil)
	if err != nil {
		return err
	}
	ebiten.SetWindowSize(g.width, g.height)
	g.scenes = []Scene{scene}
	return nil
}

func (g *GameData) switchToBattle(viewerMode bool) error {
	battle, err := scenes.NewBattleScene(g.world, g.width, g.height, g.speed, g.highScore, g.debug, g.switchToTitle)
	if err != nil {
		return err
	}
	battle.Init()
	g.scenes = []Scene{battle}

	if viewerMode || g.client != nil {
		var world donburi.World
		if g.client != nil {
			world = g.client.World
		} else {
			world = g.world
		}
		scene, err := scenes.NewViewerScene(world, g.width, g.height, g.debug, g.client == nil)
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

	dir, _ := filepath.Split(highScoreFile)

	if err := ensureDir(dir); err != nil {
		return err
	}
	return os.WriteFile(highScoreFile, []byte(str), 0644)
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
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
