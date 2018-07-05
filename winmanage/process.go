package winmanage

import (
	"os"

	"gitlab.com/256/Underbot/winmanage/gamestate"
)

// A struct holding the process to an UndertaleWindow, as well as the GameState for the game
type process struct {
	Process *os.Process
	state   gamestate.GameState
}

// Creates new process
func newProcess(Process *os.Process) process {
	tempProcess := process{Process, gamestate.Playing}
	tempProcess.check()
	return tempProcess
}

// Check the process struct for validity
func (Process *process) check() {
	if Process == nil {
		panic("given nil process during process check")
	}
	if Process.Process == nil {
		panic("The underlying process is a nil pointer")
	}
	if (Process.Process == &os.Process{}) {
		panic("The underlying process is an empty struct")
	}
	if Process.Process.Pid == 0 {
		panic("Underlying has PID of 0")
	}
	Process.state.Check()
}
