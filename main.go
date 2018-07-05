package main

import (
	"fmt"
	"image"
	"strings"

	"gitlab.com/256/Underbot/cv"

	"gitlab.com/256/Underbot/winmanage"

	"gitlab.com/256/Underbot/profiling"

	"github.com/pkg/errors"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

var mainWindow *winmanage.UndertaleWindow
var prints = 0

// Draws the screen cast of the Undertale window to the screen with the CV information added
func update(screen *ebiten.Image) error {
	prints = 0

	if ebiten.IsRunningSlowly() {
		fmt.Println("Running slowly!")
	} else {
		window, err := ebiten.NewImageFromImage(screenCast(), ebiten.FilterDefault)
		if err != nil {
			panic(errors.Wrap(err, "Failed to make image from image"))
		}
		screen.DrawImage(window, &ebiten.DrawImageOptions{})
	}

	printDebugInfo(screen)
	handleInput(screen)
	return nil
}

func DebugPrint(image *ebiten.Image, str string) error {
	prints++
	ebitenutil.DebugPrint(image, strings.Repeat("\n", prints)+str)
	return nil
}

func printDebugInfo(screen *ebiten.Image) {
	DebugPrint(screen, fmt.Sprintf("FPS: %v", ebiten.CurrentFPS()))
	DebugPrint(screen, "Click over an object to get information.")
	DebugPrint(screen, "Press P to pause the game")
	DebugPrint(screen, "Press R to resume the game")
	DebugPrint(screen, "Other recognized keys will be forwarded to the game")
	for _, recogObj := range cv.GetRecognizedObjects() {
		DebugPrint(screen, fmt.Sprintf("Recognized %s in object %v", recogObj.RecogObj.Name, recogObj.ID))
	}
}

//var keyForwards = []ebiten.Key{ebiten.KeyZ, ebiten.KeyX, ebiten.KeyLeft, ebiten.KeyRight}

func handleInput(screen *ebiten.Image) {
	// When the "left mouse button" is pressed...
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cursorPoint := image.Point{x, y}
		parents := allParents(cursorPoint, cv.GetObjects())
		for _, parent := range parents {
			rect := parent.Bounds
			if parent.Recognized {
				DebugPrint(screen, fmt.Sprintf("Parent %v (recognized as %s): %v x %v", parent.ID, parent.RecogObj.Name, rect.Dx(), rect.Dy()))
			} else {
				DebugPrint(screen, fmt.Sprintf("Parent %v: %v x %v", parent.ID, rect.Dx(), rect.Dy()))
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		mainWindow.Pause()
	} else if ebiten.IsKeyPressed(ebiten.KeyR) {
		mainWindow.Resume()
	} else if ebiten.IsKeyPressed(ebiten.KeyX) {
		DebugPrint(screen, "Forwarding X key")
		mainWindow.Press(xgbConn, x, "x")
	}
}

func allParents(point image.Point, objs []cv.Object) []cv.Object {
	var parents = []cv.Object{}
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

func main() {
	profiling.HandleProfiling()
	mainWindow = winmanage.Get(x)
	ebiten.SetRunnableInBackground(true)
	ebiten.Run(update, mainWindow.Width, mainWindow.Height, 1, "Underbot")
}
