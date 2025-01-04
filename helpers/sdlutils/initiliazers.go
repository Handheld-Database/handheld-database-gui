package sdlutils

import (
	"handheldui/output"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func InitSDL() error {
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_JOYSTICK | sdl.INIT_GAMECONTROLLER); err != nil {
		return output.Errorf("error initializing SDL: %w", err)
	}
	return nil
}

func InitTTF() error {
	if err := ttf.Init(); err != nil {
		return output.Errorf("error initializing SDL_ttf: %w", err)
	}
	return nil
}

func InitMixer() error {
	if err := mix.Init(mix.INIT_MP3 | mix.INIT_OGG); err != nil {
		return output.Errorf("failed to initialize mixer: %w", err)
	}
	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		return output.Errorf("failed to open audio: %w", err)
	}
	mix.Volume(-1, mix.MAX_VOLUME)
	return nil
}

// add font size and return font
func InitFont(fontTtf []byte, font **ttf.Font, size int) error {
	rwops, err := sdl.RWFromMem(fontTtf)
	if err != nil {
		return output.Errorf("error creating RWops from memory: %w", err)
	}
	f, err := ttf.OpenFontRW(rwops, 1, size)
	if err != nil {
		return output.Errorf("error loading font: %w", err)
	}
	*font = f
	return nil
}
