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
	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type PlayerData struct {
	Money int
	Score int
	Dead  bool
}
type PlayerRenderData struct {
}

var Player = donburi.NewComponentType[PlayerData]()
var PlayerRender = donburi.NewComponentType[PlayerRenderData]()

func NewPlayer(world donburi.World) error {
	entity := world.Create(Player, Position, Health, SpriteRender, PlayerRender, InfoRender, NameComponent)
	err := srvsync.NetworkSync(world, &entity, Player, Position, Health, SpriteRender, PlayerRender, InfoRender, NameComponent)
	if err != nil {
		return err
	}
	entry := world.Entry(entity)

	be := Board.MustFirst(entry.World)
	board := Board.Get(be)

	Position.Set(entry, &PositionData{X: 0, Y: board.Height - yBorderBottom})
	Player.Set(entry, &PlayerData{Money: 500})
	Health.Set(entry, NewHealthData(50))
	name := Name("base")
	NameComponent.Set(entry, &name)
	SpriteRender.Set(entry, &SpriteRenderData{})
	PlayerRender.Set(entry, &PlayerRenderData{})
	InfoRender.Set(entry, &InfoRenderData{})

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
	if p.Dead {
		return nil
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		cost := towerManager.GetCost("Ranged")
		if p.Money >= cost {
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
				p.Money -= cost
			}
		} else {
			// TODO move into place tower
			fmt.Printf("Not enough money for tower cost %v, remaining %v\n", cost, p.Money)
			assets.PlaySound("invalid2")
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
	bounds := img.Bounds()
	rect := bounds.Add(image.Pt(x-bounds.Dx()/2, y-bounds.Dy()/2))
	boardEntry := Board.MustFirst(world)
	board := Board.Get(boardEntry)
	if !rect.In(board.Bounds()) {
		assets.PlaySound("invalid1")
		message := fmt.Sprintf("Invalid tower location %v, %v, image out of bounds", x, y)
		return &PlacementError{message}
	} else {
		collision := DetectCollisions(world, rect, filter.Contains(Player))
		if collision != nil {
			assets.PlaySound("invalid2")
			message := fmt.Sprintf("Invalid tower location %v, %v, collision with entity collision", x, y)
			return &PlacementError{message}
		}
	}
	return NewTower(world, rect.Min.X, rect.Min.Y)
}

func (p *PlayerData) IsDead() bool {
	return p.Dead
}
func (p *PlayerData) GetMoney() int {
	return p.Money
}
func (p *PlayerData) AddMoney(money int) {
	p.Money += money
}

func (p *PlayerData) Kill() {
	p.Dead = true
}

func (pr *PlayerRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	p := Player.Get(entry)
	if p.Dead {
		rect := GetRect(entry)
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
	str = fmt.Sprintf("SCORE %05d", p.Score)
	op = &text.DrawOptions{}
	x, y := text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(halfWidth-x/2, TextBorder+y)
	text.Draw(screen, str, assets.ScoreFace, op)
}

func (p *PlayerData) GetScore() int {
	return p.Score
}
func (p *PlayerData) AddScore(score int) {
	p.Score += score
}
