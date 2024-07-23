package helpers

import (
	"fmt"
	"handheldui/output"
	"handheldui/vars"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func InitSDL() error {
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_JOYSTICK | sdl.INIT_GAMECONTROLLER); err != nil {
		return fmt.Errorf("error initializing SDL: %w", err)
	}
	return nil
}

func InitTTF() error {
	if err := ttf.Init(); err != nil {
		return fmt.Errorf("error initializing SDL_ttf: %w", err)
	}
	return nil
}

func InitMixer() error {
	if err := mix.Init(mix.INIT_MP3 | mix.INIT_OGG); err != nil {
		return fmt.Errorf("failed to initialize mixer: %w", err)
	}
	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		return fmt.Errorf("failed to open audio: %w", err)
	}
	mix.Volume(-1, mix.MAX_VOLUME)
	return nil
}

// add font size and return font
func InitFont(fontTtf []byte, font **ttf.Font, size int) error {
	rwops, err := sdl.RWFromMem(fontTtf)
	if err != nil {
		return fmt.Errorf("error creating RWops from memory: %w", err)
	}
	f, err := ttf.OpenFontRW(rwops, 1, size)
	if err != nil {
		return fmt.Errorf("error loading font: %w", err)
	}
	*font = f
	return nil
}

// LoadTexture loads an image and creates an SDL texture from it
func LoadTexture(renderer *sdl.Renderer, imagePath string) (*sdl.Texture, error) {
	imgSurface, err := img.Load(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error loading image: %w", err)
	}
	defer imgSurface.Free()

	texture, err := renderer.CreateTextureFromSurface(imgSurface)
	if err != nil {
		return nil, fmt.Errorf("error creating texture: %w", err)
	}
	return texture, nil
}

// LoadFont loads a font from RWops and returns the font object
func LoadFont(rwops *sdl.RWops, size int) (*ttf.Font, error) {
	font, err := ttf.OpenFontRW(rwops, 1, size)
	if err != nil {
		return nil, fmt.Errorf("error loading font: %w", err)
	}
	return font, nil
}

// DrawText is a function that draws text on the screen based on the provided position, color, and font.
func DrawText(renderer *sdl.Renderer, text string, position sdl.Point, color sdl.Color, font *ttf.Font) {
	// Render the text to a surface
	textSurface, err := RenderText(text, color, font)
	if err != nil {
		output.Printf("Error rendering text: %v\n", err)
		return
	}
	defer textSurface.Free()

	// Create a texture from the surface
	textTexture, err := renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		output.Printf("Error creating texture: %v\n", err)
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
		return nil, fmt.Errorf("error rendering text: %w", err)
	}
	return textSurface, nil
}

func RenderTexture(renderer *sdl.Renderer, imagePath string, startQuadrant, endQuadrant string) {
	// Load the texture image
	textureSurface, err := sdl.LoadBMP(imagePath)
	if err != nil {
		output.Printf("Error loading texture image: %v\n", err)
		return
	}
	defer textureSurface.Free()

	textureTexture, err := renderer.CreateTextureFromSurface(textureSurface)
	if err != nil {
		output.Printf("Error creating texture from image: %v\n", err)
		return
	}
	defer textureTexture.Destroy()

	// Get screen width and height
	screenWidth, screenHeight := vars.ScreenWidth, vars.ScreenHeight
	halfWidth, halfHeight := screenWidth/2, screenHeight/2

	// Define rectangles for each quadrant
	quadrants := map[string]sdl.Rect{
		"Q1": {X: halfWidth, Y: 0, W: halfWidth, H: halfHeight},          // Q1
		"Q2": {X: 0, Y: 0, W: halfWidth, H: halfHeight},                  // Q2
		"Q3": {X: 0, Y: halfHeight, W: halfWidth, H: halfHeight},         // Q3
		"Q4": {X: halfWidth, Y: halfHeight, W: halfWidth, H: halfHeight}, // Q4
	}

	// Check if the quadrants are valid
	startRect, startOk := quadrants[startQuadrant]
	endRect, endOk := quadrants[endQuadrant]

	if !startOk || !endOk {
		output.Printf("Unknown quadrant(s): %s, %s\n", startQuadrant, endQuadrant)
		return
	}

	// Calculate the rectangle covering the area between the quadrants
	dstRect := sdl.Rect{
		X: min(startRect.X, endRect.X),
		Y: min(startRect.Y, endRect.Y),
		W: max(startRect.X+startRect.W, endRect.X+endRect.W) - min(startRect.X, endRect.X),
		H: max(startRect.Y+startRect.H, endRect.Y+endRect.H) - min(startRect.Y, endRect.Y),
	}

	// Get the dimensions of the texture
	textureWidth, textureHeight := textureSurface.W, textureSurface.H

	// Calculate the source rectangle of the texture
	srcRect := sdl.Rect{
		X: 0,
		Y: 0,
		W: int32(textureWidth),
		H: int32(textureHeight),
	}

	// Render the texture adjusted to the area between the quadrants
	renderer.Copy(textureTexture, &srcRect, &dstRect)
}

// Helper functions to calculate min and max
func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func RenderTextureAdjusted(renderer *sdl.Renderer, imagePath string, x, y, width, height int32) {
	// Load the texture image
	textureSurface, err := sdl.LoadBMP(imagePath)
	if err != nil {
		output.Printf("Error loading texture image: %v\n", err)
		return
	}
	defer textureSurface.Free()

	textureTexture, err := renderer.CreateTextureFromSurface(textureSurface)
	if err != nil {
		output.Printf("Error creating texture from image: %v\n", err)
		return
	}
	defer textureTexture.Destroy()

	// Draw the texture at the specified position and size
	renderer.Copy(textureTexture, nil, &sdl.Rect{X: x, Y: y, W: width, H: height})
}

// WrapText splits a long text into multiple lines based on the specified maximum width.
func WrapText(text string, font *ttf.Font, maxWidth int) []string {
	words := strings.Fields(text)
	var lines []string
	var currentLine string

	for _, word := range words {
		lineWithWord := currentLine + word + " "
		lineWidth := textWidth(font, lineWithWord)

		if lineWidth > maxWidth {
			if len(currentLine) > 0 {
				lines = append(lines, strings.TrimSpace(currentLine))
			}
			currentLine = word + " "
		} else {
			currentLine = lineWithWord
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, strings.TrimSpace(currentLine))
	}

	return lines
}

// textWidth calculates the width of a string of text based on the provided font.
func textWidth(font *ttf.Font, text string) int {
	surface, err := font.RenderUTF8Blended(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return 0
	}
	defer surface.Free()

	return int(surface.W)
}
