package winmanage

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/pkg/errors"
	"gitlab.com/256/Underbot/sys"
)

// The height and width that the window should be resized to (these values usually keeps the processing at 60fps)
const (
	height = 480
	width  = 640
)

// The title of the Window
var title string

// The UndertaleWindow instance that will be worked upon
var mainWindow sys.Window

// GetMainWindow adds capability to get the main window from other packages
func GetMainWindow() *sys.Window {
	return &mainWindow
}

// Get an instance of Window based on the window clicked
func Get(serv sys.Server, name string) (sys.Window, error) {
	title = name
	win, err := getWinID(serv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the window")
	}
	err = modifyWindow(win)
	if err != nil {
		return nil, errors.Wrap(err, "failed to modify the window")
	}
	return win, nil
}

// Gets Window from window ID
func modifyWindow(win sys.Window) error {
	// Print the name of the window for debugging
	name, err := win.Name()
	if err != nil {
		return errors.Wrap(err, "failed to get the name")
	}
	fmt.Printf("You selected %s\n", name)

	// Resize the window to the wanted specifications
	err = win.Resize(width, height)
	if err != nil {
		return errors.Wrap(err, "failed to resize the window")
	}
	return nil
}

// Get window based on where user clicked
func getWinID(serv sys.Server) (sys.Window, error) {
	fmt.Println("Shift-click the window to act upon")
	err := waitForSelect()
	if err != nil {
		return nil, errors.Wrap(err, "could not wait for the mouse to be clicked")
	}
	time.Sleep(time.Second / 2)
	return serv.ActiveWindow()
}

// Stalls until the left mouse button is pressed
func waitForSelect() error {
	ctrl := robotgo.AddEvent("shift")
	mleft := robotgo.AddEvent("mleft")
	if ctrl == 0 && mleft == 0 {
		return nil
	}
	return errors.New("ctrl and mleft was not 0")
}
