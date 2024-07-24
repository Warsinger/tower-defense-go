package components

import (
	"image"
	"image/color"
	"tower-defense/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
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
	ticker     int
	AttackType AttackType
}

type RangeRenderData struct {
}

var Attack = donburi.NewComponentType[AttackData]()
var Health = donburi.NewComponentType[HealthData]()

func (a *AttackData) GetRect(e *donburi.Entry) image.Rectangle {
	r := Render.Get(e)
	rect := r.GetPrimaryRenderer().GetRect(e)
	ptRange := image.Pt(a.Range, a.Range)
	rect.Min = rect.Min.Sub(ptRange)
	rect.Max = rect.Max.Add(ptRange)
	return rect
}

func (a *AttackData) GetTicker() int {
	return a.ticker
}
func (a *AttackData) IncrementTicker() {
	a.ticker++
}

func (a *AttackData) CheckCooldown() {
	if a.ticker >= a.Cooldown {
		a.ticker = 0
	}
}

func (rr *RangeRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	a := Attack.Get(entry)
	aRect := a.GetRect(entry)
	aPt := util.MidpointRect(aRect)

	vector.StrokeCircle(screen, float32(aPt.X), float32(aPt.Y), float32(aRect.Dx()/2), 1, color.White, true)
}

func (rr *RangeRenderData) GetRect(e *donburi.Entry) image.Rectangle {
	a := Attack.Get(e)

	return a.GetRect(e)
}
