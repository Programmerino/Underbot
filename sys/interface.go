// Package sys provides a way of supporting new window protocols easy to do; just satisfy the interfaces!
package sys

import (
	"image"
	"os"
)

// Server provides an interface for the top-level of the protocol. Think X11 Server
type Server interface {
	ActiveWindow() (Window, error) // Get the active window
}

// Window is an instance of a window such as Chrome
type Window interface {
	GetImage() (image.RGBA, error) // Should return image of the window
	// Should return a point with the X-coordinate referring to the width, and Y with height
	Center() (image.Point, error)
	Process() (*os.Process, error)  // Returns the process of the window
	Name() (string, error)          // Returns the name of the window (titlebar)
	Resize(width, height int) error // Resizes the window to the height and width specified
	SetActive() error               // Makes the window foreground/active
	Pause() error                   // Should pause the game
	Resume() error                  // Should resume the game
	Press(string) error             // Emulates a key press
	WxH() (int, int, error)         // Gets the width and height of the window
	// An ID tied to the underlying window in some way
	// For example, for checking if two window instances are referring to the same window
	ID() (int, error)
}
