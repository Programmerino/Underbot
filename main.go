package main

import (
	"fmt"
	"image"
	"runtime/pprof"
	"strings"

	"github.com/hajimehoshi/ebiten/inpututil"

	"gitlab.com/256/Underbot/ai"
	"gitlab.com/256/Underbot/cv/object"
	"gitlab.com/256/Underbot/sys"
	impl "gitlab.com/256/Underbot/sys/Impl"

	"gitlab.com/256/Underbot/cv"

	"gitlab.com/256/Underbot/winmanage"

	"github.com/pkg/errors"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

// Counts how many times the debugPrint command has been used
var prints = 0

// The UndertaleWindow instance that will be worked upon
var mainWindow sys.Window

// The title of the debugging window
const title = "Underbot"

// A slice of keys that should be forwarded to the game
var keyForwards = []ebiten.Key{
	ebiten.KeyZ,
	ebiten.KeyX,
	ebiten.KeyEnter,
	ebiten.KeyUp,
	ebiten.KeyDown,
	ebiten.KeyLeft,
	ebiten.KeyRight,
}

// Draws the screen cast of the Undertale window to the screen with the CV information added
func update(screen *ebiten.Image) error {
	if !ebiten.IsRunningSlowly() {
		prints = 0

		img, err := screenCast()
		if err != nil {
			return errors.Wrap(err, "failed to get the image from the window")
		}

		// Create ebiten.Image from the image holding the undertale window with CV drawings on it
		window, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		if err != nil {
			return errors.Wrap(err, "failed to make image from image")
		}

		// Draw the image to the screen
		err = screen.DrawImage(window, &ebiten.DrawImageOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to draw the final image to the screen")
		}

		err = printDebugInfo(screen)
		if err != nil {
			return errors.Wrap(err, "failed to print debugging information to the screen")
		}
		err = handleInput(screen)
		if err != nil {
			return errors.Wrap(err, "failed to handle user input")
		}
	}
	return nil
}

// Prints text to the screen, but keeps track of newlines to avoid overlapping text
func debugPrint(image *ebiten.Image, str string) error {
	prints++
	err := ebitenutil.DebugPrint(image, strings.Repeat("\n", prints)+str)
	if err != nil {
		return errors.Wrap(err, "failed to print message to screen")
	}
	return nil
}

// Prints various important details to the screen for debugging
func printDebugInfo(screen *ebiten.Image) error {
	err := debugPrint(screen, fmt.Sprintf("FPS: %v", ebiten.CurrentFPS()))
	if err != nil {
		return errors.Wrap(err, "failed to print FPS")
	}
	if ai.Disabled {
		err := debugPrint(screen, "State: DISABLED")
		if err != nil {
			return errors.Wrap(err, "failed to print the disabled state")
		}
	} else {
		err := debugPrint(screen, fmt.Sprintf("State: %s", ai.CurrentState.Name))
		if err != nil {
			return errors.Wrap(err, "failed to print the current state")
		}
	}
	err = debugPrint(screen, "Click over an object to get information.")
	if err != nil {
		return errors.Wrap(err, "failed to print instructions")
	}
	err = debugPrint(screen, "Press P to pause the game")
	if err != nil {
		return errors.Wrap(err, "failed to print instructions")
	}
	err = debugPrint(screen, "Press R to resume the game")
	if err != nil {
		return errors.Wrap(err, "failed to print instructions")
	}
	err = debugPrint(screen, "Press A to toggle the AI")
	if err != nil {
		return errors.Wrap(err, "failed to print instructions")
	}
	err = debugPrint(screen, "Press G to toggle the grid")
	if err != nil {
		return errors.Wrap(err, "failed to print instructions")
	}
	err = debugPrint(screen, "Other recognized keys will be forwarded to the game")
	if err != nil {
		return errors.Wrap(err, "failed to print instructions")
	}
	if ai.Disabled {
		err = debugPrint(screen, "The AI is currently DISABLED")
		if err != nil {
			return errors.Wrap(err, "failed to print AI status")
		}
	} else {
		err = debugPrint(screen, "The AI is currently ENABLED")
		if err != nil {
			return errors.Wrap(err, "failed to print AI status")
		}
	}

	for _, recogObj := range cv.GetRecognizedObjects() {
		err = debugPrint(screen, fmt.Sprintf("Recognized %s in object %v", recogObj.RecogObj.Type.Name, recogObj.ID))
		if err != nil {
			return errors.Wrap(err, "failed to print recognized object")
		}
	}
	return nil
}

func handleInput(screen *ebiten.Image) error {
	// Shows the parent objects of the location where the pointer is and debugging information about those objects
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cursorPoint := image.Point{x, y}
		parents := allParents(cursorPoint, cv.GetObjects())
		for _, parent := range parents {
			rect := parent.Bounds
			if parent.Recognized {
				parentText := fmt.Sprintf("Parent %v (recognized as %s): %v x %v",
					parent.ID, parent.RecogObj.Type.Name,
					rect.Dx(), rect.Dy())
				err := debugPrint(screen, parentText)
				if err != nil {
					return errors.Wrap(err, "failed to print parent object with recognition")
				}
			} else {
				err := debugPrint(screen, fmt.Sprintf("Parent %v: %v x %v", parent.ID, rect.Dx(), rect.Dy()))
				if err != nil {
					return errors.Wrap(err, "failed to print parent object")
				}
			}
		}
	}

	// Handles keypresses
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		err := mainWindow.Pause() // Pause the game
		if err != nil {
			return errors.Wrap(err, "failed to pause the game")
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		err := mainWindow.Resume() // Resume the game after a pause
		if err != nil {
			return errors.Wrap(err, "failed to resume the game")
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		ai.Disabled = !ai.Disabled
	} else if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		ai.GridShow = !ai.GridShow
	}

	for _, key := range keyForwards {
		if inpututil.IsKeyJustPressed(key) {
			err := debugPrint(screen, fmt.Sprintf("Forwarding %s key", key.String()))
			if err != nil {
				return errors.Wrap(err, "failed to print key forward debug message")
			}
			err = mainWindow.Press(key.String())
			if err != nil {
				return errors.Wrap(err, "failed to forward key")
			}
		}
	}
	return nil
}

// Gets all the objects that the 'point' is within
func allParents(point image.Point, objs []object.Object) []object.Object {
	var parents = []object.Object{}
	for _, obj := range objs {
		if point.In(obj.Bounds) {
			parents = append(parents, obj)
		}
	}
	return parents
}

// Gets the image from the window, and runs the ProcessImage function on it
func screenCast() (image.Image, error) {
	image, err := mainWindow.GetImage()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the image from the window")
	}
	err = cv.ProcessImage(&image, mainWindow)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process the image")
	}
	return &image, nil
}

/*
 Handles main execution.
 -cpuprofile and -memprofile can be used for profiling to a file
*/
func main() {
	// Profiling
	defer pprof.StopCPUProfile()
	defer func() {
		err := memFile.Close()
		if err != nil {
			panic(errors.Wrap(err, "failed to close the memory file"))
		}
	}()
	err := HandleProfiling()
	if err != nil {
		panic(errors.Wrap(err, "failed to profile the application"))
	}

	serv, err := impl.NewServer()
	if err != nil {
		panic(errors.Wrap(err, "failed to find/get a server for use"))
	}

	// Set the mainWindow to the Window instance from the Get function
	mainWindow, err = winmanage.Get(serv, title)
	if err != nil {
		panic(errors.Wrap(err, "failed to get the window"))
	}

	ebiten.SetRunnableInBackground(true)
	width, height, err := mainWindow.WxH()
	if err != nil {
		panic(errors.Wrap(err, "failed to get the height and width of the window"))
	}
	err = ebiten.Run(update, width, height, 1, title)
	if err != nil {
		panic(errors.Wrap(err, "failed to run the ebiten gui"))
	}
}
