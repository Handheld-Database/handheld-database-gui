package vars

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type FontColors struct {
	WHITE     sdl.Color
	PRIMARY   sdl.Color
	SECONDARY sdl.Color
	BLACK     sdl.Color
}

var (
	CurrentPlatform string
	CurrentScreen   string
	CurrentSystem   string
	CurrentGame     string
	CurrentRepo     string
	CurrentTester   string
	BodyFont        *ttf.Font
	HeaderFont      *ttf.Font
	BodyBigFont     *ttf.Font
	LongTextFont    *ttf.Font
	Colors          FontColors
	Config          *ConfigDefinition
)

func InitVars() {
	Config = nil
	CurrentPlatform = "tsp"
	CurrentScreen = "home_screen"
	CurrentSystem = ""
	CurrentGame = ""
	CurrentRepo = ""
	BodyFont = nil
	HeaderFont = nil
	BodyBigFont = nil
	LongTextFont = nil
	Colors = FontColors{
		WHITE:     sdl.Color{R: 255, G: 255, B: 255, A: 255},
		PRIMARY:   sdl.Color{R: 255, G: 214, B: 255, A: 255},
		SECONDARY: sdl.Color{R: 231, G: 192, B: 255, A: 255},
		BLACK:     sdl.Color{R: 0, G: 0, B: 0, A: 255},
	}
}
