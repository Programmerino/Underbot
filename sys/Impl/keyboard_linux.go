// +build linux

package impl

import (
	"fmt"
	"strings"
	"time"

	"github.com/galaktor/gostwriter"
	"github.com/galaktor/gostwriter/key"
	"github.com/pkg/errors"
)

// Get the string for each
var keycodes map[string]*gostwriter.K

var keyboard *gostwriter.Keyboard

func keyInit() error {
	keycodes = make(map[string]*gostwriter.K)

	// The keys that will be needed for pressing
	neededKeys := []key.Code{
		key.CODE_Z, key.CODE_X, key.CODE_UP,
		key.CODE_LEFT, key.CODE_RIGHT, key.CODE_DOWN,
		key.CODE_ENTER}
	// The lowercase string representations of these keys
	stringReps := []string{"z", "x", "up", "left", "right", "down", "enter"}

	// Create keyboard instance
	var err error
	keyboard, err = gostwriter.New(fmt.Sprintf("%s Keyboard", "Underbot"))
	if err != nil {
		return errors.Wrap(err, "failed to create new keyboard")
	}

	// Pre-calculate the K instances for each needed key, and place them in the map
	for i, key := range neededKeys {
		k, err := keyboard.Get(key)
		if err != nil {
			return errors.Wrap(err, "failed to use Get()")
		}

		// Make the string key of the map equal to the K instance pointer
		keycodes[stringReps[i]] = k
	}
	return nil
}

// Press key in the Undertale window (should be lowercase)
func (win window) Press(key string) error {
	var res error
	go func() {
		fmt.Printf("Pressing %s\n", key)
		// Get the ID of the debugging window (could be enhanced by caching this information)
		activeWin, err := win.parent.activeWindow()
		if err != nil {
			res = errors.Wrap(err, "failed to get the active window")
		}
		acID, err := activeWin.ID()
		if err != nil {
			res = errors.Wrap(err, "failed to get the ID of the active window")
		}
		winID, err := win.ID()
		if err != nil {
			res = errors.Wrap(err, "failed to get the ID of the window")
		}
		if acID == winID {
			err := win.justPress(key)
			if err != nil {
				res = errors.Wrap(err, "failed to use justPress")
			}
		} else {
			fmt.Println("Refocusing")
			err = win.SetActive()
			if err != nil {
				res = errors.Wrap(err, "failed to set the active window")
			}
			time.Sleep(time.Millisecond * 250)
			err := win.justPress(key)
			if err != nil {
				res = errors.Wrap(err, "failed to use justPress")
			}
		}
	}()
	return res
}

// No active window handling, just pressing the key on whatever window is active
func (win window) justPress(key string) error {
	lower := strings.ToLower(key)
	if keycodes[lower] == nil {
		return errors.New("the key given was not one included in the keycodes map")
	}
	// Type the key
	err := keycodes[lower].Press()
	if err != nil {
		return errors.Wrap(err, "failed to push the key")
	}
	time.Sleep(time.Millisecond * 40)
	err = keycodes[lower].Release()
	if err != nil {
		return errors.Wrap(err, "failed to release the key")
	}
	return nil
}
