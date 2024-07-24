package components

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"tower-defense/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

type BulletData struct {
	start, end image.Point
	speed      int
	creep      bool
}

var Bullet = donburi.NewComponentType[BulletData]()

type BulletRenderData struct {
	color color.Color
	size  int
}

var creepBulletColor = color.RGBA{255, 0, 0, 255}
var towerBulletColor = color.RGBA{40, 255, 40, 255}

func NewBullet(w donburi.World, start, end image.Point, creep bool) {
	bulletEntity := w.Create(Bullet, Position, Velocity, Render, Attack)
	bullet := w.Entry(bulletEntity)

	Position.SetValue(bullet, PositionData{start.X, start.Y})

	Render.SetValue(bullet, *NewRenderer(NewBulletRender(creep)))
	Attack.SetValue(bullet, AttackData{Power: 1, AttackType: RangedSingle, Range: 1, Cooldown: 30})
	Bullet.SetValue(bullet, BulletData{start: start, end: end, speed: 4, creep: creep})
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
	dist := util.DistancePoints(bd.start, bd.end)
	ratio := dist / float64(bd.speed)
	fmt.Printf("dist: %v, ratio: %v, start: %v, end: %v\n", dist, ratio, bd.start, bd.end)

	newX := pos.x + int(math.Ceil(float64(bd.end.X-bd.start.X)/ratio))
	newY := pos.y + int(math.Ceil(float64(bd.end.Y-bd.start.Y)/ratio))
	fmt.Printf("newX, newY: %v, %v\n", newX, newY)
	be := Board.MustFirst(entry.World)
	board := Board.Get(be)

	if newX < 0 || newX > board.Width || newY < 0 || newY > board.Height {
		entry.Remove()
	} else {
		// if enemy in range, attack it
		a := Attack.Get(entry)
		if bd.IsCreep() {
			a.AttackEnemyIntersect(entry, Tower, nil, AfterBulletAttack)
		} else {
			a.AttackEnemyIntersect(entry, Creep, OnKillCreep, AfterBulletAttack)
		}

		pos.x = newX
		pos.y = newY
	}
	return nil
}
func AfterBulletAttack(bulletEntry *donburi.Entry) {
	bulletEntry.Remove()
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
