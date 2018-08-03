// +build linux

package impl

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/pkg/errors"
	"gitlab.com/256/Underbot/sys"
)

// Server is an implementation of Server
type Server struct {
	conn *xgbutil.XUtil // Connection to the X server
}

func (x Server) init() error {
	err := keyInit()
	if err != nil {
		return errors.Wrap(err, "failed to initialize server")
	}
	return nil
}

func (x Server) check() error {

	if x.conn == nil {
		return errors.New("the conn is a nil pointer")
	}
	if x.conn.Conn() == nil {
		return errors.New("The xgb conn is a nil pointer")
	}
	return nil
}

// NewServer returns a server instance
func NewServer() (Server, error) {
	conn, err := xgbutil.NewConn()
	if err != nil {
		return Server{}, errors.Wrap(err, "failed to connect to X server")
	}
	x := Server{conn: conn}
	err = x.check()
	if err != nil {
		return Server{}, errors.Wrap(err, "check for server failed")
	}
	err = x.init()
	if err != nil {
		return Server{}, errors.Wrap(err, "further initialization of server failed")
	}
	return x, nil
}

// ActiveWindow masks the true type into sys.Window to satisfy interface using activeWindow()
func (x Server) ActiveWindow() (sys.Window, error) {
	return x.activeWindow()
}

// ActiveWindow gets the active window, or foreground window
func (x Server) activeWindow() (window, error) {
	// Get the xproto window ID from the active window (the one last clicked usually)
	winID, err := ewmh.ActiveWindowGet(x.conn)
	if err != nil {
		return window{}, errors.Wrap(err, "error getting active window")
	}
	return newWindow(x, winID)
}
