package sdlutils

import (
	"handheldui/output"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// LoadFont loads a font from RWops and returns the font object
func LoadFont(rwops *sdl.RWops, size int) (*ttf.Font, error) {
	font, err := ttf.OpenFontRW(rwops, 1, size)
	if err != nil {
		return nil, output.Errorf("error loading font: %w", err)
	}
	return font, nil
}

// DrawText is a function that draws text on the screen based on the provided position, color, and font.
func DrawText(renderer *sdl.Renderer, text string, position sdl.Point, color sdl.Color, font *ttf.Font) {
	// Render the text to a surface
	textSurface, err := RenderText(text, color, font)
	if err != nil {
		output.Errorf("Error rendering text: %v\n", err)
		return
	}
	defer textSurface.Free()

	// Create a texture from the surface
	textTexture, err := renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		output.Errorf("Error creating texture: %v\n", err)
		return
	}
	defer textTexture.Destroy()

	// Set the destination rectangle for the texture
	destinationRect := sdl.Rect{
		X: position.X,
		Y: position.Y,
		W: int32(textSurface.W),
		H: int32(textSurface.H),
	}

	// Copy the texture to the renderer
	renderer.Copy(textTexture, nil, &destinationRect)
}

// RenderText renders text to an SDL surface
func RenderText(text string, color sdl.Color, font *ttf.Font) (*sdl.Surface, error) {
	textSurface, err := font.RenderUTF8Blended(text, color)
	if err != nil {
		return nil, output.Errorf("error rendering text: %w", err)
	}
	return textSurface, nil
}
