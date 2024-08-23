package components

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

type GameStats struct {
	stats map[string]int
	/*
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
	*/
	StartTime time.Time
	GameTime  time.Duration
}

const statsFile = "score/stats.txt"

var (
	gameStats  *GameStats
	validStats = []string{
		"BulletsExpired",
		"CreepBulletsFired",
		"CreepsKilled",
		"CreepsSpawned",
		"CreepWaves",
		"Games",
		"HighCreepLevel",
		"HighScore",
		"HighTowerLevel",
		"MoneySpent",
		"PlayerDeaths",
		"TowerBulletsFired",
		"TowersAmmoOut",
		"TowersBuilt",
		"TowersHealed",
		"TowersKilled",
		"TowersUpgraded",
	}
	displayNames = makeDisplayNames(validStats)
)

func SetGameStats(gs *GameStats) {
	gameStats = gs
}
func GetGameStats() *GameStats {
	return gameStats
}
func NewGameStats(old *GameStats) *GameStats {
	gs := &GameStats{stats: make(map[string]int, 16), StartTime: time.Now()}
	if old != nil {
		gs.stats["HighScore"] = old.stats["HighScore"]
		gs.stats["HighCreepLevel"] = old.stats["HighCreepLevel"]
		gs.stats["HighTowerLevel"] = old.stats["HighTowerLevel"]
	}
	return gs
}
func (gs *GameStats) GetStat(name string) int {
	return gs.stats[name]
}

func (gs *GameStats) UpdateHighs(score, creepLevel, maxTowerLevel int) {
	gs.stats["HighScore"] = max(score, gs.stats["HighScore"])
	gs.stats["HighCreepLevel"] = max(creepLevel, gs.stats["HighCreepLevel"])
	gs.stats["HighTowerLevel"] = max(maxTowerLevel, gs.stats["HighTowerLevel"])

}
func (gs *GameStats) Update(other *GameStats) {
	gs.UpdateHighs(other.stats["HighScore"], other.stats["HighCreepLevel"], other.stats["HighTowerLevel"])
	gs.stats["Games"]++
	gs.UpdateStats(other, "Game", "High")
	gs.GameTime += other.GameTime
}
func (gs *GameStats) UpdateStats(other *GameStats, excludePrefixes ...string) {
	for _, name := range slices.Sorted(maps.Keys(gs.stats)) {
		var exclude bool = false
		for _, prefix := range excludePrefixes {
			if strings.HasPrefix(name, prefix) {
				exclude = true
			}
		}
		if !exclude {
			gs.stats[name] += other.stats[name]
		}
	}
}

func (gs *GameStats) UpdateStat(name string, count int) {
	gs.stats[name] += count
}
func (gs *GameStats) IncrementStat(name string) {
	gs.stats[name]++
}

func (gs *GameStats) RunningTime() time.Duration {
	return time.Since(gs.StartTime)
}
func (gs *GameStats) FinalizeTime() {
	gs.GameTime = time.Since(gs.StartTime)
}

func (gs *GameStats) Reset() {
	gs.initStats("High", "Game")
	gs.StartTime = time.Now()
	gs.GameTime = 0
}

func LoadStats() *GameStats {
	gameStats := NewGameStats(nil)
	bytes, err := os.ReadFile(statsFile)
	if err == nil {
		strings.Split(string(bytes), "\n")
		for _, line := range strings.Split(string(bytes), "\n") {
			values := strings.Split(line, "=")
			if len(values) != 2 {
				continue
			}
			name, value := values[0], values[1]
			if name == "GameTime" {
				gameStats.GameTime, err = time.ParseDuration(value)
				if err != nil {
					fmt.Printf("WARN %s formatting err %s %v\n", name, value, err)
				}
			} else if slices.Contains(validStats, name) {
				gameStats.stats[name] = parseScore(name, value)
			} else {
				fmt.Printf("Invalid stat loaded %s\n", line)
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

	return os.WriteFile(statsFile, []byte(gs.StatsLines("=", false, false)), 0644)
}

func (gs *GameStats) initStats(excludePrefixes ...string) {
	for _, name := range slices.Sorted(maps.Keys(gs.stats)) {
		var exclude bool = false
		for _, prefix := range excludePrefixes {
			if strings.HasPrefix(name, prefix) {
				exclude = true
			}
		}
		if !exclude {
			gs.stats[name] = 0
		}
	}
}

func (gs *GameStats) StatsLines(delim string, runningTime, forDisplay bool, excludePrefixes ...string) string {
	var b strings.Builder

	for _, name := range slices.Sorted(maps.Keys(gs.stats)) {
		var exclude bool = false
		for _, prefix := range excludePrefixes {
			if strings.HasPrefix(name, prefix) {
				exclude = true
			}
		}
		if !exclude {
			displayName := name
			if forDisplay {
				displayName = displayNames[name]
			}
			fmt.Fprintf(&b, "%s%s%d\n", displayName, delim, gs.stats[name])
		}
	}
	var duration time.Duration
	if runningTime {
		duration = gs.RunningTime()
	} else {
		duration = gs.GameTime
	}
	fmt.Fprintf(&b, "GameTime%s%v\n", delim, duration.Round(time.Second))
	return b.String()
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

func makeDisplayNames(validStats []string) map[string]string {
	displayNames := make(map[string]string, len(validStats))
	for _, name := range validStats {
		displayNames[name] = makeDisplayName(name)
	}
	return displayNames
}

func makeDisplayName(name string) string {
	// Regular expression to match word boundaries in camel case strings
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)

	// Replace the matched boundaries with a space and the matched characters
	return re.ReplaceAllString(name, "${1} ${2}")

}
