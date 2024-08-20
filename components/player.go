package components

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"tower-defense/assets"
	"tower-defense/config"
	"tower-defense/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type PlayerData struct {
	Money       int
	Score       int
	Dead        bool
	TowerLevels int
}
type PlayerRenderData struct {
}

var Player = donburi.NewComponentType[PlayerData]()
var PlayerRender = donburi.NewComponentType[PlayerRenderData]()

func NewPlayer(world donburi.World, startingTowerLevel int) error {
	entity := world.Create(Player, Position, Health, Attack, SpriteRender, PlayerRender, InfoRender)
	err := srvsync.NetworkSync(world, &entity, Player, Position, Health, Attack, SpriteRender, PlayerRender, InfoRender)
	if err != nil {
		return err
	}
	entry := world.Entry(entity)

	be := Board.MustFirst(entry.World)
	board := Board.Get(be)

	Position.Set(entry, &PositionData{X: 0, Y: board.Height - yBorderBottom})
	Player.Set(entry, &PlayerData{Money: 500, TowerLevels: startingTowerLevel})
	Health.Set(entry, NewHealthData(100))
	Attack.Set(entry, &AttackData{Power: 1, AttackType: RangedSingle, Range: 15, cooldown: util.NewCooldownTimer(10), noLead: true})
	SpriteRender.Set(entry, &SpriteRenderData{Name: "base"})
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

func (p *PlayerData) UserSpeedUpdate(entry *donburi.Entry) error {
	if p.Dead {
		return nil
	}
	config := config.GetConfig(entry.World)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		x, y := ebiten.CursorPosition()
		_, err := p.TryPlaceTower(entry.World, x, y, config.Sound, config.Debug)
		if err != nil {
			return err
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		// find tower below the click and heal it if we have enough money
		x, y := ebiten.CursorPosition()
		towerEntry := findTower(entry.World, x, y)
		if towerEntry != nil {
			_ = p.TryHealTower(towerEntry, config.Sound, config.Debug)

		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		// find tower below the click and upgrade it if we have enough money
		x, y := ebiten.CursorPosition()
		towerEntry := findTower(entry.World, x, y)
		if towerEntry != nil {
			_ = p.TryUpgradeTower(towerEntry, config.Sound, config.Debug)
		}
	}

	return nil
}

func (p *PlayerData) TryHealTower(entry *donburi.Entry, sound, debug bool) bool {
	healed := false
	cost := towerManager.GetHealCost("Ranged")
	if p.Money >= cost {
		tower := Tower.Get(entry)
		if tower.Heal(entry, debug) {
			p.Money -= cost
			healed = true
		}
	} else {
		if debug {
			fmt.Printf("Not enough money to upgrade tower cost %v, remaining %v\n", cost, p.Money)
		}
		if sound {
			assets.PlaySound("invalid2")
		}
	}
	return healed
}

func (p *PlayerData) TryUpgradeTower(entry *donburi.Entry, sound, debug bool) bool {
	upgraded := false
	cost := towerManager.GetUpgradeCost("Ranged")
	if p.Money >= cost {
		tower := Tower.Get(entry)
		if tower.Upgrade(entry, debug) {
			p.Money -= cost
			p.TowerLevels++
			upgraded = true
		}
	} else {
		if debug {
			fmt.Printf("Not enough money to upgrade tower cost %v, remaining %v\n", cost, p.Money)
		}
		if sound {
			assets.PlaySound("invalid2")
		}
	}
	return upgraded
}

func (p *PlayerData) TryPlaceTower(world donburi.World, x, y int, sound, debug bool) (bool, error) {
	placed := false
	cost := towerManager.GetCost("Ranged")
	if p.Money >= cost {
		err := p.PlaceTower(world, x, y, sound)

		if err != nil {
			switch err.(type) {
			case *PlacementError:
				if debug {
					fmt.Println(err.Error())
				}
			default:
				return false, err
			}
		} else {
			p.Money -= cost
			placed = true
		}
	} else {
		if debug {
			fmt.Printf("Not enough money for tower cost %v, remaining %v\n", cost, p.Money)
		}
		if sound {
			assets.PlaySound("invalid2")
		}
	}
	return placed, nil
}

func (p *PlayerData) GameSpeedUpdate(entry *donburi.Entry) error {
	a := Attack.Get(entry)
	a.AttackEnemyRange(entry, nil, Creep)
	return nil
}

type PlacementError struct {
	message string
}

func (e *PlacementError) Error() string {
	return e.message
}

func (p *PlayerData) PlaceTower(world donburi.World, x, y int, sound bool) error {
	img := assets.GetImage("tower")
	bounds := img.Bounds()
	rect := bounds.Add(image.Pt(x-bounds.Dx()/2, y-bounds.Dy()/2))
	boardEntry := Board.MustFirst(world)
	board := Board.Get(boardEntry)
	if !rect.In(board.Bounds()) {
		if sound {
			assets.PlaySound("invalid1")
		}
		message := fmt.Sprintf("Invalid tower location %v, %v, image out of bounds", x, y)
		return &PlacementError{message}
	} else {
		collision := DetectCollisionsWorld(world, rect, filter.Contains(Player))
		if collision != nil {
			if sound {
				assets.PlaySound("invalid2")
			}
			message := fmt.Sprintf("Invalid tower location %v, %v, collision with entity", x, y)
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

func (player *PlayerData) GetCreepLevel() int {
	return int(math.Trunc(float64(player.TowerLevels)/5)) + 1
}

func (pr *PlayerRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry, debug bool) {
	player := Player.Get(entry)
	if player.Dead {
		rect := GetRect(entry)
		vector.StrokeLine(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Max.X), float32(rect.Max.Y), 3, color.RGBA{255, 0, 0, 255}, true)
		vector.StrokeLine(screen, float32(rect.Max.X), float32(rect.Min.Y), float32(rect.Min.X), float32(rect.Max.Y), 3, color.RGBA{255, 0, 0, 255}, true)
	}
	be := Board.MustFirst(entry.World)
	board := Board.Get(be)

	str := fmt.Sprintf("$ %d", player.GetMoney())
	nextY := DrawTextLines(screen, assets.ScoreFace, str, float64(board.Width), TextBorder, text.AlignStart, text.AlignStart)

	str = fmt.Sprintf("Max Tower Level %d", GetMaxTowerLevel(entry.World))
	_ = DrawTextLines(screen, assets.InfoFace, str, float64(board.Width), nextY, text.AlignStart, text.AlignStart)

	str = fmt.Sprintf("SCORE %05d", player.Score)
	_ = DrawTextLines(screen, assets.ScoreFace, str, float64(board.Width), TextBorder, text.AlignCenter, text.AlignStart)

	str = fmt.Sprintf("Creep Level %d", player.GetCreepLevel())
	if debug {
		str = fmt.Sprintf("%s (%d)", str, player.TowerLevels)
	}
	_ = DrawTextLines(screen, assets.InfoFace, str, float64(board.Width), nextY, text.AlignCenter, text.AlignStart)
}

func (p *PlayerData) GetScore() int {
	return p.Score
}
func (p *PlayerData) AddScore(score int) {
	p.Score += score
}
