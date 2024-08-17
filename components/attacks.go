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
)

type HealthData struct {
	Health    int
	MaxHealth int
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
	cooldown   *util.CooldownTimer
	AttackType AttackType
	noLead     bool
}

type LevelData struct {
	Level int
}

type RangeRenderData struct {
}

var Attack = donburi.NewComponentType[AttackData]()
var Health = donburi.NewComponentType[HealthData]()
var Level = donburi.NewComponentType[LevelData]()
var RangeRender = donburi.NewComponentType[RangeRenderData]()

func NewHealthData(health int) *HealthData {
	return &HealthData{Health: health, MaxHealth: health}
}

func (a *AttackData) GetExpandedRect(e *donburi.Entry) image.Rectangle {
	rect := GetRect(e)
	ptRange := image.Pt(a.Range, a.Range)
	rect.Min = rect.Min.Sub(ptRange)
	rect.Max = rect.Max.Add(ptRange)
	return rect
}

func (rr *RangeRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	config := config.GetConfig(entry.World)

	if config.Debug {
		a := Attack.Get(entry)
		aRect := a.GetExpandedRect(entry)
		aPt := util.MidpointRect(aRect)

		vector.StrokeCircle(screen, float32(aPt.X), float32(aPt.Y), float32(aRect.Dx()/2), 1, color.White, true)
		//vector.StrokeRect(screen, float32(aRect.Min.X), float32(aRect.Min.Y), float32(aRect.Dx()), float32(aRect.Dy()), 1, color.White, true)
	}
}

func (a *AttackData) FindEnemyRange(entry *donburi.Entry, enemyType ...component.IComponentType) *donburi.Entry {
	// query for enemies then find the closest one
	aRect := a.GetExpandedRect(entry)
	// this just sets an upper bounds on the distance
	minDist := 2000.0
	// maxRange := float64(aRect.Dx()/2 + aRect.Dy()/2)
	var foundEnemy *donburi.Entry = nil
	query := donburi.NewQuery(util.CreateOrFilter(enemyType...))
	query.Each(entry.World, func(enemyEntry *donburi.Entry) {
		// fmt.Printf("checking distance of %v\n", enemyEntry)
		eRect := GetRect(enemyEntry)

		dist := util.DistanceRects(aRect, eRect)

		// this will fire at the closest enemy but we might want to prioritize others such as lowest health or things below us
		if aRect.Overlaps(eRect) && dist < minDist {
			// fmt.Printf("enemy at distance %v\n", dist)
			minDist = dist
			foundEnemy = enemyEntry
			// fmt.Println("found enemy")
		}
	})
	return foundEnemy
}

func (a *AttackData) FindEnemyIntersect(entry *donburi.Entry, enemyType ...component.IComponentType) *donburi.Entry {
	// query for first enemy we intersect with
	rect := a.GetExpandedRect(entry)

	var foundEnemy *donburi.Entry = nil
	query := donburi.NewQuery(util.CreateOrFilter(enemyType...))
	query.Each(entry.World, func(enemyEntry *donburi.Entry) {
		if foundEnemy != nil {
			return
		}
		// fmt.Printf("checking distance of %v\n", enemyEntry)
		eRect := GetRect(enemyEntry)

		if rect.Overlaps(eRect) {
			foundEnemy = enemyEntry
			// fmt.Println("found enemy")
		}
	})
	return foundEnemy
}

func (a *AttackData) AttackEnemyRange(entry *donburi.Entry, afterAttack func(*donburi.Entry), enemyType ...component.IComponentType) {
	a.cooldown.CheckCooldown()
	defer a.cooldown.IncrementTicker()
	if !a.cooldown.InCooldown {
		// fmt.Printf("finding enemies in range of %v\n", entry)
		// look for a enemy in range to shoot at
		enemy := a.FindEnemyRange(entry, enemyType...)
		if enemy != nil {
			a.LaunchBullet(entry, enemy)
			a.cooldown.StartCooldown()
			if afterAttack != nil {
				afterAttack(entry)
			}
		}
	}

}

func (a *AttackData) LaunchBullet(entry *donburi.Entry, enemy *donburi.Entry) {
	// create a bullet path from the midpoint of the launcher to the midpoint of the enemy

	ownRect := GetRect(entry)
	enemyRect := GetRect(enemy)

	start := util.MidpointRect(ownRect)
	end := util.MidpointRect(enemyRect)
	const bulletSpeed = 8
	if !a.noLead && enemy.HasComponent(Velocity) {
		v := Velocity.Get(enemy)
		if !v.blocked {
			// how far ahead to lead, distance to target divided by speed
			lead := util.Abs(util.DistancePoints(start, end))/bulletSpeed - 0.5
			end.Y += int(float64(v.Y) * lead)

			// if the enemy is currently above then don't lead below our position
			if enemyRect.Max.Y < ownRect.Min.Y {
				end.Y = ownRect.Min.Y
			}
		}
	}

	creep := entry.HasComponent(Creep)
	NewBullet(entry.World, start, end, bulletSpeed, creep)
	var sound string
	if creep {
		sound = "shoot1"
	} else {
		sound = "shoot2"
	}
	assets.PlaySound(sound)
}

func (a *AttackData) AttackEnemyIntersect(entry *donburi.Entry, afterKill func(*donburi.Entry, *donburi.Entry), afterAttack func(*donburi.Entry), enemyType ...component.IComponentType) {
	a.cooldown.CheckCooldown()
	defer a.cooldown.IncrementTicker()
	if !a.cooldown.InCooldown {
		// fmt.Printf("finding enemies in range of %v\n", entry)
		// look for a enemy we interect
		enemy := a.FindEnemyIntersect(entry, enemyType...)
		if enemy != nil {
			enemyHealth := Health.Get(enemy)
			attack := Attack.Get(entry)
			enemyHealth.Health = enemyHealth.Health - attack.Power
			if enemyHealth.Health <= 0 {
				// kill enemy, remove from board, plays sound
				assets.PlaySound("explosion")

				// do some other stuff in a callback
				if afterKill != nil {
					afterKill(entry, enemy)
				}
				// HACK Don't remove player upon kill, TODO find a better way to handle this
				if !enemy.HasComponent(Player) {
					enemy.Remove()
				}
			}
			a.cooldown.StartCooldown()
			if afterAttack != nil {
				afterAttack(entry)
			}
		}
	}
}
