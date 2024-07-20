package main

import (
	_ "embed"
	"handheldui/helpers"
	"handheldui/input"
	"handheldui/screens"
	"handheldui/vars"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

//go:embed assets/fonts/Kenney-Pixel-Square.ttf
var kenneyPixelSquare []byte

//go:embed assets/fonts/Kenney-Bold.ttf
var kenneyBold []byte

func main() {
	vars.InitVars()

	vars.ScreenWidth = int32(1280)
	vars.ScreenHeight = int32(720)

	if err := helpers.InitSDL(); err != nil {
		panic(err)
	}

	if err := helpers.InitMixer(); err != nil {
		panic(err)
	}
	defer mix.CloseAudio()

	if err := helpers.InitTTF(); err != nil {
		panic(err)
	}

	if err := helpers.InitFont(kenneyPixelSquare, &vars.BodyFont, 14); err != nil {
		panic(err)
	}

	if err := helpers.InitFont(kenneyBold, &vars.HeaderFont, 18); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Lista de Consoles", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, vars.ScreenWidth, vars.ScreenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	input.PlaySound("assets/sounds/Retro_Mystic.ogg", 5, true)

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	systemsScreen, err := screens.NewSystemsScreen(renderer)
	if err != nil {
		panic(err)
	}
	gamesScreen, err := screens.NewGamesScreen(renderer)
	if err != nil {
		panic(err)
	}
	overviewScreen, err := screens.NewOverviewScreen(renderer)
	if err != nil {
		panic(err)
	}

	screensMap := map[string]func(){
		"systems_screen":  systemsScreen.Draw,
		"games_screen":    gamesScreen.Draw,
		"overview_screen": overviewScreen.Draw,
	}

	inputHandlers := map[string]func(input.InputEvent){
		"systems_screen":  systemsScreen.HandleInput,
		"games_screen":    gamesScreen.HandleInput,
		"overview_screen": overviewScreen.HandleInput,
	}

	input.StartListening()

	running := true
	for running {

		for {
			event := sdl.PollEvent()
			if event == nil {
				break
			}

			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		select {
		case inputEvent := <-input.InputChannel:
			if handler, ok := inputHandlers[vars.CurrentScreen]; ok {
				handler(inputEvent)
			}
		default:
			// Nenhum evento de entrada recebido
		}

		if drawFunc, ok := screensMap[vars.CurrentScreen]; ok {
			drawFunc()
		}

		sdl.Delay(16)
	}
}
