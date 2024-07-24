package components

import (
	"fmt"
	"image"
	"tower-defense/assets"
	"tower-defense/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type TowerData struct {
}

type TowerRenderData struct {
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

func NewTower(w donburi.World, x, y int) error {
	towerEntity := w.Create(Tower, Position, Render, Health, Attack)
	tower := w.Entry(towerEntity)

	Position.SetValue(tower, PositionData{x, y})
	Health.SetValue(tower, HealthData{50})
	Render.SetValue(tower, *NewRenderer(&SpriteData{image: assets.GetImage("tower")}, &RangeRenderData{}, &TowerRenderData{}))
	Attack.SetValue(tower, AttackData{Power: 1, AttackType: RangedSingle, Range: 50, Cooldown: 10})
	return nil
}

func (t *TowerData) Update(e *donburi.Entry) error {
	a := Attack.Get(e)
	a.CheckCooldown()
	if a.GetTicker() == 0 {
		// fmt.Printf("finding creeps in range of %v\n", e.Entity())
		// look for a creep in range to shoot at
		creep := t.FindCreepInRange(e)
		if creep != nil {
			t.AttackCreep(creep)
			a.IncrementTicker()
		}
	}

	return nil
}

func (t *TowerData) FindCreepInRange(e *donburi.Entry) *donburi.Entry {
	// query for creeps then find the closest one
	a := Attack.Get(e)
	aRect := a.GetRect(e)
	minDist := float64(a.Range + 1 + aRect.Dy()/2)
	var foundCreep *donburi.Entry = nil
	query := donburi.NewQuery(filter.Contains(Creep))
	query.Each(e.World, func(ce *donburi.Entry) {
		creep := Creep.Get(ce)
		cRect := creep.GetRect(ce)

		dist := util.DistanceRects(aRect, cRect)
		// fmt.Printf("creep at distance %v\n", dist)
		if dist < minDist {
			minDist = dist
			foundCreep = ce
			// fmt.Println("found creep")
		}
	})
	return foundCreep
}

func (t *TowerData) AttackCreep(ce *donburi.Entry) error {
	// fmt.Printf("attacking creep %v, %v\n", t, ce)
	creepHealth := Health.Get(ce)
	a := Attack.Get(ce)
	remainingHealth := creepHealth.Health - a.Power
	if remainingHealth <= 0 {
		// kill creep, remove from board, take its money
		creep := Creep.Get(ce)
		score := creep.GetScoreValue()

		pe := Player.MustFirst(ce.World)
		player := Player.Get(pe)
		player.AddMoney(score)
		player.AddScore(score)

		ce.Remove()
	} else {
		creepHealth.Health = remainingHealth
	}
	return nil
}

func (t *TowerRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	a := Attack.Get(entry)
	h := Health.Get(entry)
	r := Render.Get(entry)
	rect := r.GetRect(entry)

	// draw health and cooldown
	str := fmt.Sprintf("HP %d\\CD %d", h.Health, a.Cooldown)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y-20))
	text.Draw(screen, str, assets.InfoFace, op)
}

func (t *TowerRenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	panic("TowerRenderData.GetRect() unimplemented")
}
