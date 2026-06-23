package components

import (
	"fmt"
	"image"
	"math"
	"tower-defense/config"
	"tower-defense/util"

	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type TowerData struct {
}

var Tower = donburi.NewComponentType[TowerData]()

func getTowerCost(world donburi.World, name string) int {
	return config.GetBalance(world).Tower.Costs[name]
}

func getTowerHealCost(world donburi.World, name string) int {
	balance := config.GetBalance(world).Tower
	return balance.Costs[name] / balance.HealCostDivisor
}

func getTowerUpgradeCost(world donburi.World, name string) int {
	return getTowerCost(world, name)
}

func NewTower(world donburi.World, x, y int) error {
	towerEntity := world.Create(Tower, Position, Health, Attack, Level, SpriteRender, RangeRender, InfoRender)
	err := srvsync.NetworkSync(world, &towerEntity, Tower, Position, Health, Attack, Level, SpriteRender, RangeRender, InfoRender)
	if err != nil {
		return err
	}
	tower := world.Entry(towerEntity)

	balance := config.GetBalance(world).Tower
	Position.Set(tower, &PositionData{x, y})
	Health.Set(tower, NewHealthData(balance.Health))
	Attack.Set(tower, &AttackData{Power: balance.AttackPower, AttackType: RangedSingle, Range: balance.AttackRange, cooldown: util.NewCooldownTimer(balance.AttackCooldown)})
	Level.Set(tower, &LevelData{Level: balance.InitialLevel})
	SpriteRender.Set(tower, &SpriteRenderData{Name: "tower"})
	RangeRender.Set(tower, &RangeRenderData{})
	InfoRender.Set(tower, &InfoRenderData{})
	return nil
}

func (t *TowerData) Update(entry *donburi.Entry) error {
	a := Attack.Get(entry)
	a.AttackEnemyRange(entry, AfterTowerAttack, Creep)

	return nil
}

func (t *TowerData) Heal(entry *donburi.Entry, debug bool) bool {
	health := Health.Get(entry)
	if health.Health < health.MaxHealth {
		if debug {
			fmt.Printf("tower healed from %v to %v\n", health.Health, health.MaxHealth)
		}
		health.Health = health.MaxHealth
		GetGameStats().IncrementStat("TowersHealed")
		return true
	}
	return false
}

func GetMaxTowerLevel(world donburi.World) int {
	player := Player.Get(Player.MustFirst(world))
	return player.GetMaxTowerLevel(config.GetBalance(world))
}
func (p *PlayerData) GetMaxTowerLevel(balance *config.BalanceData) int {
	return balance.Player.MaxTowerInitialLevel + int(math.Trunc(float64(p.TowerLevels)/float64(balance.Player.MaxTowerLevelsPerBonus)))
}

func (t *TowerData) Upgrade(entry *donburi.Entry, debug bool) bool {
	level := Level.Get(entry)
	if level.Level >= GetMaxTowerLevel(entry.World) {
		if debug {
			fmt.Printf("Tower is max level %v\n", level.Level)
		}
		return false
	}
	level.Level++
	health := Health.Get(entry)
	balance := config.GetBalance(entry.World).Tower
	health.MaxHealth += balance.UpgradeMaxHealthAdd
	health.Health = health.MaxHealth
	attack := Attack.Get(entry)
	attack.Power += level.Level / balance.UpgradePowerLevelDivisor
	attack.Range += balance.UpgradeRangeAdd
	attack.cooldown.Cooldown = max(balance.UpgradeMinCooldown, attack.cooldown.Cooldown-balance.UpgradeCooldownReduction)
	GetGameStats().IncrementStat("TowersUpgraded")

	return true
}

func findTower(world donburi.World, x, y int) *donburi.Entry {
	query := donburi.NewQuery(filter.Contains(Tower))
	var foundEntry *donburi.Entry
	pt := image.Pt(x, y)
	query.Each(world, func(entry *donburi.Entry) {
		if foundEntry == nil && pt.In(GetRect(entry)) {
			foundEntry = entry
		}
	})
	return foundEntry
}

func AfterTowerAttack(towerEntry *donburi.Entry) {
	towerHealth := Health.Get(towerEntry)
	towerHealth.Health--
	if towerHealth.Health <= 0 {
		towerEntry.Remove()
		GetGameStats().IncrementStat("TowersAmmoOut")
	}
}

func OnKillCreep(towerEntry *donburi.Entry, enemyEntry *donburi.Entry) {
	enemy := Creep.Get(enemyEntry)
	score := enemy.GetScoreValue()

	pe := Player.MustFirst(enemyEntry.World)
	player := Player.Get(pe)
	player.AddMoney(score)
	player.AddScore(score)
}
