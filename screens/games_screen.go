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
	currentImage  string
	games         []map[string]interface{}
	renderer      *sdl.Renderer
	initialized   bool
	listComponent *components.ListComponent
	textureMutex  sync.Mutex
}

func NewGamesScreen(renderer *sdl.Renderer) (*GamesScreen, error) {
	listComponent := components.NewListComponent(renderer, "Games List", 20, func(index int, item map[string]interface{}) string {
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
	g.listComponent.SetItems(g.games)
	g.initialized = true
}

func (g *GamesScreen) HandleInput(event input.InputEvent) {
	if len(g.games) == 0 {
		return
	}

	switch event.KeyCode {
	case sdl.SCANCODE_DOWN:
		g.listComponent.ScrollDown()
	case sdl.SCANCODE_UP:
		g.listComponent.ScrollUp()
	case sdl.SCANCODE_B:
		g.initialized = false
		g.listComponent.SetItems([]map[string]interface{}{})
		vars.CurrentScreen = "systems_screen"
	case sdl.SCANCODE_A:
		selectedGame := g.games[g.listComponent.GetSelectedIndex()]
		vars.CurrentGame = selectedGame["key"].(string)
		vars.CurrentScreen = "overview_screen"
	}
}

func (g *GamesScreen) LoadGameImage() {
	selectedIndex := g.listComponent.GetSelectedIndex()
	if selectedIndex < len(g.games) {
		gameName := g.games[selectedIndex]["key"].(string)
		imagePath := helpers.FetchGameImage(gameName)
		if imagePath != "" {
			g.textureMutex.Lock()
			g.currentImage = imagePath
			g.textureMutex.Unlock()

			// Debug message to confirm the texture is loaded
			fmt.Printf("Imagem carregada para o jogo: %s\n", gameName)
		} else {
			g.currentImage = ""
		}
	}
}

func (g *GamesScreen) Draw() {
	g.InitGames()

	go g.LoadGameImage()

	g.renderer.SetDrawColor(255, 255, 255, 255)
	g.renderer.Clear()

	helpers.RenderTexture(g.renderer, "assets/textures/bg.bmp", "Q2", "Q4")

	helpers.DrawText(g.renderer, "Systems List", sdl.Point{X: 25, Y: 25}, vars.Colors.PRIMARY, vars.HeaderFont)

	g.listComponent.Draw(vars.Colors.WHITE, vars.Colors.SECONDARY)

	g.textureMutex.Lock()
	defer g.textureMutex.Unlock()
	if g.currentImage != "" {
		helpers.RenderTextureAdjusted(g.renderer, g.currentImage, vars.ScreenWidth-340-84, 78, 340, 340)
	} else {
		helpers.RenderTextureAdjusted(g.renderer, "assets/textures/not_found.bmp", vars.ScreenWidth-340-84, 78, 340, 340)
		fmt.Println("Nenhuma textura disponível para desenhar.")
	}

	rankAssets := map[string]string{
		"PLATINUM": "assets/textures/ui_game_display_rank_platinum_1280_720.bmp",
		"GOLD":     "assets/textures/ui_game_display_rank_gold_1280_720.bmp",
		"SILVER":   "assets/textures/ui_game_display_rank_silver_1280_720.bmp",
		"BRONZE":   "assets/textures/ui_game_display_rank_bronze_1280_720.bmp",
		"FAULTY":   "assets/textures/ui_game_display_rank_faulty_1280_720.bmp",
	}
	selectedIndex := g.listComponent.GetSelectedIndex()
	gameRank := g.games[selectedIndex]["rank"].(string)

	helpers.RenderTexture(g.renderer, "assets/textures/ui_game_display_1280_720.bmp", "Q1", "Q4")
	helpers.RenderTexture(g.renderer, rankAssets[gameRank], "Q1", "Q1")
	helpers.RenderTexture(g.renderer, "assets/textures/ui_game_display_details_1280_720.bmp", "Q4", "Q4")
	helpers.RenderTexture(g.renderer, "assets/textures/ui_controls_1280_720.bmp", "Q3", "Q4")

	var postitionRank = sdl.Point{X: vars.ScreenWidth - 420, Y: vars.ScreenHeight - 270}
	helpers.DrawText(g.renderer, gameRank, postitionRank, vars.Colors.BLACK, vars.BodyBigFont)

	g.renderer.Present()
}

func (g *GamesScreen) ShowGameInfo() {
	if len(g.games) > 0 {
		selectedIndex := g.listComponent.GetSelectedIndex()
		gameName := g.games[selectedIndex]["name"].(string)
		fmt.Printf("Selecionado jogo: %s\n", gameName)
	}
}
