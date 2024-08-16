package components

import (
	"image"
	"image/color"
	"tower-defense/config"
	"tower-defense/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/yohamta/donburi"
)

type BulletData struct {
	start, end image.Point
	speed      int
	creep      bool // TODO switch to using a component tag EnemyTag https://pkg.go.dev/github.com/yohamta/donburi@v1.4.4#readme-tags
}

type BulletRenderData struct {
	Size       int
	R, G, B, A uint8
}

var Bullet = donburi.NewComponentType[BulletData]()
var BulletRender = donburi.NewComponentType[BulletRenderData]()

var creepBulletColor = color.RGBA{255, 0, 0, 255}
var towerBulletColor = color.RGBA{40, 255, 40, 255}

func NewBullet(world donburi.World, start, end image.Point, speed int, creep bool) (*donburi.Entry, error) {
	bulletEntity := world.Create(Bullet, Position, Velocity, Attack, BulletRender)
	err := srvsync.NetworkSync(world, &bulletEntity, Bullet, Position, Attack, BulletRender)
	if err != nil {
		return nil, err
	}
	bullet := world.Entry(bulletEntity)

	Position.Set(bullet, &PositionData{start.X, start.Y})
	Velocity.Set(bullet, &VelocityData{X: 6, Y: 6})

	var color color.RGBA
	var size int
	if creep {
		color = creepBulletColor
		size = 3
	} else {
		color = towerBulletColor
		size = 4
	}
	BulletRender.Set(bullet, NewBulletRender(size, color))
	Attack.Set(bullet, &AttackData{Power: 1, AttackType: RangedSingle, Range: 1, cooldown: util.NewCooldownTimer(30)})
	Bullet.Set(bullet, &BulletData{start: start, end: end, speed: speed, creep: creep})
	return bullet, nil
}

func NewBulletRender(size int, clr color.Color) *BulletRenderData {
	brd := BulletRenderData{Size: size}
	r, g, b, a := clr.RGBA()
	brd.R, brd.G, brd.B, brd.A = uint8(r), uint8(g), uint8(b), uint8(a)
	return &brd
}

func (brd *BulletRenderData) GetColor() color.Color {
	return color.RGBA{brd.R, brd.G, brd.B, brd.A}
}

func (bd *BulletData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	dist := util.DistancePoints(bd.start, bd.end)
	ratio := dist / float64(bd.speed)
	// fmt.Printf("dist: %v, ratio: %v, start: %v, end: %v\n", dist, ratio, bd.start, bd.end)

	newX := pos.X + int(float64(bd.end.X-bd.start.X)/ratio)
	newY := pos.Y + int(float64(bd.end.Y-bd.start.Y)/ratio)
	// fmt.Printf("newX, newY: %v, %v\n", newX, newY)
	be := Board.MustFirst(entry.World)
	board := Board.Get(be)

	if newX < 0 || newX > board.Width || newY < 0 || newY > board.Height {
		entry.Remove()
	} else {
		// if enemy in range, attack it
		a := Attack.Get(entry)
		if bd.IsCreep() {
			a.AttackEnemyIntersect(entry, nil, AfterBulletAttack, Tower, Player)
		} else {
			a.AttackEnemyIntersect(entry, OnKillCreep, AfterBulletAttack, Creep)
		}

		pos.X = newX
		pos.Y = newY
	}
	return nil
}
func AfterBulletAttack(bulletEntry *donburi.Entry) {
	bulletEntry.Remove()
}

func (brd *BulletRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	pos := Position.Get(entry)
	bullet := Bullet.Get(entry)
	color := brd.GetColor()
	vector.DrawFilledCircle(screen, float32(pos.X), float32(pos.Y), float32(brd.Size), color, true)

	config := config.GetConfig(entry.World)
	if config.Debug {
		vector.StrokeLine(screen, float32(bullet.start.X), float32(bullet.start.Y), float32(bullet.end.X), float32(bullet.end.Y), 1, color, true)
	}
}

func (brd *BulletRenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	pos := Position.Get(entry)
	return image.Rect(pos.X, pos.Y, pos.X+brd.Size, pos.Y+brd.Size)
}

func (b *BulletData) IsCreep() bool {
	return b.creep
}
