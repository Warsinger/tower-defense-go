package util

import (
	"image"
	"math"

	"github.com/yohamta/donburi/component"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/exp/constraints"
)

func DistanceRects(rect1, rect2 image.Rectangle) float64 {
	pt1 := MidpointRect(rect1)
	pt2 := MidpointRect(rect2)
	return DistancePoints(pt1, pt2)
}

func DistancePoints(pt1, pt2 image.Point) float64 {
	pt3 := pt1.Sub(pt2)
	return math.Sqrt(math.Pow(float64(pt3.X), 2) + math.Pow(float64(pt3.Y), 2))
}

func MidpointRect(rect image.Rectangle) image.Point {
	midX := (rect.Max.X + rect.Min.X) / 2
	midY := (rect.Max.Y + rect.Min.Y) / 2
	return image.Pt(midX, midY)
}

func Abs[T constraints.Integer | constraints.Float](n T) T {
	if n < T(0) {
		return -n
	}
	return n
}

func CreateOrFilter(compTypes ...component.IComponentType) filter.LayoutFilter {
	count := len(compTypes)
	if count > 1 {
		filters := make([]filter.LayoutFilter, count)
		for i := 0; i < count; i++ {
			filters[i] = filter.Contains(compTypes[i])
		}
		return filter.Or(filters...)
	} else if count == 1 {
		return filter.Contains(compTypes[0])
	}
	return filter.Or()
}
