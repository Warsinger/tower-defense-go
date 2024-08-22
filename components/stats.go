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

	BulletsExpired    int
	CreepBulletsFired int
	CreepsKilled      int
	CreepsSpawned     int
	CreepWaves        int
	MoneySpent        int
	PlayerDeaths      int
	TowerBulletsFired int
	TowersAmmoOut     int
	TowersBuilt       int
	TowersHealed      int
	TowersKilled      int
	TowersUpgraded    int
	StartTime         time.Time
	GameTime          time.Duration
}

const statsFile = "score/stats.txt"

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

func (gs *GameStats) Update(other *GameStats) {
	if other.HighScore > gs.HighScore {
		gs.HighScore = other.HighScore
	}
	if other.HighCreepLevel > gs.HighCreepLevel {
		gs.HighCreepLevel = other.HighCreepLevel
	}
	if other.HighTowerLevel > gs.HighTowerLevel {
		gs.HighTowerLevel = other.HighTowerLevel
	}
	gs.BulletsExpired = other.BulletsExpired
	gs.CreepBulletsFired = other.CreepBulletsFired
	gs.CreepWaves = other.CreepWaves
	gs.CreepsKilled = other.CreepsKilled
	gs.CreepsSpawned = other.CreepsSpawned
	gs.GameTime += other.CalcDuration()
	gs.MoneySpent = other.MoneySpent
	gs.PlayerDeaths = other.PlayerDeaths
	gs.TowerBulletsFired = other.TowerBulletsFired
	gs.TowersAmmoOut = other.TowersAmmoOut
	gs.TowersBuilt = other.TowersBuilt
	gs.TowersHealed = other.TowersHealed
	gs.TowersKilled = other.TowersKilled
	gs.TowersUpgraded = other.TowersUpgraded
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
func (gs *GameStats) UpdateTowersHealed() {
	gs.TowersHealed++
}
func (gs *GameStats) UpdateTowersUpgraded() {
	gs.TowersUpgraded++
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
	gs.TowersHealed = 0
	gs.TowersUpgraded = 0
	gs.StartTime = time.Now()
	gs.GameTime = 0
}

func LoadStats() *GameStats {
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
			case "towersHealed":
				gameStats.TowersHealed = parseScore(values[0], values[1])
			case "towersUpgraded":
				gameStats.TowersUpgraded = parseScore(values[0], values[1])
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

func (gs *GameStats) SaveStats() error {
	dir, _ := filepath.Split(statsFile)

	if err := ensureDir(dir); err != nil {
		return err
	}

	str := fmt.Sprintf("score=%d\ncreepLevel=%d\ntowerLevel=%d\ncreepsSpawned=%d\ncreepsKilled=%d\ncreepWaves=%d\ntowersBuilt=%d\ntowersKilled=%d\ntowersAmmoOut=%d\ntowerBulletsFired=%d\ncreepbulletsFired=%d\nbulletsExpired=%d\nplayerDeaths%d\ntowersHealed=%d\ntowersUpgraded=%d\ngameTime=%v\n",
		gs.HighScore, gs.HighCreepLevel, gs.HighTowerLevel,
		gs.CreepsSpawned, gs.CreepsKilled, gs.CreepWaves, gs.TowersBuilt, gs.TowersKilled, gs.TowersAmmoOut, gs.TowerBulletsFired, gs.CreepBulletsFired, gs.BulletsExpired, gs.PlayerDeaths, gs.TowersHealed, gs.TowersUpgraded, gs.CalcDuration())
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
