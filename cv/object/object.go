package object

import (
	"image"
	"image/color"

	"github.com/pkg/errors"
)

// Object holds information for a detected object found in the window
type Object struct {
	Bounds     image.Rectangle // Rectangle describing the object's dimensions
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

// Check determines the validity of an object
func (obj *Object) Check() error {
	if obj == nil {
		return errors.New("object is a nil pointer")
	}
	if (obj == &Object{}) {
		return errors.New("object is an empty struct")
	}
	if (obj.Bounds == image.Rectangle{}) {
		return errors.New("bounds for the object are not specified")
	}
	if obj.ID < 1 {
		return errors.New("object has invalid ID")
	}
	if obj.Color == nil {
		return errors.New("object has invalid color")
	}
	err := obj.RecogObj.check()
	if err != nil {
		return errors.New("object's child recognizedobject is invalid")
	}
	return nil
}

// RecognizedObject represents an object which has been recognized
type RecognizedObject struct {
	Parent *Object
	Type   RecognizableObject
}

// NewRecognizedObject creates a new instance of RecognizedObject safely
func NewRecognizedObject(parent *Object, Type RecognizableObject) (RecognizedObject, error) {
	temp := RecognizedObject{Parent: parent, Type: Type}
	err := temp.check()
	if err != nil {
		return RecognizedObject{}, errors.Wrap(err, "the created recognizedobject was invalid")
	}
	return temp, nil
}

// Checks a RecognizedObject for validity
func (recogObj *RecognizedObject) check() error {
	// Causes infinite loop
	/*
		err := recogObj.Parent.Check()
		if err != nil {
			return errors.Wrap(err, "the parent object was invalid")
		}
	*/
	if (recogObj == &RecognizedObject{}) {
		return errors.New("the recognizedobject is empty")
	}
	return nil
}

// RecognizableObject determines the specifications for an object to be recognized by
type RecognizableObject struct {
	Name     string
	Size     image.Point
	Color    color.Color
	Leniance int // How far off the object can be in terms of size. If set to -1, will be default set in params package
}

// Create new RecognizedObject with parameter checking
func newSpecs(name string, width, height int, color color.Color, leniance int) RecognizableObject {
	recogObj := RecognizableObject{Name: name, Size: image.Point{width, height}, Color: color, Leniance: leniance}
	err := recogObj.check()
	if err != nil {
		panic(errors.Wrap(err, "the created recognizableobject is invalid"))
	}
	return recogObj
}

// Checks a RecognizedObject for validity
func (recogObj *RecognizableObject) check() error {
	if (recogObj == &RecognizableObject{}) {
		return errors.New("recognizableobject is empty")
	}
	if recogObj.Name == "" {
		return errors.New("recognizableobject has no name")
	}
	if (recogObj.Size == image.Point{}) {
		return errors.New("recognizableobject has an empty size")
	}
	if recogObj.Color == nil {
		return errors.New("recognizableobject has an invalid size")
	}
	return nil
}

// Hearts is a list holding all the hearts possible (eg. red, blue, green)
var Hearts = []RecognizableObject{
	RecMap["redHeart"],
	RecMap["greenHeart"],
	RecMap["blueHeart"],
}

// Frisk is a list holding all the objects that would be seen on Frisk
var Frisk = []RecognizableObject{
	RecMap["friskFrontFace"],
	RecMap["friskSideFace"],
	RecMap["friskBody"],
	RecMap["friskSideBody"],
	RecMap["friskUpperBody"],
	RecMap["friskBack"],
}

// Dialogue is a list of all objects that would indicate dialogue is ocurring
var Dialogue = []RecognizableObject{
	RecMap["dialogueBox"],
}

// RecognizableObjects holds all the possible recognizable objects in the game
var RecognizableObjects = []RecognizableObject{
	// 0: The largest rectangle in battleMenu that usually holds narration, item options, etc.
	newSpecs("narratorBox", 574, 139, color.RGBA{0, 0, 0, 255}, -1),
	// 1: Traditional heart. No gravity and moves around the fightBox
	newSpecs("redHeart", 15, 15, color.RGBA{255, 0, 0, 255}, -1),
	// 2: Green heart used in canon Undertale fights.
	// Green heart indicates game mode where arrow keys are used to shield against arrows
	newSpecs("greenHeart", 15, 15, color.RGBA{0, 192, 0, 255}, -1),
	// 3: Blue heart used in canon Sans and Papyrus fights. Blue heart indicates game mode where gravity is turned on
	newSpecs("blueHeart", 15, 15, color.RGBA{0, 60, 255, 255}, -1),
	// 4: The middle of the dialogueBox after pressing "FIGHT", indicating the best place to press Z
	newSpecs("attackGoal", 18, 83, color.RGBA{0, 0, 0, 255}, -1),
	// 5: The M used in the Game Over screen, indicating that the bot has lost
	newSpecs("gameOverM", 127, 79, color.RGBA{254, 254, 254, 255}, -1),
	// 6: The CV algorithm detects Frisk's body as two separate body parts. This is the face
	newSpecs("friskFrontFace", 27, 21, color.RGBA{255, 201, 14, 255}, -1),
	// 7: 6 but on the side
	newSpecs("friskSideFace", 19, 21, color.RGBA{255, 201, 14, 255}, -1),
	// 8: The CV algorithm detects Frisk's body as two separate body parts.
	// This is the body. This is the same for front and back
	newSpecs("friskBody", 23, 17, color.RGBA{230, 7, 248, 255}, -1),
	// 9: Same as 8, but on the side
	newSpecs("friskSideBody", 13, 17, color.RGBA{61, 18, 14, 255}, -1),
	// 10: The box used for dialogue outside of battles
	newSpecs("dialogueBox", 577, 151, color.RGBA{0, 0, 0, 255}, 20),
	// 11: The box surrounding the heart during a battle
	newSpecs("fightBox", 164, 139, color.RGBA{0, 0, 0, 255}, -1),
	// 12: The box used when choosing "Save" or "Return" after getting to a checkpoint
	newSpecs("saveBox", 413, 163, color.RGBA{0, 0, 0, 255}, -1),
	// 13: Frisk is sometimes detected including both the head and the upper half of the body, so this catches that
	newSpecs("friskUpperBody", 35, 49, color.RGBA{255, 201, 14, 255}, -1),
	// 14: Frisk's back
	newSpecs("friskBack", 39, 59, color.RGBA{61, 18, 14, 255}, -1),
	// 15: The long rectangle that you are supposed to center in the attackGoal
	newSpecs("attackPeg", 7, 123, color.RGBA{255, 255, 255, 255}, -1),
	// 16: The switch with a yellow outline shown towards the beginning of the game
	newSpecs("yellowSwitch", 7, 23, color.RGBA{0, 0, 0, 255}, -1),
	// 17: The dummy towards the beginning of the game
	newSpecs("dummy", 27, 19, color.RGBA{239, 228, 176, 255}, -1),
	// 18: Toriel - Front
	newSpecs("torielFront", 15, 7, color.RGBA{255, 255, 255, 255}, 0),
	// 19: Toriel - Side
	newSpecs("torielSide", 7, 7, color.RGBA{86, 86, 211, 255}, -1),
	// 20: The entrance frame to go between rooms
	newSpecs("entrance", 65, 105, color.RGBA{255, 255, 255, 255}, -1),
	// 21: Battle Options
	newSpecs("battleOption", 107, 39, color.RGBA{0, 0, 0, 255}, -1),
}

// RecMap is map of the above slice indexed by the name attribute
var RecMap = Map(RecognizableObjects)

// Map turns a slice of RecognizableObject into a map that can be searched with the string name
func Map(recObjects []RecognizableObject) (recMap map[string]RecognizableObject) {
	recMap = make(map[string]RecognizableObject)
	for _, obj := range recObjects {
		recMap[obj.Name] = obj
	}
	return recMap
}
