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
	entity := world.Create(Creep, Position, Velocity, Health, Attack, SpriteRender, RangeRender, InfoRender)
	err := srvsync.NetworkSync(world, &entity, Creep, Position, Health, Attack, SpriteRender, RangeRender, InfoRender)
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
	name := fmt.Sprintf("creep%v", choose)
	Creep.Set(creep, &CreepData{scoreValue: 10 * augment})
	Health.Set(creep, NewHealthData(1+2*augment+creepLevel/4))
	Attack.Set(creep, &AttackData{Power: 2 + 2*augment, AttackType: RangedSingle, Range: 10 + 10*augment, Cooldown: 5 + 5*augment})
	SpriteRender.Set(creep, &SpriteRenderData{Name: name})
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

	collision := DetectCollisionsEntry(entry, rect, util.CreateOrFilter(Bullet))
	if collision == nil {
		pos.X += newPt.X
		pos.Y += newPt.Y
	} else if collision.HasComponent(Creep) {
		// creep collides with another creep, so let it move a little to the side
		collRect := GetRect(collision)

		if rect.Min.X <= collRect.Min.X {
			pos.X = min(collRect.Min.X-rect.Dx(), pos.X-3)
		} else if rect.Max.X > collRect.Max.X {
			pos.X = max(collRect.Max.X+1, pos.X+3)
		}
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
