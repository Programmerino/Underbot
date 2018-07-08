package main

import (
	"fmt"
	"image"
	"runtime/pprof"
	"strings"

	"gitlab.com/256/Underbot/ai"
	"gitlab.com/256/Underbot/cv/object"

	"gitlab.com/256/Underbot/cv"

	"gitlab.com/256/Underbot/winmanage"

	"github.com/pkg/errors"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

// The UndertaleWindow instance that will be worked upon
var mainWindow *winmanage.UndertaleWindow

// Counts how many times the debugPrint command has been used
var prints = 0

// The title of the debugging window
const title = "Underbot"

// Draws the screen cast of the Undertale window to the screen with the CV information added
func update(screen *ebiten.Image) error {
	if !ebiten.IsRunningSlowly() {
		prints = 0

		// Create ebiten.Image from the image holding the undertale window with CV drawings on it
		window, err := ebiten.NewImageFromImage(screenCast(), ebiten.FilterDefault)
		if err != nil {
			panic(errors.Wrap(err, "Failed to make image from image"))
		}

		// Draw the image to the screen
		screen.DrawImage(window, &ebiten.DrawImageOptions{})

		printDebugInfo(screen)
		handleInput(screen)
	}
	return nil
}

// Prints text to the screen, but keeps track of newlines to avoid overlapping text
func debugPrint(image *ebiten.Image, str string) error {
	prints++
	ebitenutil.DebugPrint(image, strings.Repeat("\n", prints)+str)
	return nil
}

// Prints various important details to the screen for debugging
func printDebugInfo(screen *ebiten.Image) {
	debugPrint(screen, fmt.Sprintf("FPS: %v", ebiten.CurrentFPS()))
	debugPrint(screen, fmt.Sprintf("State: %s", ai.CurrentState.Name))
	debugPrint(screen, "Click over an object to get information.")
	debugPrint(screen, "Press P to pause the game")
	debugPrint(screen, "Press R to resume the game")
	debugPrint(screen, "Other recognized keys will be forwarded to the game")
	for _, recogObj := range cv.GetRecognizedObjects() {
		debugPrint(screen, fmt.Sprintf("Recognized %s in object %v", recogObj.RecogObj.Name, recogObj.ID))
	}
}

func handleInput(screen *ebiten.Image) {
	// Shows the parent objects of the location where the pointer is and debugging information about those objects
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cursorPoint := image.Point{x, y}
		parents := allParents(cursorPoint, cv.GetObjects())
		for _, parent := range parents {
			rect := parent.Bounds
			if parent.Recognized {
				debugPrint(screen, fmt.Sprintf("Parent %v (recognized as %s): %v x %v", parent.ID, parent.RecogObj.Name, rect.Dx(), rect.Dy()))
			} else {
				debugPrint(screen, fmt.Sprintf("Parent %v: %v x %v", parent.ID, rect.Dx(), rect.Dy()))
			}
		}
	}

	// Handles keypresses
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		mainWindow.Pause() // Pause the game
	} else if ebiten.IsKeyPressed(ebiten.KeyR) {
		mainWindow.Resume() // Resume the game after a pause
	} else if ebiten.IsKeyPressed(ebiten.KeyX) {
		// Doesn't work yet, but should press x in the game
		debugPrint(screen, "Forwarding X key")
		mainWindow.Press(xgbConn, x, "x")
	}
}

// Gets all the objects that the point is within
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
func screenCast() image.Image {
	image := GetImage(mainWindow)
	return cv.ProcessImage(image)
}

// Handles main execution. -cpuprofile and -memprofile with their respective file locations can be used to profile the cpu and memory usage
func main() {
	// Profiling
	defer pprof.StopCPUProfile()
	defer memFile.Close()
	HandleProfiling()

	// Set the mainWindow to the UndertaleWindow instance from the Get function
	mainWindow = winmanage.Get(x, title)
	ebiten.SetRunnableInBackground(true)
	ebiten.Run(update, mainWindow.Width, mainWindow.Height, 1, title)
}
