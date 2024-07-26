package components

import (
	"fmt"
	"image"
	"math/rand"

	"tower-defense/assets"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type CreepData struct {
	scoreValue int
}

var Creep = donburi.NewComponentType[CreepData]()

func NewCreep(w donburi.World, x, y int) *donburi.Entry {
	entity := w.Create(Creep, Position, Velocity, Render, Health, Attack)
	creep := w.Entry(entity)
	Position.SetValue(creep, PositionData{x: x, y: y})
	choose := rand.Intn(2) + 1
	Velocity.SetValue(creep, VelocityData{x: 0, y: 2 + choose})
	name := fmt.Sprintf("creep%v", choose)
	Render.SetValue(creep, *NewRenderer(&SpriteData{image: assets.GetImage(name)}, &RangeRenderData{}, &InfoRenderData{}))
	Creep.SetValue(creep, CreepData{scoreValue: 10 * choose})
	Health.SetValue(creep, HealthData{Health: 1 + 2*choose})
	Attack.SetValue(creep, AttackData{Power: 2 + 2*choose, AttackType: RangedSingle, Range: 10 + 10*choose, Cooldown: 5 + 5*choose})
	return creep
}

func (c *CreepData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)
	// check whether there are any collisions in the new spot

	newPt := image.Pt(v.x, v.y)
	rect := c.GetRect(entry)
	rect = rect.Add(newPt)

	collision := DetectCollisions(entry.World, rect, filter.Or(
		filter.Contains(Creep), // this allows creeps to overlap other creeps but if we don't filter here we deadlock when we get to the entity itself since we're already inside a query for creeps
		filter.Contains(Bullet),
		filter.Contains(Player),
	))
	if collision == nil {
		pos.x += newPt.X
		pos.y += newPt.Y
	}
	// TODO allow creeps to move sideways around the tower?
	v.blocked = collision != nil

	a := Attack.Get(entry)
	a.AttackEnemyRange(entry, Tower, nil)

	return nil
}

func DetectCollisions(world donburi.World, rect image.Rectangle, excludeFilter filter.LayoutFilter) *donburi.Entry {
	var collision *donburi.Entry = nil
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Render, Position),
			filter.Not(excludeFilter),
		),
	)

	query.Each(world, func(testEntry *donburi.Entry) {
		if collision == nil {
			testRect := Render.Get(testEntry).GetRect(testEntry)
			if rect.Overlaps(testRect) {
				collision = testEntry
			}
		}
	})
	return collision
}

func (a *CreepData) GetRect(entry *donburi.Entry) image.Rectangle {
	return Render.Get(entry).GetRect(entry)
}

func (a *CreepData) GetScoreValue() int {
	return a.scoreValue
}
