package vars

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type FontColors struct {
	WHITE     sdl.Color
	PRIMARY   sdl.Color
	SECONDARY sdl.Color
}

var (
	ScreenWidth     int32
	ScreenHeight    int32
	CurrentPlatform string
	CurrentScreen   string
	CurrentSystem   string
	CurrentGame     string
	BodyFont        *ttf.Font
	HeaderFont      *ttf.Font
	Colors          FontColors
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
	Colors = FontColors{
		WHITE:     sdl.Color{R: 255, G: 255, B: 255, A: 255},
		PRIMARY:   sdl.Color{R: 255, G: 214, B: 255, A: 255},
		SECONDARY: sdl.Color{R: 231, G: 192, B: 255, A: 255},
	}
}
