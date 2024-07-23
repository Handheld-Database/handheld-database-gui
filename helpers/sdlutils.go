package helpers

import (
	"fmt"
	"handheldui/vars"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func InitSDL() error {
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
		return fmt.Errorf("erro ao inicializar SDL: %w", err)
	}
	return nil
}

func InitTTF() error {
	if err := ttf.Init(); err != nil {
		return fmt.Errorf("erro ao inicializar SDL_ttf: %w", err)
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

// adicionar font size e retornar font
func InitFont(fontTtf []byte, font **ttf.Font, size int) error {
	rwops, err := sdl.RWFromMem(fontTtf)
	if err != nil {
		return fmt.Errorf("erro ao criar RWops a partir da memória: %w", err)
	}
	f, err := ttf.OpenFontRW(rwops, 1, size)
	if err != nil {
		return fmt.Errorf("erro ao carregar a fonte: %w", err)
	}
	*font = f
	return nil
}

// LoadTexture loads an image and creates an SDL texture from it
func LoadTexture(renderer *sdl.Renderer, imagePath string) (*sdl.Texture, error) {
	imgSurface, err := img.Load(imagePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar imagem: %w", err)
	}
	defer imgSurface.Free()

	texture, err := renderer.CreateTextureFromSurface(imgSurface)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar textura: %w", err)
	}
	return texture, nil
}

// LoadFont loads a font from RWops and returns the font object
func LoadFont(rwops *sdl.RWops, size int) (*ttf.Font, error) {
	font, err := ttf.OpenFontRW(rwops, 1, size)
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar fonte: %w", err)
	}
	return font, nil
}

// DrawText é uma função que desenha um texto na tela com base na posição, cor e fonte fornecidos.
func DrawText(renderer *sdl.Renderer, text string, position sdl.Point, color sdl.Color, font *ttf.Font) {
	// Renderize o texto para uma superfície
	textSurface, err := RenderText(text, color, font)
	if err != nil {
		fmt.Printf("Erro ao renderizar texto: %v\n", err)
		return
	}
	defer textSurface.Free()

	// Crie uma textura a partir da superfície
	textTexture, err := renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		fmt.Printf("Erro ao criar textura: %v\n", err)
		return
	}
	defer textTexture.Destroy()

	// Defina o retângulo de destino para a textura
	destinationRect := sdl.Rect{
		X: position.X,
		Y: position.Y,
		W: int32(textSurface.W),
		H: int32(textSurface.H),
	}

	// Copie a textura para o renderizador
	renderer.Copy(textTexture, nil, &destinationRect)
}

// RenderText renders text to an SDL surface
func RenderText(text string, color sdl.Color, font *ttf.Font) (*sdl.Surface, error) {
	textSurface, err := font.RenderUTF8Blended(text, color)
	if err != nil {
		return nil, fmt.Errorf("erro ao renderizar texto: %w", err)
	}
	return textSurface, nil
}

func RenderTexture(renderer *sdl.Renderer, imagePath string, startQuadrant, endQuadrant string) {
	// Carregar a imagem de textura
	textureSurface, err := sdl.LoadBMP(imagePath)
	if err != nil {
		fmt.Printf("Erro ao carregar imagem de textura: %v\n", err)
		return
	}
	defer textureSurface.Free()

	textureTexture, err := renderer.CreateTextureFromSurface(textureSurface)
	if err != nil {
		fmt.Printf("Erro ao criar textura de textura: %v\n", err)
		return
	}
	defer textureTexture.Destroy()

	// Obter a largura e altura da tela
	screenWidth, screenHeight := vars.ScreenWidth, vars.ScreenHeight
	halfWidth, halfHeight := screenWidth/2, screenHeight/2

	// Definir os retângulos para cada quadrante
	quadrants := map[string]sdl.Rect{
		"Q1": {X: halfWidth, Y: 0, W: halfWidth, H: halfHeight},          // Q1
		"Q2": {X: 0, Y: 0, W: halfWidth, H: halfHeight},                  // Q2
		"Q3": {X: 0, Y: halfHeight, W: halfWidth, H: halfHeight},         // Q3
		"Q4": {X: halfWidth, Y: halfHeight, W: halfWidth, H: halfHeight}, // Q4
	}

	// Verificar se os quadrantes são válidos
	startRect, startOk := quadrants[startQuadrant]
	endRect, endOk := quadrants[endQuadrant]

	if !startOk || !endOk {
		fmt.Printf("Quadrante(s) desconhecido(s): %s, %s\n", startQuadrant, endQuadrant)
		return
	}

	// Calcular o retângulo que cobre a área entre os quadrantes
	dstRect := sdl.Rect{
		X: min(startRect.X, endRect.X),
		Y: min(startRect.Y, endRect.Y),
		W: max(startRect.X+startRect.W, endRect.X+endRect.W) - min(startRect.X, endRect.X),
		H: max(startRect.Y+startRect.H, endRect.Y+endRect.H) - min(startRect.Y, endRect.Y),
	}

	// Obter as dimensões da textura
	textureWidth, textureHeight := textureSurface.W, textureSurface.H

	// Calcular o retângulo de origem da textura
	srcRect := sdl.Rect{
		X: 0,
		Y: 0,
		W: int32(textureWidth),
		H: int32(textureHeight),
	}

	// Renderizar a textura ajustada para a área entre os quadrantes
	renderer.Copy(textureTexture, &srcRect, &dstRect)
}

// Funções auxiliares para calcular min e max
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
	// Carregar a imagem de texture
	textureSurface, err := sdl.LoadBMP(imagePath)
	if err != nil {
		fmt.Printf("Erro ao carregar imagem de texture: %v\n", err)
		return
	}
	defer textureSurface.Free()

	textureTexture, err := renderer.CreateTextureFromSurface(textureSurface)
	if err != nil {
		fmt.Printf("Erro ao criar textura de texture: %v\n", err)
		return
	}
	defer textureTexture.Destroy()

	// Desenhe a textura na posição e tamanho especificados
	renderer.Copy(textureTexture, nil, &sdl.Rect{X: x, Y: y, W: width, H: height})
}

// WrapText divide um texto longo em várias linhas com base na largura máxima especificada.
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

// textWidth calcula a largura de uma string de texto com base na fonte fornecida.
func textWidth(font *ttf.Font, text string) int {
	surface, err := font.RenderUTF8Blended(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return 0
	}
	defer surface.Free()

	return int(surface.W)
}
