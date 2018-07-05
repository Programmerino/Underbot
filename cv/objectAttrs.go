package cv

import (
	"image"
	"image/color"
)

type RecognizedObject struct {
	Name  string
	Size  image.Point
	Color color.Color
}

func newSpecs(name string, width, height int, color color.Color) RecognizedObject {
	return RecognizedObject{name, image.Point{width, height}, color}
}

var recognizableObjects = []RecognizedObject{
	newSpecs("heart", 15, 15, color.RGBA{255, 0, 0, 255}),
	newSpecs("narratorBox", 607, 145, color.RGBA{0, 0, 0, 255}),
}
