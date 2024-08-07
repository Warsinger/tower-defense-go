package components

import (
	"image"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

// Component is any struct that holds some kind of data.
type PositionData struct {
	X, Y int
}

type VelocityData struct {
	X, Y    int
	blocked bool
}

type Name string

var Position = donburi.NewComponentType[PositionData]()
var Velocity = donburi.NewComponentType[VelocityData]()
var NameComponent = donburi.NewComponentType[Name]()

func DetectCollisions(world donburi.World, rect image.Rectangle, excludeFilter filter.LayoutFilter) *donburi.Entry {
	var collision *donburi.Entry = nil
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(SpriteRender, Position),
			filter.Not(excludeFilter),
		),
	)

	query.Each(world, func(testEntry *donburi.Entry) {
		if collision == nil {
			testRect := GetRect(testEntry)
			if rect.Overlaps(testRect) {
				collision = testEntry
			}
		}
	})
	return collision
}
