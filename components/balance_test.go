package components

import (
	"testing"
	"tower-defense/config"
	"tower-defense/util"

	"github.com/yohamta/donburi"
)

func testBalance() *config.BalanceData {
	balance := *config.DefaultBalance()
	balance.Player.StartingMoney = 321
	balance.Player.Health = 222
	balance.Player.AttackPower = 7
	balance.Player.AttackRange = 44
	balance.Player.AttackCooldown = 13
	balance.Player.MaxTowerInitialLevel = 9
	balance.Player.MaxTowerLevelsPerBonus = 10
	balance.Player.CreepLevelTowerLevels = 4

	balance.Tower.DefaultType = "Laser"
	balance.Tower.Costs = map[string]int{"Laser": 80}
	balance.Tower.HealCostDivisor = 4
	balance.Tower.Health = 55
	balance.Tower.AttackPower = 6
	balance.Tower.AttackRange = 70
	balance.Tower.AttackCooldown = 21
	balance.Tower.InitialLevel = 3
	balance.Tower.UpgradeMaxHealthAdd = 11
	balance.Tower.UpgradePowerLevelDivisor = 2
	balance.Tower.UpgradeRangeAdd = 8
	balance.Tower.UpgradeCooldownReduction = 5
	balance.Tower.UpgradeMinCooldown = 4

	balance.Creep.BigCreepChance = 1
	balance.Creep.BigCreepChoice = 9
	balance.Creep.BigCreepAugment = 3
	balance.Creep.BaseVelocityY = 10
	balance.Creep.HealthBase = 5
	balance.Creep.HealthAugmentMultiplier = 4
	balance.Creep.HealthLevelDivisor = 2
	balance.Creep.AttackPowerBase = 2
	balance.Creep.AttackPowerLevelOffset = 1
	balance.Creep.AttackPowerLevelDivisor = 3
	balance.Creep.AttackRangeBase = 12
	balance.Creep.AttackRangeAugmentMultiplier = 6
	balance.Creep.AttackCooldownBase = 8
	balance.Creep.AttackCooldownAugmentMultiplier = 7
	balance.Creep.ScoreValueBase = 15

	balance.SuperCreep.VelocityX = 4
	balance.SuperCreep.VelocityY = 6
	balance.SuperCreep.ScoreValue = 99
	balance.SuperCreep.Health = 88
	balance.SuperCreep.AttackPower = 12
	balance.SuperCreep.AttackRange = 34
	balance.SuperCreep.AttackCooldown = 56

	return &balance
}

func TestNewPlayerUsesBalanceValues(t *testing.T) {
	world := donburi.NewWorld()
	balance := testBalance()
	config.NewBalance(world, balance)
	if _, err := NewBoard(world, 200, 300); err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	if err := NewPlayer(world, 7); err != nil {
		t.Fatalf("NewPlayer() error = %v", err)
	}

	entry := Player.MustFirst(world)
	player := Player.Get(entry)
	health := Health.Get(entry)
	attack := Attack.Get(entry)

	if player.Money != balance.Player.StartingMoney {
		t.Errorf("player money = %v, want %v", player.Money, balance.Player.StartingMoney)
	}
	if player.TowerLevels != 7 {
		t.Errorf("tower levels = %v, want 7", player.TowerLevels)
	}
	if health.Health != balance.Player.Health || health.MaxHealth != balance.Player.Health {
		t.Errorf("health = %v/%v, want %v/%v", health.Health, health.MaxHealth, balance.Player.Health, balance.Player.Health)
	}
	if attack.Power != balance.Player.AttackPower {
		t.Errorf("attack power = %v, want %v", attack.Power, balance.Player.AttackPower)
	}
	if attack.Range != balance.Player.AttackRange {
		t.Errorf("attack range = %v, want %v", attack.Range, balance.Player.AttackRange)
	}
	if attack.cooldown.Cooldown != balance.Player.AttackCooldown {
		t.Errorf("attack cooldown = %v, want %v", attack.cooldown.Cooldown, balance.Player.AttackCooldown)
	}
}

func TestNewTowerUsesBalanceValues(t *testing.T) {
	world := donburi.NewWorld()
	balance := testBalance()
	config.NewBalance(world, balance)

	if err := NewTower(world, 11, 22); err != nil {
		t.Fatalf("NewTower() error = %v", err)
	}

	entry := Tower.MustFirst(world)
	position := Position.Get(entry)
	health := Health.Get(entry)
	attack := Attack.Get(entry)
	level := Level.Get(entry)

	if position.X != 11 || position.Y != 22 {
		t.Errorf("position = %v,%v, want 11,22", position.X, position.Y)
	}
	if health.Health != balance.Tower.Health || health.MaxHealth != balance.Tower.Health {
		t.Errorf("health = %v/%v, want %v/%v", health.Health, health.MaxHealth, balance.Tower.Health, balance.Tower.Health)
	}
	if attack.Power != balance.Tower.AttackPower {
		t.Errorf("attack power = %v, want %v", attack.Power, balance.Tower.AttackPower)
	}
	if attack.Range != balance.Tower.AttackRange {
		t.Errorf("attack range = %v, want %v", attack.Range, balance.Tower.AttackRange)
	}
	if attack.cooldown.Cooldown != balance.Tower.AttackCooldown {
		t.Errorf("attack cooldown = %v, want %v", attack.cooldown.Cooldown, balance.Tower.AttackCooldown)
	}
	if level.Level != balance.Tower.InitialLevel {
		t.Errorf("level = %v, want %v", level.Level, balance.Tower.InitialLevel)
	}
}

func TestNewCreepUsesBalanceValues(t *testing.T) {
	world := donburi.NewWorld()
	balance := testBalance()
	config.NewBalance(world, balance)

	entry, err := NewCreep(world, 5, 6, 7)
	if err != nil {
		t.Fatalf("NewCreep() error = %v", err)
	}

	creep := Creep.Get(entry)
	velocity := Velocity.Get(entry)
	health := Health.Get(entry)
	attack := Attack.Get(entry)
	sprite := SpriteRender.Get(entry)

	augment := balance.Creep.BigCreepAugment
	wantHealth := balance.Creep.HealthBase + balance.Creep.HealthAugmentMultiplier*augment + 7/balance.Creep.HealthLevelDivisor
	wantPower := balance.Creep.AttackPowerBase + (7-balance.Creep.AttackPowerLevelOffset)*augment/balance.Creep.AttackPowerLevelDivisor

	if sprite.Name != "creep9" {
		t.Errorf("sprite name = %q, want creep9", sprite.Name)
	}
	if creep.GetScoreValue() != balance.Creep.ScoreValueBase*augment {
		t.Errorf("score value = %v, want %v", creep.GetScoreValue(), balance.Creep.ScoreValueBase*augment)
	}
	if velocity.X != 0 || velocity.Y != balance.Creep.BaseVelocityY-augment+7/2 {
		t.Errorf("velocity = %v,%v, want 0,%v", velocity.X, velocity.Y, balance.Creep.BaseVelocityY-augment+7/2)
	}
	if health.Health != wantHealth || health.MaxHealth != wantHealth {
		t.Errorf("health = %v/%v, want %v/%v", health.Health, health.MaxHealth, wantHealth, wantHealth)
	}
	if attack.Power != wantPower {
		t.Errorf("attack power = %v, want %v", attack.Power, wantPower)
	}
	if attack.Range != balance.Creep.AttackRangeBase+balance.Creep.AttackRangeAugmentMultiplier*augment {
		t.Errorf("attack range = %v, want %v", attack.Range, balance.Creep.AttackRangeBase+balance.Creep.AttackRangeAugmentMultiplier*augment)
	}
	if attack.cooldown.Cooldown != balance.Creep.AttackCooldownBase+balance.Creep.AttackCooldownAugmentMultiplier*augment {
		t.Errorf("attack cooldown = %v, want %v", attack.cooldown.Cooldown, balance.Creep.AttackCooldownBase+balance.Creep.AttackCooldownAugmentMultiplier*augment)
	}
}

func TestNewSuperCreepUsesBalanceValues(t *testing.T) {
	world := donburi.NewWorld()
	balance := testBalance()
	config.NewBalance(world, balance)

	entry, err := NewSuperCreep(world, 3, 4)
	if err != nil {
		t.Fatalf("NewSuperCreep() error = %v", err)
	}

	creep := Creep.Get(entry)
	velocity := Velocity.Get(entry)
	health := Health.Get(entry)
	attack := Attack.Get(entry)
	sprite := SpriteRender.Get(entry)

	if sprite.Name != "supercreep" {
		t.Errorf("sprite name = %q, want supercreep", sprite.Name)
	}
	if creep.GetScoreValue() != balance.SuperCreep.ScoreValue {
		t.Errorf("score value = %v, want %v", creep.GetScoreValue(), balance.SuperCreep.ScoreValue)
	}
	if velocity.X != balance.SuperCreep.VelocityX || velocity.Y != balance.SuperCreep.VelocityY {
		t.Errorf("velocity = %v,%v, want %v,%v", velocity.X, velocity.Y, balance.SuperCreep.VelocityX, balance.SuperCreep.VelocityY)
	}
	if health.Health != balance.SuperCreep.Health || health.MaxHealth != balance.SuperCreep.Health {
		t.Errorf("health = %v/%v, want %v/%v", health.Health, health.MaxHealth, balance.SuperCreep.Health, balance.SuperCreep.Health)
	}
	if attack.Power != balance.SuperCreep.AttackPower {
		t.Errorf("attack power = %v, want %v", attack.Power, balance.SuperCreep.AttackPower)
	}
	if attack.Range != balance.SuperCreep.AttackRange {
		t.Errorf("attack range = %v, want %v", attack.Range, balance.SuperCreep.AttackRange)
	}
	if attack.cooldown.Cooldown != balance.SuperCreep.AttackCooldown {
		t.Errorf("attack cooldown = %v, want %v", attack.cooldown.Cooldown, balance.SuperCreep.AttackCooldown)
	}
}

func TestPlayerTowerActionsUseBalanceCosts(t *testing.T) {
	world := donburi.NewWorld()
	balance := testBalance()
	config.NewBalance(world, balance)

	player := &PlayerData{Money: 200, TowerLevels: 10}
	playerEntity := world.Create(Player)
	Player.Set(world.Entry(playerEntity), player)

	towerEntity := world.Create(Tower, Health, Attack, Level)
	towerEntry := world.Entry(towerEntity)
	Tower.Set(towerEntry, &TowerData{})
	Health.Set(towerEntry, &HealthData{Health: 10, MaxHealth: 20})
	Attack.Set(towerEntry, &AttackData{Power: 1, AttackType: RangedSingle, Range: 50, cooldown: util.NewCooldownTimer(30)})
	Level.Set(towerEntry, &LevelData{Level: 1})

	SetGameStats(NewGameStats(nil))
	t.Cleanup(func() { SetGameStats(nil) })

	if !player.TryHealTower(towerEntry, false, false) {
		t.Fatal("TryHealTower() = false, want true")
	}
	if player.Money != 180 {
		t.Errorf("money after heal = %v, want 180", player.Money)
	}

	if !player.TryUpgradeTower(towerEntry, false, false) {
		t.Fatal("TryUpgradeTower() = false, want true")
	}
	if player.Money != 100 {
		t.Errorf("money after upgrade = %v, want 100", player.Money)
	}
	if player.TowerLevels != 11 {
		t.Errorf("tower levels after upgrade = %v, want 11", player.TowerLevels)
	}
	if got := GetGameStats().GetStat("MoneySpent"); got != 100 {
		t.Errorf("MoneySpent = %v, want 100", got)
	}
}
