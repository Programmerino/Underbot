package ai

import (
	"gitlab.com/256/Underbot/cv/object"
)

// CurrentState describes what is happening in the game, such as if a battle has started or something like that
var CurrentState = States[0]

// Handle is the introduction function to the AI segment. See update() for more details
func Handle(objects []object.Object, recognizedObjects []object.Object) {
	update(objects, recognizedObjects)
}

// Update runs the appropriate function depending on the current GameState
func update(objects []object.Object, recognizedObjects []object.Object) {
	CurrentState = identify(recognizedObjects)
	CurrentState.updateFun(objects)
}

// A helper type to determine the most probable GameState.
type stateProbability struct {
	state   State
	matches int
}

// Determines the GameState that matches the most with the recognized objects
func identify(recognizedObjects []object.Object) State {
	var mostProbable stateProbability

	for _, state := range States {
		// How many matches were found between the objects recognized, and the state's signs
		matchNum := len(matches(state.signs, recognizedObjects))

		// Replaces mostProbable if it is more probable
		if mostProbable.matches < matchNum {
			mostProbable.state = state
			mostProbable.matches = matchNum
		}
	}

	return mostProbable.state
}

// Finds the matching elements between recognized objects and normal objects
func matches(objects1 []object.RecognizedObject, objects2 []object.Object) []object.RecognizedObject {
	var matching []object.RecognizedObject

	// Compares every object in the first slice with every object in the second slice.
	// If two identical objects are found within the slices, then add them to the matching slice
	for _, obj1 := range objects1 {
		for _, obj2 := range objects2 {
			if obj1 == obj2.RecogObj {
				matching = append(matching, obj1)
			}
		}
	}

	return matching
}
