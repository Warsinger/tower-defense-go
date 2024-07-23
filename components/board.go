package components

import "github.com/yohamta/donburi"

type BoardInfo struct {
	Width, Height int
}

var Board = donburi.NewComponentType[BoardInfo]()

func NewBoard(w donburi.World) (BoardInfo, error) {
	entity := w.Create(Board)
	entry := w.Entry(entity)
	b := BoardInfo{Width: 600, Height: 800}
	Board.SetValue(entry, b)
	return b, nil
}
