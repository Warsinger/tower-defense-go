package components

import "github.com/yohamta/donburi"

type BoardInfo struct {
	Width, Height int
}

const (
	yBorderBottom         = 45
	TextBorder    float64 = 5
)

var Board = donburi.NewComponentType[BoardInfo]()

func NewBoard(w donburi.World, width, height int) (BoardInfo, error) {
	entity := w.Create(Board)
	entry := w.Entry(entity)
	b := BoardInfo{Width: width, Height: height}
	Board.SetValue(entry, b)
	return b, nil
}
