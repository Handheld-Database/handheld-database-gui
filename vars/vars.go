package vars

import "github.com/veandco/go-sdl2/ttf"

var (
	ScreenWidth     int32
	ScreenHeight    int32
	CurrentPlatform string
	CurrentScreen   string
	CurrentSystem   string
	CurrentGame     string
	BodyFont        *ttf.Font
	HeaderFont      *ttf.Font
)

func InitVars() {
	ScreenWidth = 0
	ScreenHeight = 0
	CurrentPlatform = "tsp"
	CurrentScreen = "systems_screen"
	CurrentSystem = ""
	CurrentGame = ""
	BodyFont = nil
	HeaderFont = nil
}
