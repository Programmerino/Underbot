package cv

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/pkg/errors"
	"gocv.io/x/gocv"
	"golang.org/x/image/bmp"
)

var (
	colors [][]uint8
)

const (
	/*
		Objects will be colored based on:
		0 - Random colors for each object
		1 - Recognized objects are colored accordingly to what is specified in objectAttrs.go, and unrecognized objects are gray
	*/
	coloring = 1
)

var leniance = 3 // How different the object can be to be detected anyways

var rnd *rand.Rand
var gray = color.RGBA{193, 193, 193, 255}

var objects []Object
var recognizedObjects []Object

func init() {
	rnd = rand.New(rand.NewSource(99))
	for i := 0; i < 300; i++ {
		colors = append(colors, []uint8{randCol(), randCol(), randCol()})
	}
}

func randCol() uint8 {
	return uint8(rnd.Uint32())
}

// Describing an object on screen
type Object struct {
	// Rectangle describing where it is on screen
	Bounds     image.Rectangle
	ID         int
	Color      color.Color
	Recognized bool
	RecogObj   RecognizedObject
}

func newObject(rect image.Rectangle, ID int, color color.Color, recogobj RecognizedObject) Object {
	obj := Object{}
	obj.Bounds = rect
	obj.ID = ID
	obj.Color = color
	obj.Recognized = (recogobj != RecognizedObject{})
	if obj.Recognized {
		obj.RecogObj = recogobj
	}
	return obj
}

// Returns all the objects' rectangles detected
func GetObjects() []Object {
	return objects
}

// Returns all the identified objects
func GetRecognizedObjects() []Object {
	return recognizedObjects
}

// ProcessImage processes image, taking necessary action on the game, and then returns image with eye candy about what the CV sees
func ProcessImage(img image.RGBA) image.Image {
	var srcGray gocv.Mat
	srcGray = gocv.NewMat()
	buf := new(bytes.Buffer)
	err := bmp.Encode(buf, &img)
	if err != nil {
		panic(errors.Wrap(err, "Failed to encode image to bitmap format"))
	}
	src, err := gocv.IMDecode(buf.Bytes(), 1)
	if err != nil {
		panic(errors.Wrap(err, "Failed to decode incoming image for OpenCV"))
	}
	gocv.CvtColor(src, &srcGray, gocv.ColorBGRToGray)
	thresMap := gocv.NewMat()
	defer thresMap.Close()
	gocv.Threshold(srcGray, &thresMap, 50, 255, gocv.ThresholdBinary)
	//  findContours( threshold_output, contours, hierarchy, CV_RETR_TREE, CV_CHAIN_APPROX_SIMPLE, Point(0, 0) );
	contours := gocv.FindContours(thresMap, gocv.RetrievalTree, gocv.ChainApproxSimple)
	//fmt.Printf("Detected %v objects\n", len(contours))
	objects = []Object{}
	recognizedObjects = []Object{}
	for i, object := range contours {
		if len(object) > len(colors) {
			colors = append(colors, []uint8{randCol(), randCol(), randCol()})
		}
		rect := getRectangle(object)
		var dispColor color.Color
		if coloring == 0 {
			dispColor = color.RGBA{colors[i][0], colors[i][1], colors[i][2], 255}
		} else if coloring == 1 {
			dispColor = gray
		}

		objColor := centerColor(img.SubImage(rect))
		// Replace objColor with the actual center color of the object
		obj := newObject(rect, i, objColor, RecognizedObject{})
		obj.recognize()
		if obj.Recognized {
			if (obj.RecogObj.Color == color.RGBA{0, 0, 0, 255}) {
				dispColor = color.RGBA{colors[i][0], colors[i][1], colors[i][2], 255}
			} else {
				dispColor = obj.Color
			}
		}
		drawObject(&img, dispColor, obj)
		objects = append(objects, obj)
	}
	err = srcGray.Close()
	if err != nil {
		panic(errors.Wrap(err, "Failed to close Mat"))
	}
	err = src.Close()
	if err != nil {
		panic(errors.Wrap(err, "Failed to close Mat"))
	}
	return &img
}

func (obj *Object) recognize() {
	size := image.Point{obj.Bounds.Dx(), obj.Bounds.Dy()}
	for _, recogObj := range recognizableObjects {

		// If the object is colored properly
		if (recogObj.Color != color.RGBA{0, 0, 0, 255}) {
			if !pntWithin(recogObj.Size, size, leniance) {
				continue
			} else {
				if recogObj.Color == obj.Color {
					recTreatment(obj, recogObj)
				}
			}
		} else if pntWithin(recogObj.Size, size, leniance) {
			recTreatment(obj, recogObj)
		}
		return
	}
}

// Determines if the two numbers are within a certain range of difference
func within(num1 int, num2 int, rng int) bool {
	return (math.Abs(float64(num2-num1)) <= float64(rng))
}

func pntWithin(pnt1 image.Point, pnt2 image.Point, rng int) bool {
	return ((within(pnt1.X, pnt2.X, rng)) && (within(pnt1.Y, pnt2.Y, rng)))
}

func recTreatment(obj *Object, recogObj RecognizedObject) {
	//fmt.Printf("Recognized %s\n", recogObj.Name)
	obj.Recognized = true
	obj.RecogObj = recogObj
	recognizedObjects = append(recognizedObjects, *obj)
}

func centerColor(tmpImg image.Image) color.Color {
	img, ok := tmpImg.(*image.RGBA)
	if !ok {
		panic("Failed to convert image to RGBA in centerColor")
	}
	centerPoint := rectangleCenter(img.Rect)
	return img.At(centerPoint.X, centerPoint.Y)
}

func rectangleCenter(rect image.Rectangle) (point image.Point) {
	point.X = (rect.Min.X + rect.Max.X) / 2
	point.Y = (rect.Min.Y + rect.Max.Y) / 2
	return point
}

func drawPoint(img *image.RGBA, x, y int, c color.Color) {
	img.Set(x, y, c)
}

func drawObject(img *image.RGBA, color color.Color, obj Object) {
	Rect(img, color, obj.Bounds.Min.X, obj.Bounds.Min.Y, obj.Bounds.Max.X, obj.Bounds.Max.Y)
}

func getRectangle(points []image.Point) image.Rectangle {
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

// VLine draws a veritcal line
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
