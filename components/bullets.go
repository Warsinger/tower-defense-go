package components

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

type BulletData struct {
	length, width int
	creep         bool
}

var Bullet = donburi.NewComponentType[BulletData]()

type BulletRenderData struct {
	color color.Color
}

func (bd *BulletData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)
	newY := pos.y + v.y
	be := Board.MustFirst(entry.World)
	board := Board.Get(be)

	if newY < 0 || newY > board.Height {
		entry.Remove()
	} else {
		pos.y = newY
	}
	return nil
}

func (brd *BulletRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	b := Bullet.Get(entry)
	pos := Position.Get(entry)
	vector.StrokeLine(screen, float32(pos.x), float32(pos.y), float32(pos.x), float32(pos.y+b.length), float32(b.width), brd.color, true)
}

func (brd *BulletRenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	pos := Position.Get(entry)
	b := Bullet.Get(entry)
	return image.Rect(pos.x, pos.y, pos.x+b.width, pos.y+b.length)
}

func (b *BulletData) IsCreep() bool {
	return b.creep
}
