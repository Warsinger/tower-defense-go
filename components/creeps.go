package components

import (
	"fmt"
	"image"
	"math/rand/v2"

	"tower-defense/config"
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

	balance := config.GetBalance(world).Creep
	choose := balance.SmallCreepFirstChoice
	augment := 1
	if rand.Float32() < balance.BigCreepChance {
		choose = balance.BigCreepChoice
		augment = balance.BigCreepAugment
	} else {
		choose += rand.IntN(balance.SmallCreepVariants)
	}
	Velocity.Set(creep, &VelocityData{X: 0, Y: balance.BaseVelocityY - augment + creepLevel/2})
	name := fmt.Sprintf("creep%v", choose)
	Creep.Set(creep, &CreepData{scoreValue: balance.ScoreValueBase * augment})
	Health.Set(creep, NewHealthData(balance.HealthBase+balance.HealthAugmentMultiplier*augment+creepLevel/balance.HealthLevelDivisor))
	Attack.Set(creep, &AttackData{
		Power:      balance.AttackPowerBase + (creepLevel-balance.AttackPowerLevelOffset)*augment/balance.AttackPowerLevelDivisor,
		AttackType: RangedSingle,
		Range:      balance.AttackRangeBase + balance.AttackRangeAugmentMultiplier*augment,
		cooldown:   util.NewCooldownTimer(balance.AttackCooldownBase + balance.AttackCooldownAugmentMultiplier*augment),
	})
	SpriteRender.Set(creep, &SpriteRenderData{Name: name})
	RangeRender.Set(creep, &RangeRenderData{})
	InfoRender.Set(creep, &InfoRenderData{})
	return creep, nil
}
func NewSuperCreep(world donburi.World, x, y int) (*donburi.Entry, error) {
	entity := world.Create(Creep, Position, Velocity, Health, Attack, SpriteRender, RangeRender, InfoRender)
	err := srvsync.NetworkSync(world, &entity, Creep, Position, Health, Attack, SpriteRender, RangeRender, InfoRender)
	if err != nil {
		return nil, err
	}
	creep := world.Entry(entity)
	Position.Set(creep, &PositionData{X: x, Y: y})
	balance := config.GetBalance(world).SuperCreep
	Velocity.Set(creep, &VelocityData{X: balance.VelocityX, Y: balance.VelocityY})
	name := "supercreep"
	Creep.Set(creep, &CreepData{scoreValue: balance.ScoreValue})
	Health.Set(creep, NewHealthData(balance.Health))
	Attack.Set(creep, &AttackData{Power: balance.AttackPower, AttackType: RangedSingle, Range: balance.AttackRange, cooldown: util.NewCooldownTimer(balance.AttackCooldown)})
	SpriteRender.Set(creep, &SpriteRenderData{Name: name})
	RangeRender.Set(creep, &RangeRenderData{})
	InfoRender.Set(creep, &InfoRenderData{})
	return creep, nil
}

const maxTryMove = 10

func (c *CreepData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)
	newPt := image.Pt(pos.X+v.X, pos.Y+v.Y)
	if c.TryMoveTo(entry, pos, newPt, maxTryMove) {
		v.blocked = false
	} else {
		v.blocked = true
	}
	a := Attack.Get(entry)
	a.AttackEnemyRange(entry, nil, Tower, Player)

	return nil
}

// try to move to the new point, returning true if successful, keep trying new points up to a maximum
func (c *CreepData) TryMoveTo(entry *donburi.Entry, curPos *PositionData, newPt image.Point, maxTry int) bool {
	if maxTry <= 0 {
		return false
	}
	be := Board.MustFirst(entry.World)
	board := Board.Get(be)

	// check whether there are any collisions in the new spot
	rect := GetRect(entry)
	newRect := image.Rect(newPt.X, newPt.Y, newPt.X+rect.Dx(), newPt.Y+rect.Dy())

	collision := DetectCollisionsEntry(entry.World, entry.Entity(), newRect, util.CreateOrFilter(Bullet))
	if collision == nil {
		curPos.X = newPt.X
		curPos.Y = newPt.Y
		return true
	} else if collision.HasComponent(Creep) {
		// creep collides with another creep so let it move a little to the side
		collRect := GetRect(collision)

		var newX int
		if rect.Min.X <= collRect.Min.X {
			newX = max(min(collRect.Min.X-rect.Dx(), curPos.X-3), 0)
		} else if rect.Max.X > collRect.Max.X {
			newX = min(max(collRect.Max.X+1, curPos.X+3), board.Width-rect.Dx())
		}

		return c.TryMoveTo(entry, curPos, image.Pt(newX, newPt.Y), maxTry-1)
	} else {
		// BUG creeps can still get stuck on each other just above the player and not be able to attack, and the player can't attack them
		// if it can go a little way then let it move as far as possible
		// collRect := GetRect(collision)

		// if rect.Max.Y >= collRect.Min.Y {
		// 	newY := collRect.Min.Y - 1
		// 	return c.TryMoveTo(entry, curPos, image.Pt(newPt.X, newY), maxTry-1)
		// }

		// try to creep forward just a bit
		return c.TryMoveTo(entry, curPos, image.Pt(curPos.X, curPos.Y+1), maxTry-1)
	}

	// TODO allow creeps to move sideways around the tower? (if so don't allow for player)
	// return false
}

func (c *CreepData) GetScoreValue() int {
	return c.scoreValue
}
