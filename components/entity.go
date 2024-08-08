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

var Position = donburi.NewComponentType[PositionData]()
var Velocity = donburi.NewComponentType[VelocityData]()

func DetectCollisionsEntry(entry *donburi.Entry, rect image.Rectangle, excludeFilter filter.LayoutFilter) *donburi.Entry {
	var collision *donburi.Entry = nil
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(SpriteRender, Position),
			filter.Not(excludeFilter),
		),
	)

	query.Each(entry.World, func(testEntry *donburi.Entry) {
		if collision == nil && testEntry.Entity().Id() != entry.Entity().Id() {
			testRect := GetRect(testEntry)
			if rect.Overlaps(testRect) {
				collision = testEntry
			}
		}
	})
	return collision
}

func DetectCollisionsWorld(world donburi.World, rect image.Rectangle, excludeFilter filter.LayoutFilter) *donburi.Entry {
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
