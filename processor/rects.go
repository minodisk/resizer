package processor

import "image"

func RotateRect(r image.Rectangle) image.Rectangle {
	s := r.Size()
	return image.Rectangle{r.Min, image.Point{s.Y, s.X}}
}
