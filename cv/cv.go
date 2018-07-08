package cv

import (
	"bytes"
	"image"
	"image/color"
	"math/rand"
	"time"

	"gitlab.com/256/Underbot/cv/num"
	"gitlab.com/256/Underbot/cv/params"
	"gitlab.com/256/Underbot/cv/rect"

	"gitlab.com/256/Underbot/cv/object"

	"gitlab.com/256/Underbot/ai"

	"github.com/pkg/errors"
	"gocv.io/x/gocv"
	"golang.org/x/image/bmp"
)

var (
	colors [][]uint8 // Random colors that can be used for the debugging rectangles
)

var rnd *rand.Rand
var gray = color.RGBA{193, 193, 193, 255}

var objects []object.Object // A slice of the objects detected in the game
// RecognizedObjects is a slice of the detected objects that have been recognized as something
var RecognizedObjects []object.Object

// Fills the slice of random colors up to 300 random colors. This is necessary as random colors can't be generated as quickly on the spot
func init() {
	rnd = rand.New(rand.NewSource(time.Now().Unix())) // A random instance seeded by the current unix timestamp
	for i := 0; i < 300; i++ {
		addToColor()
	}
}

func addToColor() {
	colors = append(colors, []uint8{randCol(), randCol(), randCol()})
}

func randomColor(i int) color.Color {
	return color.RGBA{colors[i][0], colors[i][1], colors[i][2], 255}
}

//randCol generates a random number between 0 and 256 as a uint8
func randCol() uint8 {
	return uint8(rnd.Int31n(255))
}

// GetObjects returns all the objects' rectangles detected
func GetObjects() []object.Object {
	return objects
}

// GetRecognizedObjects returns all the identified objects
func GetRecognizedObjects() []object.Object {
	return RecognizedObjects
}

// ProcessImage processes image, taking necessary action on the game, and then returns image with debugging information about what the CV sees
func ProcessImage(img image.RGBA) image.Image {
	// Converts incoming image into a Mat
	src := imageToMat(img)
	defer src.Close()

	// Converts the Mat into a thresholded image
	thresMat := threshold(src)
	defer thresMat.Close()

	// Find the contours (individual items on screen)
	contours := gocv.FindContours(thresMat, gocv.RetrievalTree, gocv.ChainApproxSimple)

	// Resets the global variables
	objects = []object.Object{}
	RecognizedObjects = []object.Object{}

	// Iterate through the detected objects (literal objects, not the ones in the object package yet)
	for i, obj := range contours {
		// Generates more random colors if needed
		if len(obj) > len(colors) {
			addToColor()
		}

		// Gets surrounding rectangle of object
		rec := rect.GetRectangle(obj)

		// The color the surrounding rectangle should have
		var dispColor color.Color

		// Determine coloring based on coloring parameter
		if params.Coloring == 0 {
			// Make the surrounding rectangle a random color
			dispColor = randomColor(i)
		} else if params.Coloring == 1 {
			// Make the surrounding rectangle gray (will change if the object is detected)
			dispColor = gray
		}

		// The color of the object detected
		objColor := rect.CenterColor(img.SubImage(rec))

		// Create new object instance to be build upon
		obj := object.Object{Bounds: rec, ID: i + 1, Color: objColor, Recognized: false, RecogObj: object.RecognizedObject{}}

		// Determine if the object is a RecognizableObject, and sets the proper field values
		recognize(&obj)

		// Change the coloring if the object is recognized
		if obj.Recognized {
			// If the object's recognition is black, then give it a random color instead
			if isBlack(obj.RecogObj.Color) {
				dispColor = randomColor(i)
			} else {
				dispColor = obj.Color
			}
		}

		// Draw a rectangle around the object
		rect.DrawObject(&img, dispColor, obj)

		// Add the object to the global list of objects
		objects = append(objects, obj)
	}
	go ai.Handle(objects, RecognizedObjects)
	return &img
}

// Determines if the color input is black
func isBlack(col color.Color) bool {
	return (col == color.RGBA{0, 0, 0, 255})
}

// Converts an image into a thresholded black and white Mat
func threshold(src gocv.Mat) gocv.Mat {
	// Placeholder for colorless original Mat
	srcGray := gocv.NewMat()
	defer srcGray.Close()

	// Put original image into srcGray without color
	gocv.CvtColor(src, &srcGray, gocv.ColorBGRToGray)

	// Places threshold onto the image
	thresMat := gocv.NewMat()
	gocv.Threshold(srcGray, &thresMat, 50, 255, gocv.ThresholdBinary)

	return thresMat
}

// Converts an image into a GoCV mat
func imageToMat(img image.RGBA) gocv.Mat {
	src, err := gocv.IMDecode(imageToBmp(img), 1)
	if err != nil {
		panic(errors.Wrap(err, "Failed to decode incoming image for OpenCV"))
	}
	return src
}

// Converts image to a slice of bytes representing the bitmap encoding of the image
func imageToBmp(img image.RGBA) []byte {
	// Placeholder for the bytes from the encoding to go to
	buf := new(bytes.Buffer)

	// Fill the buffer with the bitmap encoding
	err := bmp.Encode(buf, &img)
	if err != nil {
		panic(errors.Wrap(err, "Failed to encode image to bitmap format"))
	}
	return buf.Bytes()
}

// Function handling the actions that should be taken if an object is recognized
func recTreatment(obj *object.Object, recogObj object.RecognizedObject) {
	obj.Recognized = true
	obj.RecogObj = recogObj

	// Add object to the global slice of recognized objects
	RecognizedObjects = append(RecognizedObjects, *obj)
}

// Determines if an object is a RecognizableObject, and take action if so
func recognize(obj *object.Object) {
	// Gets the height and width of the object and places it into a Point object
	size := image.Point{obj.Bounds.Dx(), obj.Bounds.Dy()}

	// Iterate over the possible recognizable objects
	for _, recogObj := range object.RecognizableObjects {
		if !isBlack(obj.Color) { // If the object is colored properly, then check for size and coloring equality
			if !num.PntWithin(recogObj.Size, size, params.Leniance) {
				continue
			} else {
				if recogObj.Color == obj.Color {
					recTreatment(obj, recogObj)
				}
			}
		} else if num.PntWithin(recogObj.Size, size, params.Leniance) { // If the recognized object is black, then only check for size
			recTreatment(obj, recogObj)
		}
	}
}
