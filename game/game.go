package game

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

	gameStats := LoadScores()

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
		g.SaveScores()
	}

	title, err := scenes.NewTitleScene(g.world, g.width, g.height, g.gameStats, gameOptions, g.switchToBattle)
	if err != nil {
		return err
	}

	ebiten.SetWindowSize(g.width, g.height)
	g.scenes = []Scene{title}
	return nil
}

const statsFile = "score/stats.txt"

func LoadScores() *comp.GameStats {
	gameStats := &comp.GameStats{}
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
				gameStats.HighScore = parseScore(values[0], values[1])
			case "creepLevel":
				gameStats.HighCreepLevel = parseScore(values[0], values[1])
			case "towerLevel":
				gameStats.HighTowerLevel = parseScore(values[0], values[1])
			case "creepsSpawned":
				gameStats.CreepsSpawned = parseScore(values[0], values[1])
			case "creepsKilled":
				gameStats.CreepsKilled = parseScore(values[0], values[1])
			case "creepWaves":
				gameStats.CreepWaves = parseScore(values[0], values[1])
			case "towersBuilt":
				gameStats.TowersBuilt = parseScore(values[0], values[1])
			case "towersKilled":
				gameStats.TowersKilled = parseScore(values[0], values[1])
			case "towerBulletsFired":
				gameStats.TowerBulletsFired = parseScore(values[0], values[1])
			case "creepBulletsFired":
				gameStats.CreepBulletsFired = parseScore(values[0], values[1])
			case "bulletsExpired":
				gameStats.BulletsExpired = parseScore(values[0], values[1])
			case "playerDeaths":
				gameStats.PlayerDeaths = parseScore(values[0], values[1])
			case "gameTime":
				gameStats.GameTime, err = time.ParseDuration(values[1])
				if err != nil {
					fmt.Printf("WARN %s formatting err %s %v\n", values[0], values[1], err)
				}
			}
		}
	}

	return gameStats
}

func parseScore(label, val string) int {
	score, err := strconv.Atoi(val)
	if err != nil {
		fmt.Printf("WARN %s formatting err %s %v\n", label, val, err)
	}
	return score
}

func (g *GameData) SaveScores() error {
	dir, _ := filepath.Split(statsFile)

	if err := ensureDir(dir); err != nil {
		return err
	}

	str := fmt.Sprintf("score=%d\ncreepLevel=%d\ntowerLevel=%d\ncreepsSpawned=%d\ncreepsKilled=%d\ncreepWaves=%d\ntowersBuilt=%d\ntowersKilled=%d\ntowersAmmoOut=%d\ntowerBulletsFired=%d\ncreepbulletsFired=%d\nbulletsExpired=%d\nplayerDeaths%d\ngameTime=%v\n",
		g.gameStats.HighScore, g.gameStats.HighCreepLevel, g.gameStats.HighTowerLevel,
		g.gameStats.CreepsSpawned, g.gameStats.CreepsKilled, g.gameStats.CreepWaves, g.gameStats.TowersBuilt, g.gameStats.TowersKilled, g.gameStats.TowersAmmoOut, g.gameStats.TowerBulletsFired, g.gameStats.CreepBulletsFired, g.gameStats.BulletsExpired, g.gameStats.PlayerDeaths, g.gameStats.CalcDuration())
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
	// TODO move this into the title scene so we can determine if the config has focus
	if !scenes.IsModalOpen() && inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if err := g.SaveScores(); err != nil {
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
