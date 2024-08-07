package scenes

import (
	"fmt"
	"math/rand/v2"

	"tower-defense/assets"
	comp "tower-defense/components"
	"tower-defense/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type BattleScene struct {
	world           donburi.World
	highScore       int
	width           int
	height          int
	speed           int
	creepTimer      int
	tickCounter     int
	config          *config.ConfigData
	battleState     *BattleSceneState
	endGameCallback func(int) error
}

type BattleSceneState struct {
	GameOver bool
	Paused   bool
}

var BattleState = donburi.NewComponentType[BattleSceneState]()

const minSpeed = 0
const maxSpeed = 60
const maxCreepTimer = 90
const startCreepTimer = 30

func NewBattleScene(world donburi.World, width, height, speed, highScore int, debug bool, endGameCallback func(int) error) (*BattleScene, error) {
	_, err := comp.NewBoard(world, width, height)
	if err != nil {
		return nil, err
	}
	bss := &BattleSceneState{}

	if speed < minSpeed {
		speed = max(1, minSpeed)
	} else if speed > maxSpeed {
		speed = maxSpeed
	}

	return &BattleScene{
		world:           world,
		highScore:       highScore,
		width:           width,
		height:          height,
		speed:           speed,
		creepTimer:      maxCreepTimer - startCreepTimer,
		config:          config.NewConfig(world, debug),
		battleState:     bss,
		endGameCallback: endGameCallback,
	}, nil
}

func (b *BattleScene) Init() error {
	b.Clear()

	entity := b.world.Create(BattleState)
	err := srvsync.NetworkSync(b.world, &entity, BattleState)
	if err != nil {
		return err
	}
	BattleState.Set(b.world.Entry(entity), b.battleState)

	err = comp.NewPlayer(b.world)
	if err != nil {
		return err
	}

	return nil
}

func (b *BattleScene) Clear() error {
	b.battleState.GameOver = false
	b.battleState.Paused = false
	b.creepTimer = maxCreepTimer - startCreepTimer
	b.tickCounter = 0

	query := donburi.NewQuery(filter.Or(
		filter.Contains(comp.Bullet),
		filter.Contains(comp.Player),
		filter.Contains(comp.Tower),
		filter.Contains(comp.Creep),
		filter.Contains(BattleState),
	))
	query.Each(b.world, func(e *donburi.Entry) {
		e.Remove()
	})
	return nil
}

func (b *BattleScene) Update() error {
	pe := comp.Player.MustFirst(b.world)
	player := comp.Player.Get(pe)

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		b.endGameCallback(player.GetScore())
	}

	if b.battleState.GameOver {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		b.speed = (min(b.speed+5, maxSpeed))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		b.speed = (max(b.speed-5, minSpeed))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		b.config.SetGridLines(!b.config.IsGridLines())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		b.config.SetDebug(!b.config.IsDebug())
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		b.battleState.Paused = !b.battleState.Paused
	}

	if b.battleState.Paused {
		return nil
	}

	// update player separately from other entities to allow user interactions outside of speed controls
	err := player.UserSpeedUpdate(pe)
	if err != nil {
		return err
	}

	if b.speed != 0 && float32(b.tickCounter) > float32(ebiten.TPS())/float32(b.speed) {
		b.tickCounter = 0
		err := b.UpdateEntities()
		if err != nil {
			return err
		}
		// have player attack at game speed
		err = player.GameSpeedUpdate(pe)
		if err != nil {
			return err
		}

	} else {
		b.tickCounter++
	}

	b.highScore = max(player.GetScore(), b.highScore)

	return nil
}
func (b *BattleScene) UpdateEntities() error {
	// query for all entities and have them do their updates
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
	entries := make([]*donburi.Entry, 0, query.Count(b.world))
	// update all entities
	query.Each(b.world, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})

	for _, entry := range entries {
		if entry.HasComponent(comp.Creep) {
			creep := comp.Creep.Get(entry)
			err = creep.Update(entry)
			if err != nil {
				return err
			}

		}
		if entry.HasComponent(comp.Tower) {
			tower := comp.Tower.Get(entry)
			err = tower.Update(entry)
			if err != nil {
				return err
			}

		}

		if entry.HasComponent(comp.Bullet) {
			b := comp.Bullet.Get(entry)
			err = b.Update(entry)
			if err != nil {
				return err
			}
		}
	}
	// if the player's health drops to 0 then it is dead and the game is over
	pe := comp.Player.MustFirst(b.world)
	player := comp.Player.Get(pe)
	playerHealth := comp.Health.Get(pe)
	if playerHealth.Health <= 0 {
		player.Kill()
		b.End()
	}

	b.creepTimer++
	if b.creepTimer >= maxCreepTimer {
		b.SpawnCreeps(player.GetCreepLevel())
		b.creepTimer = 0
	}

	return err
}

type creeepSpawnChance struct {
	count  int
	chance float32
}

func (b *BattleScene) SpawnCreeps(creepLevel int) error {
	levelBump := float32(creepLevel) / 20

	spawnChance := []creeepSpawnChance{
		{8, -0.4},
		{7, -0.25},
		{6, -0.1},
		{5, 0.1},
		{4, .25},
		{3, 0.5},
		{2, 0.7}}

	val := rand.Float32() - levelBump
	var count = 1
	for _, spawnChance := range spawnChance {
		if val < spawnChance.chance {
			count = spawnChance.count
			break
		}
	}

	for i := 0; i < count; i++ {
		be := comp.Board.MustFirst(b.world)
		board := comp.Board.Get(be)

		x := rand.IntN(board.Width/count) + board.Width/count*(i)
		y := comp.SpawnBorder
		if x < comp.SpawnBorder {
			x = comp.SpawnBorder
		} else if x > board.Width-comp.SpawnBorder {
			x = board.Width - comp.SpawnBorder
		}
		_, err := comp.NewCreep(b.world, x, y, creepLevel)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BattleScene) End() {
	assets.PlaySound("killed")
	b.battleState.GameOver = true
}

func (b *BattleScene) Draw(screen *ebiten.Image) {
	comp.DrawBoard(screen, b.world, b.config, b.DrawText)
}

func (b *BattleScene) DrawText(screen *ebiten.Image) {
	be := comp.Board.MustFirst(b.world)
	board := comp.Board.Get(be)
	halfWidth, halfHeight := float64(board.Width/2), float64(board.Height/2)

	// draw high score
	str := fmt.Sprintf("HIGH %05d", b.highScore)
	op := &text.DrawOptions{}
	x, _ := text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(float64(board.Width)-x-comp.TextBorder, comp.TextBorder)
	text.Draw(screen, str, assets.ScoreFace, op)

	b.battleState.Draw(screen, halfWidth, halfHeight)

	if b.config.IsDebug() {
		str := fmt.Sprintf("Speed %v\nTPS %2.1f", b.speed, ebiten.ActualTPS())
		ebitenutil.DebugPrintAt(screen, str, 5, 50)
	}
}

func (bss *BattleSceneState) Draw(screen *ebiten.Image, halfWidth, halfHeight float64) {
	if bss.GameOver {
		// draw game over
		str := "GAME OVER"
		op := &text.DrawOptions{}
		x, y := text.Measure(str, assets.ScoreFace, op.LineSpacing)
		op.GeoM.Translate(halfWidth-x/2, halfHeight-y/2)
		text.Draw(screen, str, assets.ScoreFace, op)
		str = "Press R to reset game"
		op = &text.DrawOptions{}
		x, _ = text.Measure(str, assets.InfoFace, op.LineSpacing)
		op.GeoM.Translate(halfWidth-x/2, halfHeight+y/2)
		text.Draw(screen, str, assets.InfoFace, op)
	} else if bss.Paused {
		// draw paused
		str := "PAUSED"
		op := &text.DrawOptions{}
		x, y := text.Measure(str, assets.ScoreFace, op.LineSpacing)
		op.GeoM.Translate(halfWidth-x/2, halfHeight-y/2)
		text.Draw(screen, str, assets.ScoreFace, op)
	}
}
