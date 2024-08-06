package components

import (
	"image"

	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/yohamta/donburi"
)

type BoardData struct {
	Width, Height int
}

const (
	yBorderBottom         = 70
	TextBorder    float64 = 5
	SpawnBorder           = 60
)

var Board = donburi.NewComponentType[BoardData]()

func NewBoard(world donburi.World, width, height int) (*BoardData, error) {
	entity := world.Create(Board)
	_ = srvsync.NetworkSync(world, &entity, Board)
	entry := world.Entry(entity)
	board := &BoardData{Width: width, Height: height}
	Board.Set(entry, board)
	return board, nil
}

func (b *BoardData) Bounds() image.Rectangle {
	return image.Rect(0, 0, b.Width, b.Height)
}
