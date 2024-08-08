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

func NewCreep(world donburi.World, x, y, creepLevel int) (*donburi.Entry, error) {
	entity := world.Create(Creep, Position, Velocity, Health, Attack, SpriteRender, RangeRender, InfoRender, NameComponent)
	err := srvsync.NetworkSync(world, &entity, Creep, Position, Health, Attack, SpriteRender, RangeRender, InfoRender, NameComponent)
	if err != nil {
		return nil, err
	}
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
	Velocity.Set(creep, &VelocityData{X: 0, Y: 5 - augment + creepLevel/2})
	name := Name(fmt.Sprintf("creep%v", choose))
	NameComponent.Set(creep, &name)
	Creep.Set(creep, &CreepData{scoreValue: 10 * augment})
	Health.Set(creep, NewHealthData(1+2*augment+creepLevel/4))
	Attack.Set(creep, &AttackData{Power: 2 + 2*augment, AttackType: RangedSingle, Range: 10 + 10*augment, Cooldown: 5 + 5*augment})
	SpriteRender.Set(creep, &SpriteRenderData{})
	RangeRender.Set(creep, &RangeRenderData{})
	InfoRender.Set(creep, &InfoRenderData{})
	return creep, nil
}

func (c *CreepData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)
	// check whether there are any collisions in the new spot

	newPt := image.Pt(v.X, v.Y)
	rect := GetRect(entry)
	rect = rect.Add(newPt)

	// HACK: Creep must be in the exclusion filter, this allows creeps to overlap other creeps
	// but if we don't filter here we deadlock when we get to the entity itself since we're already inside a query for creeps
	collision := DetectCollisionsEntry(entry, rect, util.CreateOrFilter(Bullet))
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

func (a *CreepData) GetScoreValue() int {
	return a.scoreValue
}
