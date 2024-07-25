package components

import (
	"fmt"
	"image"
	"math/rand"

	"tower-defense/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
)

type CreepData struct {
	scoreValue int
}
type CreepRenderData struct {
}

var Creep = donburi.NewComponentType[CreepData]()

func NewCreep(w donburi.World, x, y int) error {
	entity := w.Create(Creep, Position, Velocity, Render, Health, Attack)
	entry := w.Entry(entity)
	Position.SetValue(entry, PositionData{x: x, y: y})
	Velocity.SetValue(entry, VelocityData{x: 0, y: 3})
	choose := rand.Intn(2) + 1
	name := fmt.Sprintf("creep%v", choose)
	Render.SetValue(entry, *NewRenderer(&SpriteData{image: assets.GetImage(name)}, &RangeRenderData{}, &CreepRenderData{}))
	Creep.SetValue(entry, CreepData{scoreValue: 10})
	Health.SetValue(entry, HealthData{5})
	Attack.SetValue(entry, AttackData{Power: 5, AttackType: RangedSingle, Range: 50, Cooldown: 15})
	return nil
}

func (c *CreepData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)
	pos.x += v.x
	pos.y += v.y

	a := Attack.Get(entry)
	a.AttackEnemyRange(entry, Tower, nil)

	return nil
}

func (a *CreepData) GetRect(entry *donburi.Entry) image.Rectangle {
	sprite := Render.Get(entry)
	return sprite.GetPrimaryRenderer().GetRect(entry)
}

func (a *CreepData) GetScoreValue() int {
	return a.scoreValue
}

func (cr *CreepRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
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

func (cr *CreepRenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	panic("CreepRenderData.GetRect() unimplemented")
}
