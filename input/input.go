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
	go listenForControllerEvents()
}

func listenForKeyboardEvents() {
	previousKeyState := make([]uint8, sdl.NUM_SCANCODES)

	keyMappings := map[sdl.Scancode]string{
		sdl.SCANCODE_DOWN: "DOWN",
		sdl.SCANCODE_UP:   "UP",
		sdl.SCANCODE_A:    "A",
		sdl.SCANCODE_B:    "B",
	}

	for {
		currentKeyState := sdl.GetKeyboardState()

		for scancode, keyCode := range keyMappings {
			if currentKeyState[scancode] != 0 && previousKeyState[scancode] == 0 {
				InputChannel <- InputEvent{KeyCode: keyCode}
				PlaySound(getRespectiveSound(keyCode), 10, false)
			}
		}

		copy(previousKeyState, currentKeyState)
		sdl.Delay(50)
	}
}

func listenForControllerEvents() {
	controller := openController()
	defer controller.Close()

	controllerMappings := map[sdl.GameControllerButton]string{
		sdl.CONTROLLER_BUTTON_DPAD_DOWN: "DOWN",
		sdl.CONTROLLER_BUTTON_DPAD_UP:   "UP",
		sdl.CONTROLLER_BUTTON_A:         "A",
		sdl.CONTROLLER_BUTTON_B:         "B",
	}

	for {
		sdl.PumpEvents()
		for button, keyCode := range controllerMappings {
			if controller.Button(button) == sdl.PRESSED {
				InputChannel <- InputEvent{KeyCode: keyCode}
				PlaySound(getRespectiveSound(keyCode), 10, false)
			}
		}

		sdl.Delay(150)
	}
}

func openController() *sdl.GameController {
	for i := 0; i < sdl.NumJoysticks(); i++ {
		if sdl.IsGameController(i) {
			controller := sdl.GameControllerOpen(i)
			if controller != nil {
				return controller
			}
		}
	}
	return nil
}

func getRespectiveSound(key string) string {
	soundMappings := map[string]string{
		"DOWN": "assets/sounds/SFX_UI_MenuSelections.wav",
		"UP":   "assets/sounds/SFX_UI_MenuSelections.wav",
		"A":    "assets/sounds/SFX_UI_Confirm.wav",
		"B":    "assets/sounds/SFX_UI_Cancel.wav",
	}
	return soundMappings[key]
}
