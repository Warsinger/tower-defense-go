package game

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"tower-defense/assets"
	comp "tower-defense/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	ecslib "github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type GameData struct {
	world       donburi.World
	ecs         *ecslib.ECS
	gameOver    bool
	paused      bool
	highScore   int
	width       int
	height      int
	speed       int
	debug       bool
	lines       bool
	creepTimer  int
	tickCounter int
}

const minSpeed = 0
const maxSpeed = 60

func NewGame(width, height, speed int, debug, lines bool) (*GameData, error) {
	world := donburi.NewWorld()
	ecs := ecslib.NewECS(world)
	board, err := comp.NewBoard(world, width, height)
	if err != nil {
		return nil, err
	}
	err = assets.LoadAssets()
	if err != nil {
		return nil, err
	}

	highScore := LoadScores()

	ebiten.SetWindowSize(int(board.Width), int(board.Height))
	ebiten.SetWindowTitle("Tower Defense")

	if speed < minSpeed {
		speed = max(1, minSpeed)
	} else if speed > maxSpeed {
		speed = maxSpeed
	}

	return &GameData{
		world:     world,
		ecs:       ecs,
		highScore: highScore,
		width:     width,
		height:    height,
		speed:     speed,
		debug:     debug,
		lines:     lines,
	}, nil
}

const highScoreFile = "score/highscore.txt"

func LoadScores() int {
	var highScore int = 0
	bytes, err := os.ReadFile(highScoreFile)
	if err == nil {
		highScore, err = strconv.Atoi(string(bytes))
		if err != nil {
			fmt.Printf("WARN high score formatting err %v\n", err)
		}
	}

	return highScore
}

func (g *GameData) SaveScores() error {
	str := strconv.Itoa(g.highScore)

	err := os.WriteFile(highScoreFile, []byte(str), 0644)
	return err
}

func (g *GameData) Init() error {
	err := comp.NewPlayer(g.world)
	if err != nil {
		return err
	}

	return nil
}

const maxCreepTimer = 120

func (g *GameData) Clear() error {
	g.gameOver = false
	g.paused = false
	g.creepTimer = maxCreepTimer - 5
	g.tickCounter = 0

	query := donburi.NewQuery(filter.Or(
		filter.Contains(comp.Bullet),
		filter.Contains(comp.Player),
		filter.Contains(comp.Tower),
		filter.Contains(comp.Creep),
	))
	query.Each(g.world, func(e *donburi.Entry) {
		e.Remove()
	})
	return nil
}

func (g *GameData) GetWorld() donburi.World {
	return g.world
}
func (g *GameData) GetECS() *ecslib.ECS {
	return g.ecs
}

func (g *GameData) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Clear()
		g.Init()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if err := g.EndGame(); err != nil {
			return err
		}
		return ebiten.Termination
	}

	if g.gameOver {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.paused = !g.paused
	}
	if g.paused {
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		g.speed = (min(g.speed+5, maxSpeed))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		g.speed = (max(g.speed-5, minSpeed))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	// update player separately from other entities to allow user interactions outside of speed controls
	pe := comp.Player.MustFirst(g.world)
	player := comp.Player.Get(pe)
	err := player.Update(pe)
	if err != nil {
		return err
	}

	if g.speed != 0 && float32(g.tickCounter) > float32(ebiten.TPS())/float32(g.speed) {
		g.tickCounter = 0
		err := g.UpdateEntities()
		if err != nil {
			return err
		}

	} else {
		g.tickCounter++
	}

	g.highScore = max(player.GetScore(), g.highScore)

	return nil
}
func (g *GameData) UpdateEntities() error {
	// query for all entities that have position and velocity and ???
	// and have them do their updates
	query := donburi.NewQuery(
		filter.And(
			filter.Or(
				filter.Contains(comp.Creep),
				filter.Contains(comp.Tower),
				filter.Contains(comp.Bullet),
			),
		),
	)
	var err error = nil
	// update all entities
	query.Each(g.world, func(entry *donburi.Entry) {
		if entry.HasComponent(comp.Creep) {
			creep := comp.Creep.Get(entry)
			err = creep.Update(entry)
			if err != nil {
				return
			}

		}
		if entry.HasComponent(comp.Tower) {
			tower := comp.Tower.Get(entry)
			err = tower.Update(entry)
			if err != nil {
				return
			}

		}

		if entry.HasComponent(comp.Bullet) {
			b := comp.Bullet.Get(entry)
			err = b.Update(entry)
			if err != nil {
				return
			}
		}
	})
	// after updating all positions check for collisions
	// get all the bullets, for each bullet loop through all the creeps (or other objects) and see if there are collisions
	// if there is a collition, remove both objects (or subtract from their health)
	// accumultate points for killing creeps
	err = g.DetectCollisions()
	if err != nil {
		return err
	}

	// check for all creeps destroyed or creeps reaching bottom
	pe := comp.Player.MustFirst(g.world)
	player := comp.Player.Get(pe)
	pRender := comp.Render.Get(pe)
	pRect := pRender.GetRect(pe)
	query = donburi.NewQuery(filter.Contains(comp.Creep))
	query.Each(g.world, func(ae *donburi.Entry) {
		creep := comp.Creep.Get(ae)
		cRect := creep.GetRect(ae)

		if cRect.Max.Y >= pRect.Min.Y {
			player.Kill()
			g.EndGame()
		}
	})
	g.creepTimer++
	if g.creepTimer >= maxCreepTimer {
		g.SpawnCreeps()
		g.creepTimer = 0
	}

	return err
}

const muiltiSpawnChance = 0.8

func (g *GameData) SpawnCreeps() {
	var count = 1
	if rand.Float32() > muiltiSpawnChance {
		count = 2
	}
	for i := 0; i < count; i++ {
		be := comp.Board.MustFirst(g.world)
		board := comp.Board.Get(be)
		const border = 60
		x := rand.Intn(board.Width-border) + border/2
		y := border
		comp.NewCreep(g.world, x, y)
	}
}

func (g *GameData) EndGame() error {
	assets.PlaySound("killed")
	g.gameOver = true
	if err := g.SaveScores(); err != nil {
		return err
	}
	return nil
}

func (g *GameData) CleanBoard() {
	query := donburi.NewQuery(filter.Or(
		filter.Contains(comp.Bullet),
		filter.Contains(comp.Creep),
	))

	query.Each(g.world, func(be *donburi.Entry) {
		be.Remove()
	})
}

func (g *GameData) DetectCollisions() error {
	var err error = nil
	query := donburi.NewQuery(filter.Contains(comp.Bullet))
	query.Each(g.world, func(bulletEntry *donburi.Entry) {
		bulletRender := comp.Render.Get(bulletEntry)
		bulletRect := bulletRender.GetRect(bulletEntry)
		bullet := comp.Bullet.Get(bulletEntry)

		query := donburi.NewQuery(filter.Or(
			filter.Contains(comp.Creep),
		))
		query.Each(g.world, func(e *donburi.Entry) {
			pe := comp.Player.MustFirst(g.world)
			pRender := comp.Render.Get(pe)
			player := comp.Player.Get(pe)
			if bullet.IsCreep() {
				playerRect := pRender.GetRect(pe)
				if bulletRect.Overlaps(playerRect) {
					player.Kill()
					err = g.EndGame()
					if err != nil {
						return
					}
				}
			} else if e.HasComponent(comp.Creep) {
				creep := comp.Creep.Get(e)
				creepRect := creep.GetRect(e)
				if bulletRect.Overlaps(creepRect) {
					player.AddScore(creep.GetScoreValue())

					// remove bullet and creep
					e.Remove()
					bulletEntry.Remove()
					assets.PlaySound("explosion")
				}
			}
		})
	})

	return err
}

func (g *GameData) Draw(screen *ebiten.Image) {
	screen.Clear()

	img := assets.GetImage("backgroundV")
	opts := &ebiten.DrawImageOptions{}
	screen.DrawImage(img, opts)

	// query for all entities
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Position, comp.Render),
		),
	)

	// draw all entities
	query.Each(g.world, func(entry *donburi.Entry) {
		r := comp.Render.Get(entry)
		r.Draw(screen, entry)
	})

	g.DrawText(screen)
}

func (g *GameData) DrawText(screen *ebiten.Image) {
	be := comp.Board.MustFirst(g.world)
	board := comp.Board.Get(be)
	halfWidth, halfHeight := float64(board.Width/2), float64(board.Height/2)

	// draw high score
	str := fmt.Sprintf("HIGH %05d", g.highScore)
	op := &text.DrawOptions{}
	x, _ := text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(float64(board.Width)-x-comp.TextBorder, comp.TextBorder)
	text.Draw(screen, str, assets.ScoreFace, op)

	if g.gameOver {
		// draw game over
		str := "GAME OVER"
		op := &text.DrawOptions{}
		x, y := text.Measure(str, assets.ScoreFace, op.LineSpacing)
		op.GeoM.Translate(halfWidth-x/2, halfHeight-y/2)
		text.Draw(screen, str, assets.ScoreFace, op)
	} else if g.paused {
		// draw paused
		str := "PAUSED"
		op := &text.DrawOptions{}
		x, y := text.Measure(str, assets.ScoreFace, op.LineSpacing)
		op.GeoM.Translate(halfWidth-x/2, halfHeight-y/2)
		text.Draw(screen, str, assets.ScoreFace, op)
	}

	if g.debug {
		str := fmt.Sprintf("Speed %v\nTPS %2.1f", g.speed, ebiten.ActualTPS())
		ebitenutil.DebugPrintAt(screen, str, 5, 50)
	}
}

func (g *GameData) Layout(width, height int) (int, int) {
	return width, height
}
