package components

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type GameStats struct {
	HighScore      int
	HighCreepLevel int
	HighTowerLevel int

	CreepsSpawned     int
	CreepsKilled      int
	CreepWaves        int
	TowersBuilt       int
	TowersKilled      int
	TowersAmmoOut     int
	TowerBulletsFired int
	CreepBulletsFired int
	BulletsExpired    int
	PlayerDeaths      int
	MoneySpent        int
	StartTime         time.Time
	GameTime          time.Duration
}

var gameStats *GameStats

func SetGameStats(gs *GameStats) {
	gameStats = gs
}
func GetGameStats() *GameStats {
	return gameStats
}
func NewGameStats(highScore, highCreepLevel, highTowerLevel int) *GameStats {
	return &GameStats{HighScore: highScore, HighCreepLevel: highCreepLevel, HighTowerLevel: highTowerLevel}
}

func (gs *GameStats) UpdateCreepsSpawned(count int) {
	gs.CreepsSpawned += count
}
func (gs *GameStats) UpdateCreepsKilled() {
	gs.CreepsKilled++
}
func (gs *GameStats) UpdateCreepWaves() {
	gs.CreepWaves++
}
func (gs *GameStats) UpdateTowersBuilt() {
	gs.TowersBuilt++
}
func (gs *GameStats) UpdateTowersKilled() {
	gs.TowersKilled++
}
func (gs *GameStats) UpdateTowersAmmoOut() {
	gs.TowersAmmoOut++
}
func (gs *GameStats) UpdateTowerBulletsFired() {
	gs.TowerBulletsFired++
}
func (gs *GameStats) UpdateCreepBulletsFired() {
	gs.CreepBulletsFired++
}
func (gs *GameStats) UpdateBulletsExpired() {
	gs.BulletsExpired++
}
func (gs *GameStats) UpdatePlayerDeaths() {
	gs.PlayerDeaths++
}
func (gs *GameStats) UpdateMoneyspent(money int) {
	gs.CreepsSpawned += money
}

func (gs *GameStats) CalcDuration() time.Duration {
	return time.Since(gs.StartTime)
}

func (gs *GameStats) Reset() {
	gs.CreepsSpawned = 0
	gs.CreepsKilled = 0
	gs.CreepWaves = 0
	gs.TowersBuilt = 0
	gs.TowersKilled = 0
	gs.TowerBulletsFired = 0
	gs.CreepBulletsFired = 0
	gs.BulletsExpired = 0
	gs.PlayerDeaths = 0
	gs.MoneySpent = 0
	gs.StartTime = time.Now()
	gs.GameTime = 0
}

const statsFile = "score/stats.txt"

func LoadScores() *GameStats {
	gameStats := &GameStats{}
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

func (gs *GameStats) SaveScores() error {
	dir, _ := filepath.Split(statsFile)

	if err := ensureDir(dir); err != nil {
		return err
	}

	str := fmt.Sprintf("score=%d\ncreepLevel=%d\ntowerLevel=%d\ncreepsSpawned=%d\ncreepsKilled=%d\ncreepWaves=%d\ntowersBuilt=%d\ntowersKilled=%d\ntowersAmmoOut=%d\ntowerBulletsFired=%d\ncreepbulletsFired=%d\nbulletsExpired=%d\nplayerDeaths%d\ngameTime=%v\n",
		gs.HighScore, gs.HighCreepLevel, gs.HighTowerLevel,
		gs.CreepsSpawned, gs.CreepsKilled, gs.CreepWaves, gs.TowersBuilt, gs.TowersKilled, gs.TowersAmmoOut, gs.TowerBulletsFired, gs.CreepBulletsFired, gs.BulletsExpired, gs.PlayerDeaths, gs.CalcDuration())
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
