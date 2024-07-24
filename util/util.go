package util

import (
	"image"
	"math"
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
