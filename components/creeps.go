package components

import (
	"fmt"
	"image"
	"math/rand/v2"

	"tower-defense/util"

	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/yohamta/donburi"
)

type CreepData struct {
	scoreValue int
}

var Creep = donburi.NewComponentType[CreepData]()

func NewCreep(world donburi.World, x, y int) *donburi.Entry {
	entity := world.Create(Creep, Position, Velocity, Render, Health, Attack)
	_ = srvsync.NetworkSync(world, &entity, Creep, Position, Render, Health, Attack)
	creep := world.Entry(entity)
	Position.Set(creep, &PositionData{X: x, Y: y})

	const bigCreepChance = 0.3
	choose := 1
	augment := 1
	if rand.Float32() < bigCreepChance {
		choose = 4
		augment = 2
	} else {
		choose += rand.IntN(3)
	}
	Velocity.Set(creep, &VelocityData{X: 0, Y: 5 - augment})
	name := fmt.Sprintf("creep%v", choose)
	Render.Set(creep, NewRenderer(NewSprite(name), &RangeRenderData{}, &InfoRenderData{}))
	Creep.Set(creep, &CreepData{scoreValue: 10 * augment})
	Health.Set(creep, NewHealthData(1+2*augment))
	Attack.Set(creep, &AttackData{Power: 2 + 2*augment, AttackType: RangedSingle, Range: 10 + 10*augment, Cooldown: 5 + 5*augment})
	return creep
}

func (c *CreepData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)
	// check whether there are any collisions in the new spot

	newPt := image.Pt(v.X, v.Y)
	rect := c.GetRect(entry)
	rect = rect.Add(newPt)

	// HACK: Creep must be in the exclusion filter, this allows creeps to overlap other creeps
	// but if we don't filter here we deadlock when we get to the entity itself since we're already inside a query for creeps
	collision := DetectCollisions(entry.World, rect, util.CreateOrFilter(Creep, Bullet))
	if collision == nil {
		pos.X += newPt.X
		pos.Y += newPt.Y
	}
	// TODO allow creeps to move sideways around the tower? (if so don't allow for player)
	v.blocked = collision != nil

	a := Attack.Get(entry)
	a.AttackEnemyRange(entry, nil, Tower, Player)

	return nil
}

func (a *CreepData) GetRect(entry *donburi.Entry) image.Rectangle {
	return Render.Get(entry).GetRect(entry)
}

func (a *CreepData) GetScoreValue() int {
	return a.scoreValue
}
