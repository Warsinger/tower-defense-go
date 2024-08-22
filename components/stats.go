package components

import "time"

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
	gs.StartTime = time.Now()
	gs.GameTime = 0
}
