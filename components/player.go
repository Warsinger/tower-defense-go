package components

import (
	"fmt"
	"image"
	"image/color"

	"tower-defense/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type PlayerData struct {
	money int
	score int
	dead  bool
}
type PlayerRenderData struct {
}

var Player = donburi.NewComponentType[PlayerData]()

func NewPlayer(w donburi.World) error {
	entity := w.Create(Player, Position, Render, Health)
	entry := w.Entry(entity)

	be := Board.MustFirst(entry.World)
	board := Board.Get(be)

	Position.SetValue(entry, PositionData{x: 0, y: board.Height - yBorderBottom})
	Render.SetValue(entry, *NewRenderer(&SpriteData{image: assets.GetImage("base")}, &PlayerRenderData{}, &InfoRenderData{}))
	Player.SetValue(entry, PlayerData{money: 500})
	Health.SetValue(entry, HealthData{50})
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
			err := p.PlaceTower(entry.World, x, y)

			if err != nil {
				switch err.(type) {
				case *PlacementError:
					fmt.Println(err.Error())
				default:
					return err
				}
			} else {
				p.money -= cost
			}
		} else {
			fmt.Printf("Not enough money for tower cost %v, remaining %v\n", cost, p.money)
		}
	}

	return nil
}

type PlacementError struct {
	message string
}

func (e *PlacementError) Error() string {
	return e.message
}

func (p *PlayerData) PlaceTower(world donburi.World, x, y int) error {
	img := assets.GetImage("tower")
	rect := img.Bounds().Add(image.Pt(x, y))
	boardEntry := Board.MustFirst(world)
	board := Board.Get(boardEntry)
	if !rect.In(board.Bounds()) {
		// TODO sound for invalid operation
		assets.PlaySound("invalid1")
		message := fmt.Sprintf("Invalid tower location %v, %v, image out of bounds", x, y)
		return &PlacementError{message}
	} else {
		collision := DetectCollisions(world, rect, filter.Contains(Player))
		if collision != nil {
			// TODO sound for invalid operation
			assets.PlaySound("invalid2")
			message := fmt.Sprintf("Invalid tower location %v, %v, collision with entity collision", x, y)
			return &PlacementError{message}
		}
	}
	return NewTower(world, x, y)
}

func (p *PlayerData) IsDead() bool {
	return p.dead
}
func (p *PlayerData) GetMoney() int {
	return p.money
}
func (p *PlayerData) AddMoney(money int) {
	p.money += money
}

func (p *PlayerData) Kill() {
	p.dead = true
}

func (p *PlayerRenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	return Render.Get(entry).GetRect(entry)
}

func (pr *PlayerRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	p := Player.Get(entry)
	if p.dead {
		rect := pr.GetRect(entry)
		vector.StrokeLine(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Max.X), float32(rect.Max.Y), 3, color.RGBA{255, 0, 0, 255}, true)
		vector.StrokeLine(screen, float32(rect.Max.X), float32(rect.Min.Y), float32(rect.Min.X), float32(rect.Max.Y), 3, color.RGBA{255, 0, 0, 255}, true)
	}

	// draw player money
	str := fmt.Sprintf("$ %d", p.GetMoney())
	op := &text.DrawOptions{}
	op.GeoM.Translate(TextBorder, TextBorder)
	text.Draw(screen, str, assets.ScoreFace, op)

	// draw score
	be := Board.MustFirst(entry.World)
	board := Board.Get(be)
	halfWidth, _ := float64(board.Width/2), float64(board.Height/2)
	str = fmt.Sprintf("SCORE %05d", p.score)
	op = &text.DrawOptions{}
	x, y := text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(halfWidth-x/2, TextBorder+y)
	text.Draw(screen, str, assets.ScoreFace, op)
}

func (p *PlayerData) GetScore() int {
	return p.score
}
func (p *PlayerData) AddScore(score int) {
	p.score += score
}
