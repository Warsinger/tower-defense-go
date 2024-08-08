package components

import (
	"fmt"
	"image"

	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type TowerData struct {
	Level int
}

var Tower = donburi.NewComponentType[TowerData]()

type CostList map[string]int

type TowerManagerData struct {
}

var costList = CostList{"Melee": 50, "Ranged": 50}

func (tm *TowerManagerData) GetCostList() CostList {
	return costList
}

func (tm *TowerManagerData) GetCost(name string) int {
	return costList[name]
}
func (tm *TowerManagerData) GetHealCost(name string) int {
	return costList[name] / 2
}

func (tm *TowerManagerData) GetUpgradeCost(name string) int {
	return costList[name]
}

var towerManager = &TowerManagerData{}

func NewTower(world donburi.World, x, y int) error {
	towerEntity := world.Create(Tower, Position, Health, Attack, SpriteRender, RangeRender, InfoRender, NameComponent)
	err := srvsync.NetworkSync(world, &towerEntity, Tower, Position, Health, Attack, SpriteRender, RangeRender, InfoRender, NameComponent)
	if err != nil {
		return err
	}
	tower := world.Entry(towerEntity)

	Position.Set(tower, &PositionData{x, y})
	Health.Set(tower, NewHealthData(20))
	Attack.Set(tower, &AttackData{Power: 1, AttackType: RangedSingle, Range: 50, Cooldown: 30})
	Tower.Set(tower, &TowerData{Level: 1})
	name := Name("tower")
	NameComponent.Set(tower, &name)
	SpriteRender.Set(tower, &SpriteRenderData{})
	RangeRender.Set(tower, &RangeRenderData{})
	InfoRender.Set(tower, &InfoRenderData{})
	return nil
}

func (t *TowerData) Update(entry *donburi.Entry) error {
	a := Attack.Get(entry)
	a.AttackEnemyRange(entry, AfterTowerAttack, Creep)

	return nil
}

func (t *TowerData) Heal(entry *donburi.Entry) bool {
	health := Health.Get(entry)
	if health.Health < health.MaxHealth {
		fmt.Printf("tower healed from %v to %v\n", health.Health, health.MaxHealth)
		health.Health = health.MaxHealth
		return true
	}
	return false
}

const maxLevel = 5

func (t *TowerData) Upgrade(entry *donburi.Entry) bool {
	if t.Level+1 >= maxLevel {
		return false
	}
	t.Level++
	health := Health.Get(entry)
	health.MaxHealth += 10
	health.Health = health.MaxHealth
	attack := Attack.Get(entry)
	attack.Power++
	attack.Range += 3
	attack.Cooldown = max(3, attack.Cooldown-3)

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
