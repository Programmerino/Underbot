package ai

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	"github.com/pkg/errors"
	"gitlab.com/256/Underbot/ai/pathfinding"
	"gitlab.com/256/Underbot/cv/object"
	"gitlab.com/256/Underbot/cv/params"
	"gitlab.com/256/Underbot/cv/rect"
	"gitlab.com/256/Underbot/sys"
)

// GridShow determines if the pathfinding grid is enabled
var GridShow bool

// Should be ran in every update function (other than unknown)
func genUpdate() {
	unknownFrames = 0
}

// BattleMenuUpdate is the function run every frame when the battle menu is detected
func BattleMenuUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	recObjects, err := GetWanted(objects, []object.RecognizableObject{
		object.RecMap["battleOption"],
		object.RecMap["redHeart"],
	})
	if err != nil {

	}
	recMap := Map(recObjects)

	var farthest int
	var fightOption object.RecognizedObject
	for i, option := range recMap["battleOption"] {
		if i == 0 {
			farthest = option.Parent.ID
			fightOption = option
			continue
		}
		if farthest < option.Parent.ID {
			farthest = option.Parent.ID
			fightOption = option
		}
	}

	redHeart := rect.RectangleCenter(recMap["redHeart"][0].Parent.Bounds)
	fightCenter := rect.RectangleCenter(fightOption.Parent.Bounds)

	if redHeart.X > fightCenter.X {
		err := win.Press("left")
		if err != nil {
			return errors.Wrap(err, "couldn't press left arrow key")
		}
	} else {
		err := win.Press("z")
		if err != nil {
			return errors.Wrap(err, "couldn't press z key")
		}
	}

	return nil
}

// EmptyUpdate does nothing
func EmptyUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	return nil
}

// DialogueUpdate is the function run every frame when dialogue is detected
func DialogueUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	err := win.Press("x")
	if err != nil {
		return errors.Wrap(err, "failed to press x key")
	}
	err = win.Press("z")
	if err != nil {
		return errors.Wrap(err, "failed to press z key")
	}
	return nil
}

// InBattleUpdate is the function run every frame when dialogue is detected
func InBattleUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	heartObjects, err := GetWanted(objects, append([]object.RecognizableObject{}, object.Hearts...))
	if len(heartObjects) == 0 {
		// Stall until items can be found
		failedRetrieval++
		// If stalling takes too long, then try to get unstuck
		if failedRetrieval > params.FailedLimit {
			fmt.Println("The AI is unsure about what is happening. Trying to get unstuck...")
			unstuck(win)
		}
		return nil
	}
	failedRetrieval = 0

	// Draw the tiles on the screen
	tiles, err := pathfinding.MakeTiles(*img)
	if err != nil {
		return errors.Wrap(err, "failed to create the screen tiles")
	}

	// Calculate which tile Frisk is in

	//currentTile := pathfinding.GetCurrentTile(image.Point{averageX, averageY})
	//_, err = pathfinding.GetPath(currentTile, pathfinding.GetGoal())
	if err != nil {
		return errors.Wrap(err, "failed to generate path")
	}
	/*
		for _, pathTile := range path {

		}
	*/
	if GridShow {
		for _, tile := range tiles {
			rect.DrawRectangle(img, tile.Color, tile.Rectangle)
		}
	}
	return nil
}

var failedRetrieval int

// SaveUpdate is the function that presses the save button
func SaveUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	recObjects, err := GetWanted(objects, []object.RecognizableObject{
		object.RecMap["redHeart"],
		object.RecMap["saveBox"],
	})
	if err != nil {
		// Stall until items can be found
		failedRetrieval++
		// If stalling takes too long, then try to get unstuck
		if failedRetrieval > params.FailedLimit {
			fmt.Println("The AI is unsure about what is happening. Trying to get unstuck...")
			unstuck(win)
		}
		return nil
	}
	failedRetrieval = 0
	recMap := Map(recObjects)

	// Determines if the heart is on the left side, indicating that it is selecting save
	if rect.RectangleCenter(recMap["redHeart"][0].Parent.Bounds).X < rect.RectangleCenter(recMap["saveBox"][0].Parent.Bounds).X {
		err := win.Press("z")
		if err != nil {
			return errors.Wrap(err, "failed to press z key")
		}
		err = win.Press("z")
		if err != nil {
			return errors.Wrap(err, "failed to press z key")
		}
	} else {
		err := win.Press("left")
		if err != nil {
			return errors.Wrap(err, "failed to press left key")
		}
		err = win.Press("z")
		if err != nil {
			return errors.Wrap(err, "failed to press z key")
		}
		err = win.Press("z")
		if err != nil {
			return errors.Wrap(err, "failed to press z key")
		}
	}
	return nil
}

// OutsideBattleUpdate is the function that pathfinds through the game
func OutsideBattleUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	friskObjects, _ := GetWanted(objects, object.Frisk)
	if len(friskObjects) == 0 {
		// Stall until items can be found
		failedRetrieval++
		// If stalling takes too long, then try to get unstuck
		if failedRetrieval > params.FailedLimit {
			fmt.Println("The AI is unsure about what is happening. Trying to get unstuck...")
			unstuck(win)
		}
		return nil
	}
	failedRetrieval = 0

	// The point in the middle of all the frisk objects
	var averageX int
	var averageY int
	for _, obj := range friskObjects {
		point := rect.RectangleCenter(obj.Parent.Bounds)
		averageX += point.X
		averageY += point.Y
	}
	averageX = averageX / len(friskObjects)
	averageY = averageY / len(friskObjects)

	img.Set(averageX, averageY, color.RGBA{255, 0, 0, 255})
	rect.VLine(img, color.RGBA{255, 0, 0, 255}, averageX, averageY-10, averageY+10)
	rect.HLine(img, color.RGBA{255, 0, 0, 255}, averageX-10, averageY, averageX+10)

	// Draw the tiles on the screen
	tiles, err := pathfinding.MakeTiles(*img)
	if err != nil {
		return errors.Wrap(err, "failed to create the screen tiles")
	}

	// Calculate which tile Frisk is in

	currentTile := pathfinding.GetCurrentTile(image.Point{averageX, averageY})
	_, err = pathfinding.GetPath(currentTile, pathfinding.GetGoal())
	if err != nil {
		return errors.Wrap(err, "failed to generate path")
	}
	/*
		for _, pathTile := range path {

		}
	*/
	if GridShow {
		for _, tile := range tiles {
			rect.DrawRectangle(img, tile.Color, tile.Rectangle)
		}
	}
	return nil
}

// unstuck will do keypresses that attempt to get the game to a state that can be detected properly
func unstuck(win sys.Window) error {
	keys := []string{"z", "x", "up", "left", "right", "down", "enter"}
	randKey := keys[rand.Intn(len(keys))]
	err := win.Press(randKey)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to press %s key", randKey))
	}
	return nil
}

// Counts how many frames the state has been unknown
var unknownFrames int

// UnknownUpdate is the update function that just calls unstuck
func UnknownUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	// Stall until items can be found
	unknownFrames++
	// If stalling takes too long, then try to get unstuck
	if unknownFrames > params.FailedLimit {
		fmt.Println("The AI is unsure about what is happening. Trying to get unstuck...")
		err := unstuck(win)
		if err != nil {
			return errors.Wrap(err, "failed to try to get unstuck")
		}
	}
	return nil
}

// GetWanted returns a slice of objects based on what the recognized objects wanted are
func GetWanted(objects []object.Object, wanted []object.RecognizableObject) ([]object.RecognizedObject, error) {
	var found []object.RecognizedObject

	// Used to make sure that the types found are the same as the ones wanted
	var wantedCompare []object.RecognizableObject

	// Used to ensure that duplicate items of the same type aren't added to wantedCompare
	var typesUsed []object.RecognizableObject
	for _, wantObj := range wanted {
		for _, obj := range objects {
			if obj.RecogObj.Type == wantObj {
				found = append(found, obj.RecogObj)
				if !contains(typesUsed, obj.RecogObj.Type) {
					wantedCompare = append(wantedCompare, obj.RecogObj.Type)
					typesUsed = append(typesUsed, obj.RecogObj.Type)
				}
			}
		}
	}
	if !equal(wanted, wantedCompare) {
		return found, errors.New("could not get all the needed items")
	}
	return found, nil
}

// Map turns a slice of RecognizableObject into a map that can be searched with the string name
func Map(recObjects []object.RecognizedObject) (recMap map[string][]object.RecognizedObject) {
	recMap = make(map[string][]object.RecognizedObject)
	for _, obj := range recObjects {
		if recMap[obj.Type.Name] == nil {
			recMap[obj.Type.Name] = []object.RecognizedObject{obj}
		} else {
			recMap[obj.Type.Name] = append(recMap[obj.Type.Name], obj)
		}
	}
	return recMap
}

// equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func equal(a, b []object.RecognizableObject) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func contains(s []object.RecognizableObject, e object.RecognizableObject) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
