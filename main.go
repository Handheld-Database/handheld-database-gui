package main

import (
	_ "embed"
	"log"
	"os"
	"runtime/debug"

	"handheldui/helpers/sdlutils"
	"handheldui/input"
	"handheldui/output"
	"handheldui/screens"
	"handheldui/vars"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

//go:embed assets/fonts/Kenney-Pixel-Square.ttf
var kenneyPixelSquare []byte

//go:embed assets/fonts/Kenney-Space.ttf
var KenneySpace []byte

//go:embed assets/fonts/NotoSans_Condensed-SemiBold.ttf
var NotoSans []byte

//go:embed configs/config.json
var ConfigFile []byte

func main() {
	// Defer a function to handle panics and exit with -1
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Unhandled error: %v\n", r)
			log.Println("Stack trace:")
			debug.PrintStack()
			os.Exit(-1)
		}
	}()

	vars.InitVars()

	var err error

	vars.Config, err = vars.LoadConfig(ConfigFile)
	if err != nil {
		panic(err)
	}

	vars.ScreenWidth = int32(1280)
	vars.ScreenHeight = int32(720)

	if err := sdlutils.InitSDL(); err != nil {
		panic(err)
	}

	if err := sdlutils.InitMixer(); err != nil {
		panic(err)
	}
	defer mix.CloseAudio()

	if err := sdlutils.InitTTF(); err != nil {
		panic(err)
	}

	if err := sdlutils.InitFont(NotoSans, &vars.BodyFont, 24); err != nil {
		panic(err)
	}

	if err := sdlutils.InitFont(kenneyPixelSquare, &vars.BodyBigFont, 58); err != nil {
		panic(err)
	}

	if err := sdlutils.InitFont(NotoSans, &vars.LongTextFont, 24); err != nil {
		panic(err)
	}

	if err := sdlutils.InitFont(KenneySpace, &vars.HeaderFont, 28); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Systems List", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, vars.ScreenWidth, vars.ScreenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	output.PlaySound("assets/sounds/Retro_Mystic.ogg", 5, true)

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	homeScreen, err := screens.NewHomeScreen(renderer)
	if err != nil {
		panic(err)
	}

	repositoriesScreen, err := screens.NewRepositoriesScreen(renderer)
	if err != nil {
		panic(err)
	}

	filesScreen, err := screens.NewFilesScreen(renderer)
	if err != nil {
		panic(err)
	}

	systemsScreen, err := screens.NewSystemsScreen(renderer)
	if err != nil {
		panic(err)
	}

	gamesScreen, err := screens.NewGamesScreen(renderer)
	if err != nil {
		panic(err)
	}

	reviewsScreen, err := screens.NewReviewsScreen(renderer)
	if err != nil {
		panic(err)
	}

	overviewScreen, err := screens.NewOverviewScreen(renderer)
	if err != nil {
		panic(err)
	}

	screensMap := map[string]func(){
		"home_screen":         homeScreen.Draw,
		"repositories_screen": repositoriesScreen.Draw,
		"files_screen":        filesScreen.Draw,
		"systems_screen":      systemsScreen.Draw,
		"games_screen":        gamesScreen.Draw,
		"overview_screen":     overviewScreen.Draw,
		"reviews_screen":      reviewsScreen.Draw,
	}

	inputHandlers := map[string]func(input.InputEvent){
		"home_screen":         homeScreen.HandleInput,
		"repositories_screen": repositoriesScreen.HandleInput,
		"files_screen":        filesScreen.HandleInput,
		"systems_screen":      systemsScreen.HandleInput,
		"games_screen":        gamesScreen.HandleInput,
		"overview_screen":     overviewScreen.HandleInput,
		"reviews_screen":      reviewsScreen.HandleInput,
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
			// No event received
		}

		if drawFunc, ok := screensMap[vars.CurrentScreen]; ok {
			drawFunc()
		}

		sdl.Delay(16)
	}
}
