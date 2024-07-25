package components

import (
	"image"
	"image/color"
	"tower-defense/assets"
	"tower-defense/config"
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
	config := config.Config.Get(config.Config.MustFirst(entry.World))

	if config.IsDebug() {
		a := Attack.Get(entry)
		aRect := a.GetRect(entry)
		aPt := util.MidpointRect(aRect)

		vector.StrokeCircle(screen, float32(aPt.X), float32(aPt.Y), float32(aRect.Dx()/2), 1, color.White, true)
	}
}

func (rr *RangeRenderData) GetRect(e *donburi.Entry) image.Rectangle {
	a := Attack.Get(e)
	return a.GetRect(e)
}

func (a *AttackData) FindEnemyRange(entry *donburi.Entry, enemyType component.IComponentType) *donburi.Entry {
	// query for enemies then find the closest one
	aRect := a.GetRect(entry)
	minDist := float64(a.Range + 1 + aRect.Dy()/2)
	var foundEnemy *donburi.Entry = nil
	query := donburi.NewQuery(filter.Contains(enemyType))
	query.Each(entry.World, func(enemyEntry *donburi.Entry) {
		// fmt.Printf("checking distance of %v\n", enemyEntry)
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

func (a *AttackData) FindEnemyIntersect(entry *donburi.Entry, enemyType component.IComponentType) *donburi.Entry {
	// query for first enemy we intersect with
	aRect := a.GetRect(entry)

	var foundEnemy *donburi.Entry = nil
	query := donburi.NewQuery(filter.Contains(enemyType))
	query.Each(entry.World, func(enemyEntry *donburi.Entry) {
		if foundEnemy != nil {
			return
		}
		// fmt.Printf("checking distance of %v\n", enemyEntry)
		enemy := Render.Get(enemyEntry)
		eRect := enemy.GetRect(enemyEntry)

		if aRect.Overlaps(eRect) {
			foundEnemy = enemyEntry
			// fmt.Println("found enemy")
		}
	})
	return foundEnemy
}

func (a *AttackData) AttackEnemyRange(entry *donburi.Entry, enemyType component.IComponentType, afterAttack func(*donburi.Entry)) {
	a.CheckCooldown()
	if a.GetTicker() == 0 {
		// fmt.Printf("finding enemies in range of %v\n", entry)
		// look for a enemy in range to shoot at
		enemy := a.FindEnemyRange(entry, enemyType)
		if enemy != nil {
			a.LaunchBullet(entry, enemy)
			a.StartCooldown()
			if afterAttack != nil {
				afterAttack(entry)
			}
		}
	}
	a.IncrementTicker()
}

func (a *AttackData) LaunchBullet(entry *donburi.Entry, enemy *donburi.Entry) {
	// create a bullet path from the midpoint of the launcher to the midpoint of the enemy

	r1 := Render.Get(entry).GetRect(entry)
	r2 := Render.Get(enemy).GetRect(enemy)

	start := util.MidpointRect(r1)
	end := util.MidpointRect(r2)
	const bulletSpeed = 6
	if enemy.HasComponent(Velocity) {
		// how far ahead to lead, distance to target divided by speed
		lead := util.Abs(util.DistancePoints(start, end))/bulletSpeed - 0.5
		velocity := Velocity.Get(enemy)
		end.Y += int(float64(velocity.y) * lead)
	}

	NewBullet(entry.World, start, end, bulletSpeed, enemy.HasComponent(Tower))
	assets.PlaySound("shoot")
}

func (a *AttackData) AttackEnemyIntersect(entry *donburi.Entry, enemyType component.IComponentType, afterKill func(*donburi.Entry, *donburi.Entry), afterAttack func(*donburi.Entry)) {
	a.CheckCooldown()
	if a.GetTicker() == 0 {
		// fmt.Printf("finding enemies in range of %v\n", entry)
		// look for a enemy we interect
		enemy := a.FindEnemyIntersect(entry, enemyType)
		if enemy != nil {
			enemyHealth := Health.Get(enemy)
			attack := Attack.Get(entry)
			remainingHealth := enemyHealth.Health - attack.Power
			if remainingHealth <= 0 {
				// kill enemy, remove from board, plays sound
				assets.PlaySound("explosion")

				// do some other stuff in a callback
				if afterKill != nil {
					afterKill(entry, enemy)
				}
				enemy.Remove()
			} else {
				enemyHealth.Health = remainingHealth
			}
			a.StartCooldown()
			if afterAttack != nil {
				afterAttack(entry)
			}
		}
	}
	a.IncrementTicker()
}
