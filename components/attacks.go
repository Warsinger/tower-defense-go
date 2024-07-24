package components

import (
	"image"
	"image/color"
	"tower-defense/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/component"
	"github.com/yohamta/donburi/filter"
)

type HealthData struct {
	Health int
}
type AttackType int

const (
	MeleeSingle AttackType = iota
	RangedSingle
	MeleeArea
	RangedArea
)

type AttackData struct {
	Power      int
	Range      int
	Cooldown   int
	ticker     int
	inCooldown bool
	AttackType AttackType
}

type RangeRenderData struct {
}

var Attack = donburi.NewComponentType[AttackData]()
var Health = donburi.NewComponentType[HealthData]()

func (a *AttackData) GetRect(e *donburi.Entry) image.Rectangle {
	r := Render.Get(e)
	rect := r.GetPrimaryRenderer().GetRect(e)
	ptRange := image.Pt(a.Range, a.Range)
	rect.Min = rect.Min.Sub(ptRange)
	rect.Max = rect.Max.Add(ptRange)
	return rect
}

func (a *AttackData) GetTicker() int {
	return a.ticker
}
func (a *AttackData) IncrementTicker() {
	if a.inCooldown {
		a.ticker++
	}
}

func (a *AttackData) CheckCooldown() {
	if a.ticker >= a.Cooldown {
		a.ticker = 0
		a.inCooldown = false
	}
}

func (a *AttackData) StartCooldown() {
	a.inCooldown = true
}

func (rr *RangeRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	a := Attack.Get(entry)
	aRect := a.GetRect(entry)
	aPt := util.MidpointRect(aRect)

	vector.StrokeCircle(screen, float32(aPt.X), float32(aPt.Y), float32(aRect.Dx()/2), 1, color.White, true)
}

func (rr *RangeRenderData) GetRect(e *donburi.Entry) image.Rectangle {
	a := Attack.Get(e)
	return a.GetRect(e)
}

func (a *AttackData) FindEnemyInRange(e *donburi.Entry, enemyType component.IComponentType) *donburi.Entry {
	// query for enemies then find the closest one
	aRect := a.GetRect(e)
	minDist := float64(a.Range + 1 + aRect.Dy()/2)
	var foundEnemy *donburi.Entry = nil
	query := donburi.NewQuery(filter.Contains(enemyType))
	query.Each(e.World, func(enemyEntry *donburi.Entry) {
		enemy := Render.Get(enemyEntry)
		eRect := enemy.GetRect(enemyEntry)

		dist := util.DistanceRects(aRect, eRect)
		// fmt.Printf("enemy at distance %v\n", dist)
		if dist < minDist {
			minDist = dist
			foundEnemy = enemyEntry
			// fmt.Println("found enemy")
		}
	})
	return foundEnemy
}

func (a *AttackData) AttackEnemy(entry *donburi.Entry, enemyType component.IComponentType, afterKill func(*donburi.Entry, *donburi.Entry), afterAttack func(*donburi.Entry)) {
	a.CheckCooldown()
	if a.GetTicker() == 0 {
		// fmt.Printf("finding enemies in range of %v\n", e.Entity())
		// look for a creep in range to shoot at
		enemy := a.FindEnemyInRange(entry, enemyType)
		if enemy != nil {
			enemyHealth := Health.Get(enemy)
			attack := Attack.Get(entry)
			remainingHealth := enemyHealth.Health - attack.Power
			if remainingHealth <= 0 {
				// kill enemy, remove from board

				// do some other stuff in a callback
				if afterKill != nil {
					afterKill(entry, enemy)
				}
				enemy.Remove()
			} else {
				enemyHealth.Health = remainingHealth
			}
			a.StartCooldown()
		}
	}
	a.IncrementTicker()
	if afterAttack != nil {
		afterAttack(entry)
	}
}
