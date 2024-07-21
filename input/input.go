package input

import (
	"github.com/veandco/go-sdl2/sdl"
)

type InputEvent struct {
	KeyCode string
}

var InputChannel = make(chan InputEvent)

func StartListening() {
	go listenForKeyboardEvents()
}

func listenForKeyboardEvents() {
	for {
		keyState := sdl.GetKeyboardState()

		if keyState[sdl.SCANCODE_DOWN] != 0 {
			InputChannel <- InputEvent{KeyCode: "DOWN"}
			PlaySound("assets/sounds/SFX_UI_MenuSelections.wav", 10, false)
		}
		if keyState[sdl.SCANCODE_UP] != 0 {
			InputChannel <- InputEvent{KeyCode: "UP"}
			PlaySound("assets/sounds/SFX_UI_MenuSelections.wav", 10, false)
		}
		if keyState[sdl.SCANCODE_A] != 0 {
			InputChannel <- InputEvent{KeyCode: "A"}
			PlaySound("assets/sounds/SFX_UI_Confirm.wav", 10, false)
		}
		if keyState[sdl.SCANCODE_B] != 0 {
			InputChannel <- InputEvent{KeyCode: "B"}
			PlaySound("assets/sounds/SFX_UI_Cancel.wav", 10, false)
		}

		sdl.Delay(150)
	}
}
