package cv

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"gitlab.com/256/Underbot/cv/num"
	"gitlab.com/256/Underbot/cv/params"
	"gitlab.com/256/Underbot/cv/rect"
	"gitlab.com/256/Underbot/sys"

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

// Fills the slice of random colors up to 300 random colors.
// This is necessary as random colors can't be generated as quickly on the spot
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
	// In case there aren't enough colors, just make a not very random color instead on the spot quickly
	if len(colors) < i {
		index := len(colors) - 1
		return color.RGBA{colors[index][0], colors[index/2][1], colors[index/3][2], 255}
	}
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

// ProcessImage processes image, runs AI code,
// and then modifies image with debugging information about what the CV sees
func ProcessImage(img *image.RGBA, win sys.Window) error {
	// Converts incoming image into a Mat
	src, err := imageToMat(*img)
	if err != nil {
		return errors.Wrap(err, "failed to convert the image to a Mat")
	}
	defer func() {
		err := src.Close()
		if err != nil {
			panic(errors.Wrap(err, "failed to close the orginal Mat"))
		}
	}()

	// Converts the Mat into a thresholded image
	thresMat := threshold(src)
	defer func() {
		err := thresMat.Close()
		if err != nil {
			panic(errors.Wrap(err, "failed to close the threshold Mat"))
		}
	}()

	// Find the contours (individual items on screen)
	contours := gocv.FindContours(thresMat, gocv.RetrievalTree, gocv.ChainApproxSimple)

	// Resets the global variables
	objects = []object.Object{}
	RecognizedObjects = []object.Object{}

	// A secondary iterator that only iterates each time a random color is used.
	// This is to prevent unneeded extra colors from being created
	usedColors := 0

	// Iterate through the detected objects (literal objects, not the ones in the object package yet)
	for i, obj := range contours {
		// Generates more random colors if needed
		if (len(obj) > len(colors)) && params.Coloring == 1 {
			addToColor()
		}

		// Gets surrounding rectangle of object
		rec := rect.GetRectangle(obj)

		// The color the surrounding rectangle should have
		var dispColor color.Color

		// Determine coloring based on coloring parameter
		if params.Coloring == 0 {
			// Make the surrounding rectangle a random color
			dispColor = randomColor(usedColors)
			usedColors++
		} else if params.Coloring == 1 {
			// Make the surrounding rectangle gray (will change if the object is detected)
			dispColor = gray
		}

		// The color of the object detected
		objColor, err := rect.CenterColor(img.SubImage(rec))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not detect the center color of object %v", i))
		}
		// Create new object instance to be build upon
		obj := object.Object{Bounds: rec, ID: i + 1, Color: objColor, Recognized: false, RecogObj: object.RecognizedObject{}}

		// Determine if the object is a RecognizableObject, and sets the proper field values
		err = recognize(&obj)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed in recognizing object %v", i))
		}

		// Change the coloring if the object is recognized
		if obj.Recognized {
			// If the object's recognition is black, then give it a random color instead
			if isBlack(obj.RecogObj.Type.Color) {
				dispColor = randomColor(usedColors)
				usedColors++
			} else {
				dispColor = obj.Color
			}
		}

		// Draw a rectangle around the object
		rect.DrawObject(img, dispColor, obj)

		// Add the object to the global list of objects
		objects = append(objects, obj)
	}
	err = ai.Handle(objects, RecognizedObjects, win, img)
	if err != nil {
		return errors.Wrap(err, "ai failed to act upon the objects")
	}
	return nil
}

// Determines if the color input is black
func isBlack(col color.Color) bool {
	return (col == color.RGBA{0, 0, 0, 255})
}

// Converts an image into a thresholded black and white Mat
func threshold(src gocv.Mat) gocv.Mat {
	// Placeholder for colorless original Mat
	srcGray := gocv.NewMat()
	defer func() {
		err := srcGray.Close()
		if err != nil {
			panic(errors.Wrap(err, "failed to close the gray Mat"))
		}
	}()

	// Put original image into srcGray without color
	gocv.CvtColor(src, &srcGray, gocv.ColorBGRToGray)

	// Places threshold onto the image
	thresMat := gocv.NewMat()
	gocv.Threshold(srcGray, &thresMat, 50, 255, gocv.ThresholdBinary)

	return thresMat
}

// Converts an image into a GoCV mat
func imageToMat(img image.RGBA) (gocv.Mat, error) {
	bmp, err := imageToBmp(img)
	if err != nil {
		return gocv.Mat{}, errors.Wrap(err, "failed to convert the image to a BMP")
	}
	src, err := gocv.IMDecode(bmp, 1)
	if err != nil {
		return gocv.Mat{}, errors.Wrap(err, "Failed to decode incoming image for OpenCV")
	}
	return src, nil
}

// Converts image to a slice of bytes representing the bitmap encoding of the image
func imageToBmp(img image.RGBA) ([]byte, error) {
	// Placeholder for the bytes from the encoding to go to
	buf := new(bytes.Buffer)

	// Fill the buffer with the bitmap encoding
	err := bmp.Encode(buf, &img)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to encode image to bitmap format")
	}
	return buf.Bytes(), nil
}

// Function handling the actions that should be taken if an object is recognized
func recTreatment(obj *object.Object, recogObj object.RecognizedObject) {
	obj.Recognized = true
	obj.RecogObj = recogObj

	// Add object to the global slice of recognized objects
	RecognizedObjects = append(RecognizedObjects, *obj)
}

// Determines if an object is a RecognizableObject, and take action if so
func recognize(obj *object.Object) error {
	err := obj.Check()
	if err != nil {
		return errors.Wrap(err, "refusing to operate on invalid object")
	}
	// Gets the height and width of the object and places it into a Point object
	size := image.Point{obj.Bounds.Dx(), obj.Bounds.Dy()}

	// Iterate over the possible recognizable objects
	for _, recogObj := range object.RecognizableObjects {
		var leniance int
		if recogObj.Leniance == -1 {
			leniance = params.Leniance
		} else {
			leniance = recogObj.Leniance
		}
		if !isBlack(recogObj.Color) { // If the object is colored properly, then check for size and coloring equality
			if num.PntWithin(recogObj.Size, size, leniance) && recogObj.Color == obj.Color {
				recognized, err := object.NewRecognizedObject(obj, recogObj)
				if err != nil {
					return errors.Wrap(err, "failed to create new recognized object")
				}
				recTreatment(obj, recognized)
			}
			// If the recognized object is black, then only check for size
		} else if num.PntWithin(recogObj.Size, size, leniance) {
			recognized, err := object.NewRecognizedObject(obj, recogObj)
			if err != nil {
				return errors.Wrap(err, "failed to create new recognized object")
			}
			recTreatment(obj, recognized)
		}
	}
	return nil
}
