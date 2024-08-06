package components

import (
	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/yohamta/donburi"
)

type TowerData struct {
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

var towerManager = &TowerManagerData{}

func NewTower(world donburi.World, x, y int) error {
	towerEntity := world.Create(Tower, Position, Render, Health, Attack)
	_ = srvsync.NetworkSync(world, &towerEntity, Tower, Position, Render, Health, Attack)
	tower := world.Entry(towerEntity)

	Position.Set(tower, &PositionData{x, y})
	Health.Set(tower, NewHealthData(20))
	Render.Set(tower, NewRenderer(NewSprite("tower"), &RangeRenderData{}, &InfoRenderData{}))
	Attack.Set(tower, &AttackData{Power: 1, AttackType: RangedSingle, Range: 50, Cooldown: 30})
	return nil
}

func (t *TowerData) Update(entry *donburi.Entry) error {
	a := Attack.Get(entry)
	a.AttackEnemyRange(entry, AfterTowerAttack, Creep)

	return nil
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
