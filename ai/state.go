package ai

import (
	"image"

	"github.com/pkg/errors"

	"gitlab.com/256/Underbot/cv/object"
	"gitlab.com/256/Underbot/sys"
)

// State determines the factors that indicate what is going on in the game.
// An example of a state is when the game is at the battle menu
type State struct {
	Name      string
	signs     []object.RecognizableObject
	antiSigns []object.RecognizableObject // Objects that will never appear in this state
	// An update function called every frame for the specific state to handle
	updateFun func([]object.Object, sys.Window, *image.RGBA) error
	times     int // Specifies how many times per 10 frames the function should run. If -1, then run all the time. Limited 5
}

// NewState creates new State instance with parameter checking
func NewState(name string, signs []object.RecognizableObject, antiSigns []object.RecognizableObject,
	updateFun func([]object.Object, sys.Window, *image.RGBA) error, times int) State {

	tempState := State{Name: name, signs: signs, antiSigns: antiSigns, updateFun: updateFun, times: times}
	if err := tempState.check(); err != nil {
		panic(errors.Wrap(err, "create state was invalid"))
	}
	return tempState

}

// Checks a State instance for validity
func (state *State) check() error {
	if state.Name == "" {
		return errors.New("state has empty name")
	}
	if state.updateFun == nil {
		return errors.New("state has invalid update function")
	}
	return nil
}

// Convience object for when no signs/antisigns are needed for the state
var emptyObjects = []object.RecognizableObject{}

// List of signs for each state

var battleMenuSigns = append([]object.RecognizableObject{
	object.RecMap["narratorBox"], // narratorBox
}, object.Hearts...)

var inBattleSigns = append([]object.RecognizableObject{
	object.RecMap["fightBox"], // fightBox
}, object.Hearts...)

var dialogueSigns = object.Dialogue

var dialogueAntiSigns = []object.RecognizableObject{
	object.RecMap["attackGoal"], // attackGoal
	object.RecMap["attackPeg"],  // attackPeg
}

var outsideBattleSigns = append([]object.RecognizableObject{}, object.Frisk...)

var saveScreenSigns = []object.RecognizableObject{
	object.RecMap["saveBox"],  // saveBox
	object.RecMap["redHeart"], // redHeart
}
var attackGoalSigns = []object.RecognizableObject{
	object.RecMap["attackGoal"], // attackGoal
	object.RecMap["attackPeg"],  // attackPeg
}

// States holds the list of possible game states
var States = []State{
	// When no game state seems suitable for the recognized objects (or lack thereof)
	NewState("Unknown", emptyObjects, emptyObjects, UnknownUpdate, -1),
	// After encountering a battle when no option has been pressed yet
	NewState("battleMenu", battleMenuSigns, emptyObjects, BattleMenuUpdate, -1),
	NewState("inBattle", inBattleSigns, emptyObjects, InBattleUpdate, -1),                // When the opponent is attacking
	NewState("dialogue", dialogueSigns, dialogueAntiSigns, DialogueUpdate, 5),            // When outside battle with dialogue
	NewState("outsideBattle", outsideBattleSigns, emptyObjects, OutsideBattleUpdate, -1), // When outside battle
	// The screen where you can choose to "Save" or "Return" at a checkpoint
	NewState("saveScreen", saveScreenSigns, emptyObjects, SaveUpdate, -1),
	// After pressing fight, when z needs to be pressed with good timing
	NewState("attackGoal", attackGoalSigns, emptyObjects, BattleMenuUpdate, -1),
}
