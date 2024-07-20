package screens

import (
	"fmt"
	"handheldui/components"
	"handheldui/helpers"
	"handheldui/input"
	"handheldui/services"
	"handheldui/vars"

	"github.com/veandco/go-sdl2/sdl"
)

type GamesScreen struct {
	scrollOffset       int
	currentGameIndex   int
	games              []map[string]interface{}
	currentGameTexture *sdl.Texture
	renderer           *sdl.Renderer
	initialized        bool
	listComponent      *components.ListComponent
}

func NewGamesScreen(renderer *sdl.Renderer) (*GamesScreen, error) {
	listComponent := components.NewListComponent(renderer, "Games List", func(index int, item map[string]interface{}) string {
		return fmt.Sprintf("%d. %s", index+1, item["name"].(string))
	})

	g := &GamesScreen{
		renderer:      renderer,
		listComponent: listComponent,
	}

	return g, nil
}

func (g *GamesScreen) InitGames() {
	if g.initialized {
		return
	}

	games, err := services.FetchGames(vars.CurrentPlatform, vars.CurrentSystem)
	if err != nil {
		fmt.Printf("Erro ao buscar jogos: %v\n", err)
		return
	}
	g.games = games
	g.initialized = true
}

func (g *GamesScreen) HandleInput(event input.InputEvent) {
	if len(g.games) == 0 {
		return
	}

	switch event.KeyCode {
	case sdl.SCANCODE_DOWN:
		g.currentGameIndex = (g.currentGameIndex + 1) % len(g.games)
		g.scrollOffset = g.currentGameIndex
		g.LoadGameImage()
		g.listComponent.SetItems(g.games, g.currentGameIndex, g.scrollOffset)
	case sdl.SCANCODE_UP:
		g.currentGameIndex = (g.currentGameIndex - 1 + len(g.games)) % len(g.games)
		g.scrollOffset = g.currentGameIndex
		g.LoadGameImage()
		g.listComponent.SetItems(g.games, g.currentGameIndex, g.scrollOffset)
	case sdl.SCANCODE_B:
		g.initialized = false
		g.currentGameIndex = 0
		g.scrollOffset = g.currentGameIndex
		vars.CurrentScreen = "systems_screen"
	case sdl.SCANCODE_A:
		g.currentGameIndex = 0
		g.scrollOffset = g.currentGameIndex
		vars.CurrentGame = g.games[g.currentGameIndex]["key"].(string)
		vars.CurrentScreen = "overview_screen"
	}
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

	// Atualize o componente da lista com os dados atuais
	g.listComponent.SetItems(g.games, g.currentGameIndex, g.scrollOffset)
	g.listComponent.Draw()

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
