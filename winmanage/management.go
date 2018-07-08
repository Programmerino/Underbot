package winmanage

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/go-vgo/robotgo"
	"github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
)

// The height and width that the window should be resized to (these values usually keeps the processing at 60fps)
const (
	height = 480
	width  = 640
)

// The title of the UndertaleWindow
var title string

// Get an instance of UndertaleWindow based on the window clicked
func Get(x *xgbutil.XUtil, name string) *UndertaleWindow {
	title = name
	winID := getWinID(x)
	return getWinInfo(x, winID)
}

// Gets UndertaleWindow from window ID
func getWinInfo(x *xgbutil.XUtil, winID xproto.Window) *UndertaleWindow {
	//Get the title of the window for debugging
	name, err := ewmh.WmNameGet(x, winID)
	if err != nil {
		name = "unknown window"
		fmt.Println("Could not get the name of the window")
	}
	fmt.Printf("You selected %s\n", name)

	// Create xwindow instance from xproto window id
	xwinWin := xwindow.New(x, winID)

	// Attempt to resize the window for WM compatible DEs
	err = xwinWin.WMResize(width, height)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Could not resize window (method 1)"))
	}

	// Sleep for DE to recognize resize request before continuing
	time.Sleep(time.Millisecond * 30)

	// Get a rectangle describing the dimensions of the window for proper rendering in the case the resize fails silently
	rect, err := xwinWin.Geometry()
	if err != nil {
		panic(errors.Wrap(err, "Error getting window size"))
	}
	// Get the PID of the window
	pid, err := ewmh.WmPidGet(x, winID)
	if err != nil {
		fmt.Println("Failed to get PID")
		pid = findUndertale()
		fmt.Println("Got PID from process list")
	}

	// Use the PID to create a os.Process instance for the ability to pause the game, etc.
	process, err := os.FindProcess(int(pid))
	if err != nil {
		panic(errors.Wrap(err, "Could not find the process from X's PID"))
	}
	return newUndertaleWindow(rect.Width(), rect.Height(), winID, process)
}

// Get the ID of the active window
func getWinID(x *xgbutil.XUtil) xproto.Window {
	fmt.Println("Click the window to act upon")
	waitForMouseClick()
	time.Sleep(time.Second / 2)
	// Get the xproto window ID from the active window (the one last clicked usually)
	winID, err := ewmh.ActiveWindowGet(x)
	if err != nil {
		panic(errors.Wrap(err, "Error getting active window"))
	}
	return winID
}

// Stalls until the left mouse button is pressed
func waitForMouseClick() {
	mleft := robotgo.AddEvent("mleft")
	if mleft == 0 {
		return
	}
	panic("mleft was not 0")
}

// Indicators within a process name for an Undertale related process
var undertaleProcessNames = []string{"runner", "under", "tale"}

// Gets the PID of the Undertale process
func findUndertale() uint {
	for _, processName := range undertaleProcessNames {
		pid, err := executableNameToPid(processName)
		if err == nil {
			return pid
		}
	}
	panic("Could not find Undertale instance")
}

// Finds the PID based on the executable name of a process
func executableNameToPid(processName string) (uint, error) {
	processes, err := ps.Processes()
	if err != nil {
		panic(err)
	}
	for _, process := range processes {
		if strings.Contains(strings.ToLower(process.Executable()), strings.ToLower(processName)) && process.Executable() != title {
			return uint(process.Pid()), nil
		}
	}
	return 0, errors.New("Could not find the process")
}
