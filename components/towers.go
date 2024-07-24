package components

import (
	"tower-defense/assets"

	"github.com/yohamta/donburi"
)

type TowerData struct {
}

var Tower = donburi.NewComponentType[TowerData]()

type CostList map[string]int

type TowerManagerData struct {
}
type RangeRenderData struct {
}

var costList = CostList{"Melee": 50, "Ranged": 50}

func (tm *TowerManagerData) GetCostList() CostList {
	return costList
}

func (tm *TowerManagerData) GetCost(name string) int {
	return costList[name]
}

var towerManager = &TowerManagerData{}

func NewTower(w donburi.World, x, y int) error {
	towerEntity := w.Create(Tower, Position, Render, Health, Attack)
	tower := w.Entry(towerEntity)

	Position.SetValue(tower, PositionData{x, y})
	Health.SetValue(tower, HealthData{50})
	Render.SetValue(tower, RenderData{primaryRenderer: &SpriteData{image: assets.GetImage("tower")}, secondaryRenderer: &RangeRenderData{}})
	Attack.SetValue(tower, AttackData{Power: 1, AttackType: RangedSingle, Range: 30, Cooldown: 30})
	return nil
}
