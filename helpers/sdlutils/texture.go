package sdlutils

import (
	"fmt"
	"handheldui/output"
	"handheldui/vars"
	"log"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func getImagePath(str string) string {
	fmt.Println(vars.Config.Screen.AspectRatio)
	return strings.ReplaceAll(str, "$aspect_ratio", vars.Config.Screen.AspectRatio)
}

// LoadTexture loads an image and creates an SDL texture from it
func LoadTexture(renderer *sdl.Renderer, imagePath string) (*sdl.Texture, error) {
	imgSurface, err := img.Load(getImagePath(imagePath))
	if err != nil {
		return nil, output.Errorf("error loading image: %w", err)
	}
	defer imgSurface.Free()

	texture, err := renderer.CreateTextureFromSurface(imgSurface)
	if err != nil {
		return nil, output.Errorf("error creating texture: %w", err)
	}
	return texture, nil
}

// RenderTexture renders a texture with optional scaling and positioning
func RenderTexture(renderer *sdl.Renderer, imagePath string, destRect sdl.Rect, scale float64, crop bool) error {
	// Load the texture
	texture, err := LoadTexture(renderer, imagePath)
	if err != nil {
		return fmt.Errorf("Error loading texture: %w", err)
	}
	defer texture.Destroy()

	// Get texture dimensions
	_, _, width, height, err := texture.Query()
	if err != nil {
		return fmt.Errorf("Error querying texture: %w", err)
	}

	// Apply scaling if necessary
	if scale != 1 {
		width = int32(float64(width) * scale)
		height = int32(float64(height) * scale)
		destRect.W = width
		destRect.H = height
	}

	// Adjust the texture with cropping if necessary
	var srcRect sdl.Rect
	if crop {
		textureWidth, textureHeight := float64(width), float64(height)
		screenAspect := float64(destRect.W) / float64(destRect.H)
		textureAspect := textureWidth / textureHeight

		if textureAspect > screenAspect {
			newWidth := int32(textureHeight * screenAspect)
			srcRect = sdl.Rect{
				X: int32((textureWidth - float64(newWidth)) / 2),
				Y: 0,
				W: newWidth,
				H: int32(textureHeight),
			}
		} else {
			newHeight := int32(textureWidth / screenAspect)
			srcRect = sdl.Rect{
				X: 0,
				Y: int32((textureHeight - float64(newHeight)) / 2),
				W: int32(textureWidth),
				H: newHeight,
			}
		}
	} else {
		// Otherwise, use the entire texture
		srcRect = sdl.Rect{X: 0, Y: 0, W: width, H: height}
	}

	// Render the texture
	return renderer.Copy(texture, &srcRect, &destRect)
}

func RenderTextureCartesian(renderer *sdl.Renderer, imagePath string, startQuadrant, endQuadrant string) {

	screenWidth, screenHeight := vars.Config.Screen.Width, vars.Config.Screen.Height

	// Define the quadrants
	quadrants := map[string]sdl.Rect{
		"Q1": {X: screenWidth / 2, Y: 0, W: screenWidth / 2, H: screenHeight / 2},
		"Q2": {X: 0, Y: 0, W: screenWidth / 2, H: screenHeight / 2},
		"Q3": {X: 0, Y: screenHeight / 2, W: screenWidth / 2, H: screenHeight / 2},
		"Q4": {X: screenWidth / 2, Y: screenHeight / 2, W: screenWidth / 2, H: screenHeight / 2},
	}

	startRect, startOk := quadrants[startQuadrant]
	endRect, endOk := quadrants[endQuadrant]

	if !startOk || !endOk {
		output.Errorf("Unknown quadrant(s): %s, %s\n", startQuadrant, endQuadrant)
		return
	}

	// Calculate the combined rectangle
	dstRect := sdl.Rect{
		X: min(startRect.X, endRect.X),
		Y: min(startRect.Y, endRect.Y),
		W: max(startRect.X+startRect.W, endRect.X+endRect.W) - min(startRect.X, endRect.X),
		H: max(startRect.Y+startRect.H, endRect.Y+endRect.H) - min(startRect.Y, endRect.Y),
	}

	// Use RenderTexture to render
	err := RenderTexture(renderer, imagePath, dstRect, 1, false)
	if err != nil {
		output.Errorf("Error rendering texture: %v", err)
	}
}

func RenderTextureAdjusted(renderer *sdl.Renderer, imagePath string, rect sdl.Rect) {
	// Use RenderTexture to render
	err := RenderTexture(renderer, imagePath, rect, 1, false)
	if err != nil {
		output.Errorf("Error rendering texture: %v", err)
	}
}

func RenderScaledTexture(renderer *sdl.Renderer, imgPath string, x, y int32, scale float64) {
	// Load the texture
	texture, err := LoadTexture(renderer, imgPath)
	if err != nil {
		log.Printf("Error loading texture: %s, %v", imgPath, err)
		return
	}
	defer texture.Destroy()

	// Get texture dimensions
	_, _, width, height, err := texture.Query()
	if err != nil {
		log.Printf("Error querying texture: %v", err)
		return
	}

	// Apply scaling
	scaledWidth := int32(float64(width) * scale)
	scaledHeight := int32(float64(height) * scale)

	// Adjust the destination rectangle (centering)
	dstRect := sdl.Rect{
		X: x - scaledWidth/2,
		Y: y - scaledHeight/2,
		W: scaledWidth,
		H: scaledHeight,
	}

	// Use RenderTexture to render
	err = RenderTexture(renderer, imgPath, dstRect, scale, false)
	if err != nil {
		log.Printf("Error rendering texture: %v", err)
	}
}

func RenderTextureCover(renderer *sdl.Renderer, imagePath string) {
	// Get screen dimensions
	screenWidth, screenHeight := vars.Config.Screen.Width, vars.Config.Screen.Height

	// Load and calculate the aspect of the image
	textureSurface, err := sdl.LoadBMP(imagePath)
	if err != nil {
		fmt.Printf("Error loading texture image: %v\n", err)
		return
	}
	defer textureSurface.Free()

	textureTexture, err := renderer.CreateTextureFromSurface(textureSurface)
	if err != nil {
		fmt.Printf("Error creating texture from image: %v\n", err)
		return
	}
	defer textureTexture.Destroy()

	dstRect := sdl.Rect{X: 0, Y: 0, W: int32(screenWidth), H: int32(screenHeight)}

	// Use RenderTexture to render
	err = RenderTexture(renderer, imagePath, dstRect, 1, true)
	if err != nil {
		fmt.Printf("Error rendering texture: %v\n", err)
	}
}
