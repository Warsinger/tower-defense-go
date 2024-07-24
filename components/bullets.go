package components

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

type BulletData struct {
	creep bool
}

var Bullet = donburi.NewComponentType[BulletData]()

type BulletRenderData struct {
	color color.Color
	size  int
}

var creepBulletColor = color.RGBA{255, 0, 0, 255}
var towerBulletColor = color.RGBA{255, 215, 0, 255}

func NewBullet(w donburi.World, x, y int, creep bool) {
	bulletEntity := w.Create(Bullet, Position, Velocity, Render, Attack)
	bullet := w.Entry(bulletEntity)

	Position.SetValue(bullet, PositionData{x, y})
	Velocity.SetValue(bullet, VelocityData{x: 4, y: 4})

	Render.SetValue(bullet, *NewRenderer(NewBulletRender(creep)))
	Attack.SetValue(bullet, AttackData{Power: 1, AttackType: RangedSingle, Range: 40, Cooldown: 30})
}

func NewBulletRender(creep bool) *BulletRenderData {
	var color color.Color
	var size int
	if creep {
		color = creepBulletColor
		size = 3
	} else {
		color = towerBulletColor
		size = 4
	}
	return &BulletRenderData{color: color, size: size}
}

func (bd *BulletData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)
	newX := pos.x + v.x
	newY := pos.y + v.y
	be := Board.MustFirst(entry.World)
	board := Board.Get(be)

	if newX < 0 || newY > board.Width || newY < 0 || newY > board.Height {
		entry.Remove()
	} else {
		pos.y = newY
	}

	// if enemy in range, attack it
	a := Attack.Get(entry)

	if bd.IsCreep() {
		a.AttackEnemyIntersect(entry, Tower, nil, nil)
	} else {
		tower := Tower.Get(entry)
		a.AttackEnemyIntersect(entry, Creep, tower.OnKillEnemy, tower.AfterAttackEnemy)
	}

	return nil
}

func (brd *BulletRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	pos := Position.Get(entry)
	vector.DrawFilledCircle(screen, float32(pos.x), float32(pos.y), float32(brd.size), brd.color, true)
}

func (brd *BulletRenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	pos := Position.Get(entry)
	return image.Rect(pos.x, pos.y, pos.x+brd.size, pos.y+brd.size)
}

func (b *BulletData) IsCreep() bool {
	return b.creep
}
