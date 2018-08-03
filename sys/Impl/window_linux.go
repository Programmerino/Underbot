// +build linux

package impl

import (
	"image"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/pkg/errors"
)

// window is an implementation of Window
type window struct {
	parent Server          // The server in charge of the window
	winID  xproto.Window   // The xproto ID of the window
	xWinID *xwindow.Window // The xwindow ID of the window
}

// Newwindow creates a new window instance
func newWindow(x Server, winID xproto.Window) (window, error) {
	// Create xwindow instance from xproto window id
	xWinID := xwindow.New(x.conn, winID)

	xWin := window{parent: x, winID: winID, xWinID: xWinID}
	return xWin, xWin.check()
}

// check ensures that the window is valid
func (xWin window) check() error {
	err := xWin.parent.check()
	if err != nil {
		return errors.Wrap(err, "the parent server was invalid")
	}
	if xWin.winID == 0 {
		return errors.New("window id is zero value")
	}
	if (xWin.xWinID == &xwindow.Window{}) {
		return errors.New("xwindow is empty")
	}
	return nil
}

// Center finds the point in the middle of a xrect
func (xWin window) Center() (point image.Point, err error) {
	rect, err := xWin.rect()
	if err != nil {
		return image.Point{}, errors.Wrap(err, "failed to get the rectangle of the window")
	}
	point.X = rect.X() + (rect.Width() / 2)
	point.Y = rect.Y() + (rect.Height() / 2)
	return point, nil
}

func (xWin window) rect() (xrect.Rect, error) {
	// Get a rectangle describing the dimensions of the window for proper rendering in the case the resize fails silently
	rect, err := xWin.xWinID.Geometry()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the rectangle of the window")
	}
	return rect, nil
}

// WxH gets the width and height of the window
func (xWin window) WxH() (int, int, error) {
	rect, err := xWin.rect()
	if err != nil {
		return 0, 0, errors.Wrap(err, "could not get the rectangle of the window")
	}
	return rect.Width(), rect.Height(), nil
}

// GetImage gets a screenshot of the window
func (xWin window) GetImage() (image.RGBA, error) {
	// Get rectangle of the window
	rect, err := xWin.rect()
	if err != nil {
		return image.RGBA{}, errors.Wrap(err, "failed to get the rectangle of the window")
	}

	// Gets the image from the window
	ximg, err := xproto.GetImage(
		xWin.parent.conn.Conn(), xproto.ImageFormatZPixmap,
		xproto.Drawable(xWin.winID), int16(0), int16(0),
		uint16(rect.Width()), uint16(rect.Height()), 0xffffffff).Reply()
	if err != nil {
		return image.RGBA{}, errors.Wrap(err, "failed to get the image of the window")
	}

	// Converts the xproto image to an image.RGBA instance
	data := ximg.Data
	for i := 0; i < len(data); i += 4 {
		data[i], data[i+2], data[i+3] = data[i+2], data[i], 255
	}
	return image.RGBA{
		Pix:    data,
		Stride: 4 * rect.Width(),
		Rect:   image.Rect(0, 0, rect.Width(), rect.Height()),
	}, nil
}

// Name gets the name of the window
func (xWin window) Name() (string, error) {
	// Get the title of the window
	name, err := ewmh.WmNameGet(xWin.parent.conn, xWin.winID)
	if err != nil || len(name) == 0 {
		name, err = icccm.WmNameGet(xWin.parent.conn, xWin.winID)
		if err != nil || len(name) == 0 {
			return "", errors.Wrap(err, "failed to get the name of the window")
		}
	}
	return name, nil
}

// Resize adjusts the height and width of the window
func (xWin window) Resize(height, width int) error {
	// Attempt to resize the window for WM compatible DEs
	err := xWin.xWinID.WMResize(width, height)
	if err != nil {
		return errors.Wrap(err, "Could not resize window")
	}
	return nil
}

// SetActive makes the window active, or puts the window into the foreground
func (xWin window) SetActive() error {
	err := ewmh.ActiveWindowReq(xWin.parent.conn, xWin.winID)
	if err != nil {
		return errors.Wrap(err, "failed to focus the window")
	}
	return nil
}

// SetActive makes the window active, or puts the window into the foreground
func (xWin window) ID() (int, error) {
	return int(xWin.winID), nil
}
