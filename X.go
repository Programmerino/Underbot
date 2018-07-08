package main

import (
	"image"

	"gitlab.com/256/Underbot/winmanage"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/pkg/errors"
)

// Connections to the X server
var x *xgbutil.XUtil
var xgbConn *xgb.Conn

// Create the connections to the X server
func init() {
	var err error
	x, err = xgbutil.NewConn()
	if err != nil {
		panic(errors.Wrap(err, "Failed to connect to X server"))
	}
	xgbConn, err = xgb.NewConn()
	if err != nil {
		panic(errors.Wrap(err, "Failed to connect to X server"))
	}
}

// GetImage gets a screenshot of the window
func GetImage(window *winmanage.UndertaleWindow) image.RGBA {
	// Gets the image from the window
	ximg, err := xproto.GetImage(xgbConn, xproto.ImageFormatZPixmap, xproto.Drawable(window.WindowID), int16(0), int16(0), uint16(window.Width), uint16(window.Height), 0xffffffff).Reply()
	if err != nil {
		panic(errors.Wrap(err, "Failed to get image"))
	}

	// Converts the xproto image to an image.RGBA instance
	data := ximg.Data
	for i := 0; i < len(data); i += 4 {
		data[i], data[i+2], data[i+3] = data[i+2], data[i], 255
	}
	return image.RGBA{
		Pix:    data,
		Stride: 4 * window.Width,
		Rect:   image.Rect(0, 0, window.Width, window.Height),
	}
}
