package screens

import (
	"fmt"
	"handheldui/components"
	"handheldui/helpers"
	"handheldui/input"
	"handheldui/services"
	"handheldui/vars"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type GamesScreen struct {
	scrollOffset     int
	currentGameIndex int
	currentImage     string
	games            []map[string]interface{}
	renderer         *sdl.Renderer
	initialized      bool
	listComponent    *components.ListComponent
	textureMutex     sync.Mutex
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
		go g.LoadGameImage() // Start loading the image asynchronously
		g.listComponent.SetItems(g.games, g.currentGameIndex, g.scrollOffset)
	case sdl.SCANCODE_UP:
		g.currentGameIndex = (g.currentGameIndex - 1 + len(g.games)) % len(g.games)
		g.scrollOffset = g.currentGameIndex
		go g.LoadGameImage() // Start loading the image asynchronously
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
			g.textureMutex.Lock()
			g.currentImage = imagePath
			g.textureMutex.Unlock()

			// Debug message to confirm the texture is loaded
			fmt.Printf("Imagem carregada para o jogo: %s\n", gameName)
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

	g.textureMutex.Lock()
	defer g.textureMutex.Unlock()
	if g.currentImage != "" {
		// Debug message to confirm the draw call
		fmt.Println("Desenhando textura atual...")
		helpers.RenderTextureAdjusted(g.renderer, g.currentImage, vars.ScreenWidth-340-84, 78, 340, 340)
	} else {
		// Debug message if no texture is available
		fmt.Println("Nenhuma textura disponível para desenhar.")
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
