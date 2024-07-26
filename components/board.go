package components

import (
	"image"

	"github.com/yohamta/donburi"
)

type BoardInfo struct {
	Width, Height int
}

const (
	yBorderBottom         = 45
	TextBorder    float64 = 5
	SpawnBorder           = 60
)

var Board = donburi.NewComponentType[BoardInfo]()

func NewBoard(w donburi.World, width, height int) (BoardInfo, error) {
	entity := w.Create(Board)
	entry := w.Entry(entity)
	b := BoardInfo{Width: width, Height: height}
	Board.SetValue(entry, b)
	return b, nil
}

func (b *BoardInfo) Bounds() image.Rectangle {
	return image.Rect(0, 0, b.Width, b.Height)
}
