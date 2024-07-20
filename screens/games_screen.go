package screens

import (
	"fmt"
	"handheldui/helpers"
	"handheldui/input"
	"handheldui/services"
	"handheldui/vars"

	"github.com/veandco/go-sdl2/sdl"
)

type GamesScreen struct {
	detectedPlatform   string
	scrollOffset       int
	currentGameIndex   int
	games              []map[string]interface{}
	currentGameTexture *sdl.Texture
	renderer           *sdl.Renderer
	initialized        bool
}

func NewGamesScreen(platform string, renderer *sdl.Renderer) (*GamesScreen, error) {
	g := &GamesScreen{
		detectedPlatform: platform,
		renderer:         renderer,
	}

	return g, nil
}

func (g *GamesScreen) InitGames() {
	if g.initialized {
		return
	}

	games, err := services.FetchGames(g.detectedPlatform, vars.CurrentSystem)
	if err != nil {
		fmt.Printf("Erro ao buscar jogos: %v\n", err)
		return
	}
	g.games = games
	g.initialized = true
}

func (g *GamesScreen) HandleInput(event input.InputEvent) {
	go func() {
		if len(g.games) == 0 {
			return
		}

		switch event.KeyCode {
		case sdl.SCANCODE_DOWN:
			g.currentGameIndex = (g.currentGameIndex + 1) % len(g.games)
			g.scrollOffset = g.currentGameIndex
			g.LoadGameImage()
		case sdl.SCANCODE_UP:
			g.currentGameIndex = (g.currentGameIndex - 1 + len(g.games)) % len(g.games)
			g.scrollOffset = g.currentGameIndex
			g.LoadGameImage()
		case sdl.SCANCODE_B:
			g.currentGameIndex = 0
			g.scrollOffset = g.currentGameIndex
			vars.CurrentScreen = "systems_screen"
		case sdl.SCANCODE_A:
			g.currentGameIndex = 0
			g.scrollOffset = g.currentGameIndex
			vars.CurrentGame = g.games[g.currentGameIndex]["key"].(string)
			vars.CurrentScreen = "overview_screen"
		}
	}()
}

func (g *GamesScreen) LoadGameImage() {
	if g.currentGameIndex < len(g.games) {
		gameName := g.games[g.currentGameIndex]["key"].(string)
		imagePath := helpers.FetchGameImage(gameName)
		if imagePath != "" {
			texture, err := helpers.LoadTexture(g.renderer, imagePath)
			if err != nil {
				fmt.Printf("Erro ao carregar textura: %v\n", err)
				return
			}
			g.currentGameTexture = texture
		}
	}
}

func (g *GamesScreen) Draw() {
	g.InitGames()

	g.renderer.SetDrawColor(255, 255, 255, 255)
	g.renderer.Clear()

	helpers.RenderTexture(g.renderer, "assets/textures/bg.bmp")

	startIndex := g.scrollOffset
	endIndex := startIndex + int(vars.ScreenHeight/30)
	if endIndex > len(g.games) {
		endIndex = len(g.games)
	}
	visibleGames := g.games[startIndex:endIndex]

	for index, game := range visibleGames {
		color := sdl.Color{R: 255, G: 255, B: 255, A: 255}
		if index == g.currentGameIndex-startIndex {
			color = sdl.Color{R: 128, G: 0, B: 128, A: 255}
		}
		gameName := game["name"].(string)
		textSurface, err := helpers.RenderText(fmt.Sprintf("%s \u2605", gameName), color, vars.BodyFont)
		if err != nil {
			fmt.Printf("Erro ao renderizar texto: %v\n", err)
			return
		}
		defer textSurface.Free()

		texture, err := g.renderer.CreateTextureFromSurface(textSurface)
		if err != nil {
			fmt.Printf("Erro ao criar textura: %v\n", err)
			return
		}
		defer texture.Destroy()

		g.renderer.Copy(texture, nil, &sdl.Rect{X: 40, Y: 60 + 30*int32(index), W: textSurface.W, H: textSurface.H})
	}

	if g.currentGameTexture != nil {
		g.renderer.Copy(g.currentGameTexture, nil, &sdl.Rect{X: int32(vars.ScreenWidth) - 430, Y: 72, W: 350, H: 350})
	}

	helpers.RenderTexture(g.renderer, "assets/textures/ui_1280_720.bmp")

	g.renderer.Present()
}

func (g *GamesScreen) ShowGameInfo() {
	if len(g.games) > 0 {
		selectedIndex := g.scrollOffset
		gameName := g.games[selectedIndex]["name"].(string)
		fmt.Printf("Selecionado jogo: %s\n", gameName)
	}
}
