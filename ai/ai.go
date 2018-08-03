package ai

import (
	"image"

	"github.com/pkg/errors"
	"gitlab.com/256/Underbot/cv/object"
	"gitlab.com/256/Underbot/sys"
)

// CurrentState describes what is happening in the game, such as if a battle has started or something like that
var CurrentState = States[0]

// Disabled determines whether or not the AI should be turned on
var Disabled bool

// Handle is the introduction function to the AI segment. See update() for more details
func Handle(objects []object.Object, recognizedObjects []object.Object, win sys.Window, img *image.RGBA) error {
	if !Disabled {
		err := update(objects, recognizedObjects, win, img)
		if err != nil {
			return errors.Wrap(err, "failed to run the update function")
		}
	}
	return nil
}

// How many frames have happened (resets at 10)
var frames int

// How many times the update function has been called (resets when frames reaches 10)
var usedFrames int

// Update runs the appropriate function depending on the current GameState
func update(objects []object.Object, recognizedObjects []object.Object, win sys.Window, img *image.RGBA) error {
	CurrentState = identify(recognizedObjects)
	if CurrentState.times != -1 {
		if (usedFrames < CurrentState.times) && (frames%2 == 0) {
			err := CurrentState.updateFun(objects, win, img)
			if err != nil {
				return errors.Wrap(err, "the update function for the state failed")
			}
			usedFrames++
		}
	} else {
		err := CurrentState.updateFun(objects, win, img)
		if err != nil {
			return errors.Wrap(err, "the update function for the state failed")
		}
	}
	if frames >= 10 {
		frames = 0
		usedFrames = 0
	}
	frames++
	return nil
}

// A helper type to determine the most probable GameState.
type stateProbability struct {
	state   State
	matches int
}

// Determines the GameState that matches the most with the recognized objects
func identify(recognizedObjects []object.Object) State {
	// Initialize it with unknown state to prevent nil pointer when running the game state update function
	var mostProbable = stateProbability{States[0], 0}

	for _, state := range States {
		// How many matches were found between the objects recognized, and the state's signs
		matchNum := len(matches(state.signs, recognizedObjects))
		antiMatchNum := len(matches(state.antiSigns, recognizedObjects))

		if antiMatchNum > 0 {
			continue
		}

		// Replaces mostProbable if it is more probable
		if mostProbable.matches < matchNum {
			mostProbable.state = state
			mostProbable.matches = matchNum
		}
	}

	return mostProbable.state
}

// Finds the matching elements between recognizable objects and normal objects
func matches(objects1 []object.RecognizableObject, objects2 []object.Object) []object.Object {
	var matching []object.Object

	// Compares every object in the first slice with every object in the second slice.
	// If two identical objects are found within the slices, then add them to the matching slice
	for _, obj1 := range objects1 {
		for _, obj2 := range objects2 {
			if obj1 == obj2.RecogObj.Type {
				if !alreadyInSlice(obj2, matching) {
					matching = append(matching, obj2)
				}
			}
		}
	}

	return matching
}

func alreadyInSlice(a object.Object, list []object.Object) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
