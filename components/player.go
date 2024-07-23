package components

import (
	"fmt"
	"image"
	"image/color"

	"tower-defense/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

type PlayerData struct {
	money int
	dead  bool
}

var Player = donburi.NewComponentType[PlayerData]()

func NewPlayer(w donburi.World) error {
	entity := w.Create(Player, Position, Render, Health)
	entry := w.Entry(entity)

	be := Board.MustFirst(entry.World)
	board := Board.Get(be)

	Position.SetValue(entry, PositionData{x: board.Width / 2, y: board.Height - yBorderBottom})
	Render.SetValue(entry, RenderData{&SpriteData{image: assets.GetImage("base")}})
	Player.SetValue(entry, PlayerData{money: 500})
	Health.SetValue(entry, HealthData{500})
	return nil
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (p *PlayerData) Update(entry *donburi.Entry) error {
	if p.dead {
		return nil
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		cost := towerManager.GetCost("Ranged")
		if p.money >= cost {
			x, y := ebiten.CursorPosition()
			fmt.Printf("left placing tower at %v, %v, cost %v\n", x, y, cost)
			err := p.PlaceTower(entry.World, x, y)
			if err != nil {
				return err
			}
			p.money -= cost
		} else {
			fmt.Printf("Not enough money for tower cost %v, remaining %v\n", cost, p.money)
		}
	}

	return nil
}

func (p *PlayerData) PlaceTower(w donburi.World, x, y int) error {
	towerEntity := w.Create(Tower, Position, Render, Health, Attack)
	tower := w.Entry(towerEntity)

	Position.SetValue(tower, PositionData{x, y})
	Health.SetValue(tower, HealthData{50})
	Render.SetValue(tower, RenderData{&SpriteData{image: assets.GetImage("tower")}})
	Attack.SetValue(tower, AttackData{Power: 1, AttackType: RangedSingle})

	return nil
}

func (p *PlayerData) IsDead() bool {
	return p.dead
}
func (p *PlayerData) GetMoney() int {
	return p.money
}

func (p *PlayerData) Kill() {
	p.dead = true
}

func (p *PlayerData) GetRect(entry *donburi.Entry) image.Rectangle {
	sprite := Render.Get(entry)
	return sprite.renderer.GetRect(entry)
}

func (p *PlayerData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	if p.dead {
		rect := p.GetRect(entry)
		vector.StrokeLine(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Max.X), float32(rect.Max.Y), 3, color.RGBA{255, 0, 0, 255}, true)
		vector.StrokeLine(screen, float32(rect.Max.X), float32(rect.Min.Y), float32(rect.Min.X), float32(rect.Max.Y), 3, color.RGBA{255, 0, 0, 255}, true)
	}
}
