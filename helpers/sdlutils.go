package helpers

import (
	"fmt"
	"handheldui/vars"

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
	if err := mix.Init(mix.INIT_MP3 | mix.INIT_FLAC | mix.INIT_OGG); err != nil {
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

func RenderTexture(renderer *sdl.Renderer, imagePath string) {
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

	// Desenhe o texture em cima de tudo
	renderer.Copy(textureTexture, nil, &sdl.Rect{X: 0, Y: 0, W: vars.ScreenWidth, H: vars.ScreenHeight})
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
