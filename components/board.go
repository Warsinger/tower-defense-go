package components

import (
	"image"

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

func NewBoard(w donburi.World, width, height int) (BoardData, error) {
	entity := w.Create(Board)
	entry := w.Entry(entity)
	b := BoardData{Width: width, Height: height}
	Board.SetValue(entry, b)
	return b, nil
}

func (b *BoardData) Bounds() image.Rectangle {
	return image.Rect(0, 0, b.Width, b.Height)
}
