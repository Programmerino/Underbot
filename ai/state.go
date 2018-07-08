package ai

import (
	"gitlab.com/256/Underbot/cv/object"
)

// State determines the factors that indicate what is going on in the game. An example of a state is when the game is at the battle menu
type State struct {
	Name      string
	signs     []object.RecognizedObject
	antiSigns []object.RecognizedObject // Objects that will never appear in this state
	updateFun func([]object.Object)     // An update function called every frame for the specific state to handle
}

// NewState creates new State instance with parameter checking
func NewState(name string, signs []object.RecognizedObject, antiSigns []object.RecognizedObject, updateFun func([]object.Object)) State {
	tempState := State{name, signs, antiSigns, updateFun}
	tempState.check()
	return tempState
}

// Checks a State instance for validity
func (state *State) check() {
	if state.Name == "" {
		panic("Name is empty!")
	}
	if state.updateFun == nil {
		panic("Bad function")
	}
}

// Convience object for when no signs/antisigns are needed for the state
var emptyObjects = []object.RecognizedObject{}

// List of signs for each state

var battleMenuSigns = append([]object.RecognizedObject{
	object.RecognizableObjects[0],
}, object.Hearts...)

// States holds the list of possible game states
var States = []State{
	NewState("Unknown", []object.RecognizedObject{}, []object.RecognizedObject{}, BattleMenuUpdate), // When no game state seems suitable for the recognized objects (or lack thereof)
	NewState("battleMenu", battleMenuSigns, emptyObjects, BattleMenuUpdate),                         // After encountering a battle when no option has been pressed yet
}
