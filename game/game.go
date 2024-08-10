package game

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	gameStats            *scenes.GameStats
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

	gameStats := LoadScores()

	game := &GameData{world: donburi.NewWorld(), width: width, height: height, speed: speed, gameStats: gameStats, debug: debug}

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
		err = game.switchToTitle(gameStats)
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
	battle, err := scenes.NewBattleScene(g.world, g.width, g.height, g.speed, g.gameStats, g.debug, g.switchToTitle)
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

func (g *GameData) switchToTitle(gameStats *scenes.GameStats) error {
	if gameStats != g.gameStats {
		save := false
		if gameStats.HighScore > g.gameStats.HighScore {
			g.gameStats.HighScore = gameStats.HighScore
			save = true
		}
		if gameStats.HighCreepLevel > g.gameStats.HighCreepLevel {
			g.gameStats.HighCreepLevel = gameStats.HighCreepLevel
			save = true
		}
		if gameStats.HighTowerLevel > g.gameStats.HighTowerLevel {
			g.gameStats.HighTowerLevel = gameStats.HighTowerLevel
			save = true
		}
		if save {
			g.SaveScores()
		}
	}

	title, err := scenes.NewTitleScene(g.width, g.height, g.gameStats, g.switchToBattle)
	if err != nil {
		return err
	}

	ebiten.SetWindowSize(g.width, g.height)
	g.scenes = []Scene{title}
	return nil
}

const statsFile = "score/stats.txt"

func LoadScores() *scenes.GameStats {
	var highScore, highCreepLevel, highTowerLevel int
	bytes, err := os.ReadFile(statsFile)
	if err == nil {
		strings.Split(string(bytes), "\n")
		for _, line := range strings.Split(string(bytes), "\n") {
			values := strings.Split(line, "=")
			if len(values) != 2 {
				continue
			}
			switch values[0] {
			case "score":
				highScore, err = strconv.Atoi(values[1])
				if err != nil {
					fmt.Printf("WARN high score formatting err %v\n", err)
				}
			case "creepLevel":
				highCreepLevel, err = strconv.Atoi(values[1])
				if err != nil {
					fmt.Printf("WARN high creep level formatting err %v\n", err)
				}
			case "towerLevel":
				highTowerLevel, err = strconv.Atoi(values[1])
				if err != nil {
					fmt.Printf("WARN high tower level formatting err %v\n", err)
				}
			}
		}
	}

	return &scenes.GameStats{HighScore: highScore, HighCreepLevel: highCreepLevel, HighTowerLevel: highTowerLevel}
}
func (g *GameData) SaveScores() error {
	dir, _ := filepath.Split(statsFile)

	if err := ensureDir(dir); err != nil {
		return err
	}

	str := fmt.Sprintf("score=%d\ncreepLevel=%d\ntowerLevel=%d\n", g.gameStats.HighScore, g.gameStats.HighCreepLevel, g.gameStats.HighTowerLevel)
	return os.WriteFile(statsFile, []byte(str), 0644)
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
