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
	if other.stats["HighScore"] > gs.stats["HighScore"] {
		gs.stats["HighScore"] = other.stats["HighScore"]
	}
	if other.stats["HighCreepLevel"] > gs.stats["HighCreepLevel"] {
		gs.stats["HighCreepLevel"] = other.stats["HighCreepLevel"]
	}
	if other.stats["HighTowerLevel"] > gs.stats["HighTowerLevel"] {
		gs.stats["HighTowerLevel"] = other.stats["HighTowerLevel"]
	}
	gs.stats["BulletsExpired"] += other.stats["BulletsExpired"]
	gs.stats["CreepBulletsFired"] += other.stats["CreepBulletsFired"]
	gs.stats["CreepWaves"] += other.stats["CreepWaves"]
	gs.stats["CreepsKilled"] += other.stats["CreepsKilled"]
	gs.stats["CreepsSpawned"] += other.stats["CreepsSpawned"]
	gs.stats["MoneySpent"] += other.stats["MoneySpent"]
	gs.stats["PlayerDeaths"] += other.stats["PlayerDeaths"]
	gs.stats["TowerBulletsFired"] += other.stats["TowerBulletsFired"]
	gs.stats["TowersAmmoOut"] += other.stats["TowersAmmoOut"]
	gs.stats["TowersBuilt"] += other.stats["TowersBuilt"]
	gs.stats["TowersHealed"] += other.stats["TowersHealed"]
	gs.stats["TowersKilled"] += other.stats["TowersKilled"]
	gs.stats["TowersUpgraded"] += other.stats["TowersUpgraded"]
	gs.GameTime += other.GameTime
}

func (gs *GameStats) UpdateCreepsSpawned(count int) {
	gs.stats["CreepsSpawned"] += count
}
func (gs *GameStats) UpdateCreepsKilled() {
	gs.stats["CreepsKilled"]++
}
func (gs *GameStats) UpdateCreepWaves() {
	gs.stats["CreepWaves"]++
}
func (gs *GameStats) UpdateTowersBuilt() {
	gs.stats["TowersBuilt"]++
}
func (gs *GameStats) UpdateTowersHealed() {
	gs.stats["TowersHealed"]++
}
func (gs *GameStats) UpdateTowersUpgraded() {
	gs.stats["TowersUpgraded"]++
}
func (gs *GameStats) UpdateTowersKilled() {
	gs.stats["TowersKilled"]++
}
func (gs *GameStats) UpdateTowersAmmoOut() {
	gs.stats["TowersAmmoOut"]++
}
func (gs *GameStats) UpdateTowerBulletsFired() {
	gs.stats["TowerBulletsFired"]++
}
func (gs *GameStats) UpdateCreepBulletsFired() {
	gs.stats["CreepBulletsFired"]++
}
func (gs *GameStats) UpdateBulletsExpired() {
	gs.stats["BulletsExpired"]++
}
func (gs *GameStats) UpdatePlayerDeaths() {
	gs.stats["PlayerDeaths"]++
}
func (gs *GameStats) UpdateMoneySpent(money int) {
	gs.stats["MoneySpent"] += money
}

func (gs *GameStats) RunningTime() time.Duration {
	return time.Since(gs.StartTime)
}
func (gs *GameStats) FinalizeTime() {
	gs.GameTime = time.Since(gs.StartTime)
}

func (gs *GameStats) Reset() {
	gs.stats["BulletsExpired"] = 0
	gs.stats["CreepBulletsFired"] = 0
	gs.stats["CreepsKilled"] = 0
	gs.stats["CreepsSpawned"] = 0
	gs.stats["CreepWaves"] = 0
	gs.stats["MoneySpent"] = 0
	gs.stats["PlayerDeaths"] = 0
	gs.stats["TowerBulletsFired"] = 0
	gs.stats["TowersBuilt"] = 0
	gs.stats["TowersHealed"] = 0
	gs.stats["TowersKilled"] = 0
	gs.stats["TowersUpgraded"] = 0
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
