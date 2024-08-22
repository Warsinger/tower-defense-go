package components

import (
	"fmt"
	"image"
	"math"
	"tower-defense/util"

	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type TowerData struct {
}

var Tower = donburi.NewComponentType[TowerData]()

type CostList map[string]int

type TowerManagerData struct {
	costList CostList
}

var towerManager = &TowerManagerData{costList: CostList{"Melee": 50, "Ranged": 50}}

func (tm *TowerManagerData) GetCostList() CostList {
	return tm.costList
}

func (tm *TowerManagerData) GetCost(name string) int {
	return tm.costList[name]
}
func (tm *TowerManagerData) GetHealCost(name string) int {
	return tm.costList[name] / 2
}

func (tm *TowerManagerData) GetUpgradeCost(name string) int {
	return tm.costList[name]
}

func NewTower(world donburi.World, x, y int) error {
	towerEntity := world.Create(Tower, Position, Health, Attack, Level, SpriteRender, RangeRender, InfoRender)
	err := srvsync.NetworkSync(world, &towerEntity, Tower, Position, Health, Attack, Level, SpriteRender, RangeRender, InfoRender)
	if err != nil {
		return err
	}
	tower := world.Entry(towerEntity)

	Position.Set(tower, &PositionData{x, y})
	Health.Set(tower, NewHealthData(20))
	Attack.Set(tower, &AttackData{Power: 1, AttackType: RangedSingle, Range: 50, cooldown: util.NewCooldownTimer(30)})
	Level.Set(tower, &LevelData{Level: 1})
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
		GetGameStats().UpdateTowersHealed()
		return true
	}
	return false
}

const initMaxTowerLevel = 5

func GetMaxTowerLevel(world donburi.World) int {
	player := Player.Get(Player.MustFirst(world))
	return player.GetMaxTowerLevel()
}
func (p *PlayerData) GetMaxTowerLevel() int {
	return initMaxTowerLevel + int(math.Trunc(float64(p.TowerLevels)/20))
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
	health.MaxHealth += 5
	health.Health = health.MaxHealth
	attack := Attack.Get(entry)
	attack.Power += level.Level / 3
	attack.Range += 3
	attack.cooldown.Cooldown = max(3, attack.cooldown.Cooldown-3)
	GetGameStats().UpdateTowersUpgraded()

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
		GetGameStats().UpdateTowersAmmoOut()
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
