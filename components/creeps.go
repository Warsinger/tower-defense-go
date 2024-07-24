package components

import (
	"fmt"
	"image"
	"math/rand"

	"tower-defense/assets"

	"github.com/yohamta/donburi"
)

type CreepData struct {
	scoreValue int
}

var Creep = donburi.NewComponentType[CreepData]()

func NewCreep(w donburi.World, x, y int) error {
	entity := w.Create(Creep, Position, Velocity, Render, Health, Attack)
	entry := w.Entry(entity)
	Position.SetValue(entry, PositionData{x: x, y: y})
	Velocity.SetValue(entry, VelocityData{x: 0, y: 1})
	choose := rand.Intn(2) + 1
	name := fmt.Sprintf("creep%v", choose)
	Render.SetValue(entry, RenderData{primaryRenderer: &SpriteData{image: assets.GetImage(name)}, secondaryRenderer: &RangeRenderData{}})
	Creep.SetValue(entry, CreepData{scoreValue: 10})
	Health.SetValue(entry, HealthData{1})
	Attack.SetValue(entry, AttackData{Power: 5, AttackType: MeleeSingle, Range: 3})
	return nil
}

func (a *CreepData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)

	pos.x += v.x
	pos.y += v.y
	return nil
}

func (a *CreepData) GetRect(entry *donburi.Entry) image.Rectangle {
	sprite := Render.Get(entry)
	return sprite.primaryRenderer.GetRect(entry)
}

func (a *CreepData) GetScoreValue() int {
	return a.scoreValue
}
