// +build linux

package impl

import (
	"os"
	"strings"
	"syscall"

	"github.com/BurntSushi/xgbutil/ewmh"
	ps "github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
)

/*
0 - Resumed/Playing
1 - Paused
*/
var gamestate = 0

// Process gets the process of a window
func (xWin window) Process() (*os.Process, error) {
	// Get the PID of the window
	pid, err := ewmh.WmPidGet(xWin.parent.conn, xWin.winID)
	if err != nil {
		pid, err = findUndertale()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get the pid of the window")
		}
	}

	// Use the PID to create a os.Process instance for the ability to pause the game, etc.
	process, err := os.FindProcess(int(pid))
	if err != nil {
		return nil, errors.Wrap(err, "failed to find the process based on the pid")
	}
	return process, nil
}

// Pause pauses the game
func (xWin window) Pause() error {
	if gamestate == 0 {
		proc, err := xWin.Process()
		if err != nil {
			return errors.Wrap(err, "failed to get process to pause")
		}

		err = proc.Signal(syscall.SIGSTOP)
		if err != nil {
			return errors.Wrap(err, "failed to send sigstop")
		}
		gamestate = 1
	}
	return nil
}

// Resume resumes the game
func (xWin window) Resume() error {
	if gamestate == 1 {
		proc, err := xWin.Process()
		if err != nil {
			return errors.Wrap(err, "failed to get process to pause")
		}

		err = proc.Signal(syscall.SIGCONT)
		if err != nil {
			return errors.Wrap(err, "failed to send sigcont")
		}
		gamestate = 0
	}
	return nil
}

// Indicators within a process name for an Undertale related process
var undertaleProcessNames = []string{"runner", "under", "tale"}

// Gets the PID of the Undertale process
func findUndertale() (uint, error) {
	for _, processName := range undertaleProcessNames {
		pid, err := executableNameToPid(processName)
		if err == nil {
			return pid, nil
		}
	}
	return 0, errors.New("failed to find undertale process")
}

// Finds the PID based on the executable name of a process
func executableNameToPid(processName string) (uint, error) {
	processes, err := ps.Processes()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get the list of processes")
	}
	for _, process := range processes {
		if strings.Contains(strings.ToLower(process.Executable()),
			strings.ToLower(processName)) && process.Executable() != "Underbot" {
			return uint(process.Pid()), nil
		}
	}
	return 0, errors.New("Could not find the process")
}
