package object

import (
	"image"
	"image/color"
)

// Object holds information for a detected object found in the window
type Object struct {
	Bounds     image.Rectangle // Rectangle describing where it is on screen
	ID         int             // Used for debugging to provide a method of describing an object
	Color      color.Color     // The color of the pixel in the center of the object
	Recognized bool
	RecogObj   RecognizedObject // Holds information for the object it is recognized as if it is recognized
}

// NewObject creates new instance of an Object with parameter checking
func NewObject(rect image.Rectangle, id int, color color.Color, recogobj RecognizedObject) Object {
	obj := Object{}
	obj.Bounds = rect
	obj.ID = id
	obj.Color = color
	obj.Recognized = (recogobj != RecognizedObject{})
	if obj.Recognized {
		obj.RecogObj = recogobj
	}
	return obj
}

// Determines the validity of an object
func (obj *Object) check() {
	if (obj == &Object{}) {
		panic("Empty object")
	}
	if (obj.Bounds == image.Rectangle{}) {
		panic("Bounds are not specified")
	}
	if obj.ID < 1 {
		panic("Invalid ID")
	}
	if obj.Color == nil {
		panic("Invalid color")
	}
	obj.RecogObj.check()
}

// RecognizedObject determines the specifications for an object to be recognized by
type RecognizedObject struct {
	Name  string
	Size  image.Point
	Color color.Color
}

// Create new RecognizedObject with parameter checking
func newSpecs(name string, width, height int, color color.Color) RecognizedObject {
	return RecognizedObject{name, image.Point{width, height}, color}
}

// Checks a RecognizedObject for validity
func (recogObj *RecognizedObject) check() {
	if (recogObj == &RecognizedObject{}) {
		panic("Empty RecognizedObject")
	}
	if recogObj.Name == "" {
		panic("Empty name")
	}
	if (recogObj.Size == image.Point{}) {
		panic("Size is empty")
	}
	if recogObj.Color == nil {
		panic("Invalid color")
	}
}

// Hearts is a list holding all the hearts possible (eg. red, blue, green)
var Hearts = []RecognizedObject{
	RecognizableObjects[1],
	RecognizableObjects[2],
	RecognizableObjects[3],
}

// Frisk is a list holding all the objects that would be seen on Frisk
var Frisk = []RecognizedObject{
	RecognizableObjects[6],
	RecognizableObjects[7],
	RecognizableObjects[8],
	RecognizableObjects[9],
}

// RecognizableObjects holds all the possible recognizable objects in the game
var RecognizableObjects = []RecognizedObject{
	newSpecs("narratorBox", 574, 139, color.RGBA{0, 0, 0, 255}),       // 0: The largest rectangle in battleMenu that usually holds narration, item options, etc.
	newSpecs("redHeart", 15, 15, color.RGBA{255, 0, 0, 255}),          // 1: Traditional heart. No gravity and moves around the fightBox
	newSpecs("greenHeart", 15, 15, color.RGBA{0, 192, 0, 255}),        // 2: Green heart used in canon Undertale fights. Green heart indicates game mode where arrow keys are used to shield against arrows
	newSpecs("blueHeart", 15, 15, color.RGBA{0, 60, 255, 255}),        // 3: Blue heart used in canon Sans and Papyrus fights. Blue heart indicates game mode where gravity is turned on
	newSpecs("attackGoal", 18, 83, color.RGBA{0, 0, 0, 255}),          // 4: The middle of the dialogueBox after pressing "FIGHT", indicating the best place to press Z
	newSpecs("gameOverM", 127, 79, color.RGBA{254, 254, 254, 255}),    // 5: The M used in the Game Over screen, indicating that the bot has lost
	newSpecs("friskFrontFace", 27, 21, color.RGBA{255, 201, 14, 255}), // 6: The CV algorithm detects Frisk's body as two seperate body parts. This is the face
	newSpecs("friskSideFace", 19, 21, color.RGBA{255, 201, 14, 255}),  // 7: 6 but on the side
	newSpecs("friskBody", 23, 17, color.RGBA{230, 7, 248, 255}),       // 8: The CV algorithm detects Frisk's body as two seperate body parts. This is the body. This is the same for front and back
	newSpecs("friskSideBody", 13, 17, color.RGBA{61, 18, 14, 255}),    // 9: Same as 8, but on the side
}
