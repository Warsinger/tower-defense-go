package components

import (
	"image"
	"image/color"
	"tower-defense/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
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
	Ticker     int
	AttackType AttackType
}

var Attack = donburi.NewComponentType[AttackData]()
var Health = donburi.NewComponentType[HealthData]()

func (a *AttackData) Update(e *donburi.Entry) error {
	if a.Ticker == 0 {
		// look for a creep in range to shoot at
		creep := a.findCreepInRange(e)
		if creep != nil {
			a.Attack(creep)
			a.Ticker++
		}
	}
	if a.Ticker >= a.Cooldown {
		a.Ticker = 0
	}
	return nil
}

func (a *AttackData) findCreepInRange(e *donburi.Entry) *donburi.Entry {
	// query for creeps then find the closest one
	minDist := float64(a.Range + 1)
	var foundCreep *donburi.Entry = nil
	aRect := a.GetRect(e)
	query := donburi.NewQuery(filter.Contains(Creep))
	query.Each(e.World, func(ce *donburi.Entry) {
		creep := Creep.Get(ce)
		cRect := creep.GetRect(ce)

		dist := util.DistanceRects(aRect, cRect)
		if dist < minDist {
			minDist = dist
			foundCreep = ce
		}
	})
	return foundCreep
}

func (a *AttackData) Attack(e *donburi.Entry) error {
	return nil
}

func (a *AttackData) GetRect(e *donburi.Entry) image.Rectangle {
	r := Render.Get(e)
	rect := r.primaryRenderer.GetRect(e)
	ptRange := image.Pt(a.Range, a.Range)
	rect.Min = rect.Min.Sub(ptRange)
	rect.Max = rect.Max.Add(ptRange)
	return rect
}

func (rr *RangeRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	a := Attack.Get(entry)
	aRect := a.GetRect(entry)
	// get midpoint of rect
	aPt := util.MidpointRect(aRect)

	vector.StrokeCircle(screen, float32(aPt.X), float32(aPt.Y), float32(a.Range+aRect.Dx()/2), 1, color.White, true)
}

func (rr *RangeRenderData) GetRect(e *donburi.Entry) image.Rectangle {
	a := Attack.Get(e)

	return a.GetRect(e)
}
