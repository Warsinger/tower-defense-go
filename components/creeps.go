package components

import (
	"fmt"
	"image"
	"math/rand"

	"tower-defense/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
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
	choose := rand.Intn(2) + 1
	Velocity.SetValue(entry, VelocityData{x: 0, y: 2 + choose})
	name := fmt.Sprintf("creep%v", choose)
	Render.SetValue(entry, *NewRenderer(&SpriteData{image: assets.GetImage(name)}, &RangeRenderData{}, &CreepRenderData{}))
	Creep.SetValue(entry, CreepData{scoreValue: 10 * choose})
	Health.SetValue(entry, HealthData{Health: 1 + 2*choose})
	Attack.SetValue(entry, AttackData{Power: 2 + 2*choose, AttackType: RangedSingle, Range: 10 + 10*choose, Cooldown: 5 + 5*choose})
	return nil
}

func (c *CreepData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)
	newPt := image.Pt(v.x, v.y)
	// check whether there are any collisions in the new spot
	collision := false
	rect := c.GetRect(entry)
	fmt.Printf("old rect %v\n", rect)
	rect = rect.Add(newPt)
	fmt.Printf("new rect %v\n", rect)
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Render, Position),
			filter.Not(
				filter.Or(
					filter.Contains(Creep),
					filter.Contains(Bullet),
					filter.Contains(Player),
				),
			),
		),
	)
	query.Each(entry.World, func(testEntry *donburi.Entry) {
		if !collision {
			testRect := Render.Get(testEntry).GetRect(testEntry)
			fmt.Printf("testing overlap of %v with %v\n", rect, testRect)
			if rect.Overlaps(testRect) {
				fmt.Printf("overlap of %v with %v\n", rect, testRect)
				collision = true
			}
		}
	})
	if !collision {
		pos.x += newPt.X
		pos.y += newPt.Y
	}
	v.blocked = collision

	a := Attack.Get(entry)
	a.AttackEnemyRange(entry, Tower, nil)

	return nil
}

func (a *CreepData) GetRect(entry *donburi.Entry) image.Rectangle {
	return Render.Get(entry).GetRect(entry)
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
