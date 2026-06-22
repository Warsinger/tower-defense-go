package components

import (
	"testing"
	"tower-defense/util"

	"github.com/yohamta/donburi"
)

func newTowerTestEntry(t *testing.T, playerTowerLevels int, towerLevel int) *donburi.Entry {
	t.Helper()

	world := donburi.NewWorld()
	playerEntity := world.Create(Player)
	Player.Set(world.Entry(playerEntity), &PlayerData{Money: 500, TowerLevels: playerTowerLevels})

	towerEntity := world.Create(Tower, Health, Attack, Level)
	towerEntry := world.Entry(towerEntity)
	Tower.Set(towerEntry, &TowerData{})
	Health.Set(towerEntry, &HealthData{Health: 10, MaxHealth: 20})
	Attack.Set(towerEntry, &AttackData{Power: 1, AttackType: RangedSingle, Range: 50, cooldown: util.NewCooldownTimer(30)})
	Level.Set(towerEntry, &LevelData{Level: towerLevel})

	SetGameStats(NewGameStats(nil))
	t.Cleanup(func() { SetGameStats(nil) })

	return towerEntry
}

func TestTowerData_HealRestoresHealthAndTracksStat(t *testing.T) {
	entry := newTowerTestEntry(t, 0, 1)
	tower := Tower.Get(entry)

	if !tower.Heal(entry, false) {
		t.Fatal("Heal() = false, want true for damaged tower")
	}

	health := Health.Get(entry)
	if health.Health != health.MaxHealth {
		t.Fatalf("health after Heal() = %v, want max health %v", health.Health, health.MaxHealth)
	}
	if got := GetGameStats().GetStat("TowersHealed"); got != 1 {
		t.Errorf("TowersHealed = %v, want 1", got)
	}

	if tower.Heal(entry, false) {
		t.Fatal("Heal() = true, want false when already at max health")
	}
	if got := GetGameStats().GetStat("TowersHealed"); got != 1 {
		t.Errorf("TowersHealed after no-op heal = %v, want 1", got)
	}
}

func TestTowerData_UpgradeAppliesScalingAndTracksStat(t *testing.T) {
	entry := newTowerTestEntry(t, 20, 2)
	tower := Tower.Get(entry)

	if !tower.Upgrade(entry, false) {
		t.Fatal("Upgrade() = false, want true below max level")
	}

	level := Level.Get(entry)
	if level.Level != 3 {
		t.Errorf("level after Upgrade() = %v, want 3", level.Level)
	}

	health := Health.Get(entry)
	if health.MaxHealth != 25 || health.Health != 25 {
		t.Errorf("health after Upgrade() = %v/%v, want 25/25", health.Health, health.MaxHealth)
	}

	attack := Attack.Get(entry)
	if attack.Power != 2 {
		t.Errorf("attack power after Upgrade() = %v, want 2", attack.Power)
	}
	if attack.Range != 53 {
		t.Errorf("attack range after Upgrade() = %v, want 53", attack.Range)
	}
	if attack.cooldown.Cooldown != 27 {
		t.Errorf("attack cooldown after Upgrade() = %v, want 27", attack.cooldown.Cooldown)
	}
	if got := GetGameStats().GetStat("TowersUpgraded"); got != 1 {
		t.Errorf("TowersUpgraded = %v, want 1", got)
	}
}

func TestTowerData_UpgradeStopsAtPlayerMaxTowerLevel(t *testing.T) {
	entry := newTowerTestEntry(t, 0, 5)
	tower := Tower.Get(entry)

	if tower.Upgrade(entry, false) {
		t.Fatal("Upgrade() = true, want false at max tower level")
	}
	if got := Level.Get(entry).Level; got != 5 {
		t.Errorf("level after max-level Upgrade() = %v, want 5", got)
	}
	if got := GetGameStats().GetStat("TowersUpgraded"); got != 0 {
		t.Errorf("TowersUpgraded after max-level Upgrade() = %v, want 0", got)
	}
}

func TestAfterTowerAttackConsumesAmmoAndRemovesEmptyTower(t *testing.T) {
	entry := newTowerTestEntry(t, 0, 1)
	health := Health.Get(entry)
	health.Health = 2

	AfterTowerAttack(entry)
	if got := Health.Get(entry).Health; got != 1 {
		t.Fatalf("health after first attack = %v, want 1", got)
	}
	if !entry.Valid() {
		t.Fatal("tower removed after first attack, want still valid")
	}

	AfterTowerAttack(entry)
	if entry.Valid() {
		t.Fatal("tower still valid after ammo reaches zero, want removed")
	}
	if got := GetGameStats().GetStat("TowersAmmoOut"); got != 1 {
		t.Errorf("TowersAmmoOut = %v, want 1", got)
	}
}
