package winmanage

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/BurntSushi/xgbutil"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xtest"

	"gitlab.com/256/Underbot/winmanage/gamestate"

	"github.com/BurntSushi/xgb/xproto"
)

var keycodes map[string]byte

func init() {
	keycodes = make(map[string]byte)
	keycodes["z"] = 52
	keycodes["x"] = 53
}

// UndertaleWindow contains all the necessary info regarding the window the game is located
type UndertaleWindow struct {
	Width    int
	Height   int
	WindowID xproto.Window
	process  process
}

// Create new UndertaleWindow from parameters with checking
func newUndertaleWindow(width, height int, WindowID xproto.Window, Process *os.Process) *UndertaleWindow {
	tempProcess := newProcess(Process)
	window := UndertaleWindow{Width: width, Height: height, WindowID: WindowID, process: tempProcess}
	window.check()
	return &window
}

// Check UndertaleWindow for invalid values
func (window *UndertaleWindow) check() {
	if (window.Width == 0) || (window.Height == 0) {
		panic("size of 0 not permitted")
	}
	if window.WindowID == 0 {
		panic("window has invalid ID")
	}
	window.process.check()
}

// Pause pauses the game
func (window *UndertaleWindow) Pause() {
	if window.process.state != gamestate.Paused {
		window.process.Process.Signal(syscall.SIGSTOP)
		window.process.state = gamestate.Paused
	}
}

// Resume resumes the game
func (window *UndertaleWindow) Resume() {
	if window.process.state == gamestate.Paused {
		window.process.Process.Signal(syscall.SIGCONT)
		window.process.state = gamestate.Playing
	}
}

var cooldown = time.Millisecond * 200
var cooldownActive = false

// Press key in the Undertale window
func (window *UndertaleWindow) Press(x *xgb.Conn, xu *xgbutil.XUtil, key string) {
	if cooldownActive {
		time.Sleep(cooldown)
		cooldownActive = false
		return
	}
	fmt.Println("Press to window")
	prepareXTest(x)
	// 2 is key press
	xtest.FakeInput(x, 2, keycodes["z"], uint32(xu.TimeGet()), window.WindowID, 0, 0, 0)
	// 3 is key release
	xtest.FakeInput(x, 3, keycodes["z"], uint32(xu.TimeGet()), window.WindowID, 0, 0, 0)
}

var prepared = false

func prepareXTest(x *xgb.Conn) {
	if !prepared {
		err := xtest.Init(x)
		if err != nil {
			panic(err)
		}
	}
}
