package rect

import (
	"image"
	"image/color"

	"gitlab.com/256/Underbot/cv/object"
)

// CenterColor gets the color of the pixel in the middle of an image
func CenterColor(tmpImg image.Image) color.Color {
	// Gets the underlying type of RGBA which supports At()
	img, ok := tmpImg.(*image.RGBA)
	if !ok {
		panic("Failed to convert image to RGBA in centerColor")
	}
	centerPoint := rectangleCenter(img.Rect)
	return img.At(centerPoint.X, centerPoint.Y)
}

// rectangleCenter finds the point in the middle of a rectangle
func rectangleCenter(rect image.Rectangle) (point image.Point) {
	point.X = (rect.Min.X + rect.Max.X) / 2
	point.Y = (rect.Min.Y + rect.Max.Y) / 2
	return point
}

// DrawObject draws a rectangle around the object using a specified color onto an image
func DrawObject(img *image.RGBA, color color.Color, obj object.Object) {
	Rect(img, color, obj.Bounds.Min.X, obj.Bounds.Min.Y, obj.Bounds.Max.X, obj.Bounds.Max.Y)
}

// GetRectangle creates a rectangle around a cluster of points
func GetRectangle(points []image.Point) image.Rectangle {
	topLeftX := 0
	topLeftY := 0
	bottomRightX := 0
	bottomRightY := 0
	for i, point := range points {
		if i == 0 {
			bottomRightX = point.X
			bottomRightY = point.Y
			topLeftX = point.X
			topLeftY = point.Y
		}
		if point.X > bottomRightX {
			bottomRightX = point.X
		}
		if point.X < topLeftX {
			topLeftX = point.X
		}
		if point.Y < topLeftY {
			topLeftY = point.Y
		}
		if point.Y > bottomRightY {
			bottomRightY = point.Y
		}
	}
	return image.Rect(topLeftX, topLeftY, bottomRightX, bottomRightY)
}

// HLine draws a horizontal line
func HLine(img *image.RGBA, col color.Color, x1, y, x2 int) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// VLine draws a vertical line
func VLine(img *image.RGBA, col color.Color, x, y1, y2 int) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func Rect(img *image.RGBA, col color.Color, x1, y1, x2, y2 int) {
	HLine(img, col, x1, y1, x2)
	HLine(img, col, x1, y2, x2)
	VLine(img, col, x1, y1, y2)
	VLine(img, col, x2, y1, y2)
}
