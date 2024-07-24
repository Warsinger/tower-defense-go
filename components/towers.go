package components

import (
	"fmt"
	"image"
	"tower-defense/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
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
	Health.SetValue(tower, HealthData{20})
	Render.SetValue(tower, *NewRenderer(&SpriteData{image: assets.GetImage("tower")}, &RangeRenderData{}, &TowerRenderData{}))
	Attack.SetValue(tower, AttackData{Power: 1, AttackType: RangedSingle, Range: 40, Cooldown: 30})
	return nil
}

func (t *TowerData) Update(entry *donburi.Entry) error {
	a := Attack.Get(entry)
	a.AttackEnemyRange(entry, Creep, t.OnKillEnemy, t.AfterAttackEnemy)

	return nil
}

func (t *TowerData) AfterAttackEnemy(towerEntry *donburi.Entry) {
	towerHealth := Health.Get(towerEntry)
	towerHealth.Health--
	if towerHealth.Health <= 0 {
		towerEntry.Remove()
	}
}

func (t *TowerData) OnKillEnemy(towerEntry *donburi.Entry, enemyEntry *donburi.Entry) {
	enemy := Creep.Get(enemyEntry)
	score := enemy.GetScoreValue()

	pe := Player.MustFirst(enemyEntry.World)
	player := Player.Get(pe)
	player.AddMoney(score)
	player.AddScore(score)
}

func (t *TowerRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	a := Attack.Get(entry)
	h := Health.Get(entry)
	r := Render.Get(entry)
	rect := r.GetRect(entry)

	// draw health and cooldown
	var cd int = 0
	if a.inCooldown {
		cd = a.Cooldown - a.GetTicker()
	}
	str := fmt.Sprintf("HP %d\\CD %d", h.Health, cd)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y-20))
	text.Draw(screen, str, assets.InfoFace, op)
}

func (t *TowerRenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	panic("TowerRenderData.GetRect() unimplemented")
}
