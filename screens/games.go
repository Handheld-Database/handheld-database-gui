package screens

import (
	"fmt"
	"handheldui/components"
	"handheldui/helpers/geometry"
	"handheldui/helpers/image"
	"handheldui/helpers/sdlutils"
	"handheldui/input"
	"handheldui/output"
	"handheldui/services"
	"handheldui/vars"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type GamesScreen struct {
	currentImageCover string
	currentImageIcon  string
	games             []map[string]interface{}
	renderer          *sdl.Renderer
	initialized       bool
	listComponent     *components.ListComponent
	textureCoverMutex sync.Mutex
	textureIconMutex  sync.Mutex
}

func NewGamesScreen(renderer *sdl.Renderer) (*GamesScreen, error) {
	listComponent := components.NewListComponent(renderer, 10, func(index int, item map[string]interface{}) string {
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
		output.Errorf("Error fetching games: %v\n", err)
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
	case "DOWN":
		g.listComponent.ScrollDown()
	case "UP":
		g.listComponent.ScrollUp()
	case "L1":
		g.listComponent.PageUp()
	case "R1":
		g.listComponent.PageDown()
	case "B":
		g.initialized = false
		g.listComponent.SetItems([]map[string]interface{}{})
		vars.CurrentScreen = "systems_screen"
	case "A":
		selectedGame := g.games[g.listComponent.GetSelectedIndex()]
		vars.CurrentGame = selectedGame["key"].(string)
		vars.CurrentScreen = "reviews_screen"
	}
}

func (g *GamesScreen) LoadGameImage() {
	selectedIndex := g.listComponent.GetSelectedIndex()
	if selectedIndex < len(g.games) {
		gameName := g.games[selectedIndex]["key"].(string)
		imageCoverPath := image.FetchGameImage(gameName, "cover")
		imageIconPath := image.FetchGameImage(gameName, "icon")

		if imageCoverPath != "" {
			g.textureCoverMutex.Lock()
			g.currentImageCover = imageCoverPath
			g.textureCoverMutex.Unlock()

			// Debug message to confirm the texture is loaded
			output.Printf("Cover loaded for game: %s\n", gameName)
		} else {
			g.currentImageCover = ""
		}

		if imageIconPath != "" {
			g.textureCoverMutex.Lock()
			g.currentImageIcon = imageIconPath
			g.textureCoverMutex.Unlock()

			// Debug message to confirm the texture is loaded
			output.Printf("Icon loaded for game: %s\n", gameName)
		} else {
			g.currentImageIcon = ""
		}
	}
}

func (g *GamesScreen) Draw() {
	g.InitGames()

	go g.LoadGameImage()

	g.renderer.SetDrawColor(255, 255, 255, 255)
	g.renderer.Clear()

	if g.currentImageCover != "" {
		sdlutils.RenderTextureCover(g.renderer, g.currentImageCover)
	} else {
		sdlutils.RenderTextureCover(g.renderer, "assets/textures/bg.bmp")
	}

	sdlutils.RenderTextureCover(g.renderer, "assets/textures/bg_overlay.bmp")

	sdlutils.DrawText(g.renderer, "Games List", sdl.Point{X: 25, Y: 25}, vars.Colors.WHITE, vars.HeaderFont)

	g.listComponent.Draw(vars.Colors.SECONDARY, vars.Colors.WHITE)

	g.textureCoverMutex.Lock()
	defer g.textureCoverMutex.Unlock()

	if vars.Config.Screen.Width >= 1280 {

		sdlutils.RenderTexture(g.renderer, "assets/textures/ui_game_display_1280_720.bmp", "Q1", "Q4")

		element := geometry.NewElement(280, 280, 0, 78, "top-right")

		if g.currentImageCover != "" {
			sdlutils.RenderTextureAdjusted(g.renderer, g.currentImageCover, element.GetPosition())
		} else {
			sdlutils.RenderTextureAdjusted(g.renderer, "assets/textures/not_found.bmp", element.GetPosition())
			output.Printf("No texture available to draw.")
		}

		g.textureIconMutex.Lock()
		defer g.textureIconMutex.Unlock()

		rankAssets := map[string]string{
			"PLATINUM": "assets/textures/ui_game_display_rank_platinum_1280_720.bmp",
			"GOLD":     "assets/textures/ui_game_display_rank_gold_1280_720.bmp",
			"SILVER":   "assets/textures/ui_game_display_rank_silver_1280_720.bmp",
			"BRONZE":   "assets/textures/ui_game_display_rank_bronze_1280_720.bmp",
			"FAULTY":   "assets/textures/ui_game_display_rank_faulty_1280_720.bmp",
		}

		selectedIndex := g.listComponent.GetSelectedIndex()
		gameRank := g.games[selectedIndex]["rank"].(string)

		sdlutils.RenderTexture(g.renderer, "assets/textures/ui_game_display_1280_720_overlay.bmp", "Q1", "Q4")
		sdlutils.RenderTexture(g.renderer, rankAssets[gameRank], "Q1", "Q4")
		sdlutils.RenderTexture(g.renderer, "assets/textures/ui_controls_1280_720.bmp", "Q3", "Q4")
	}

	g.renderer.Present()
}

func (g *GamesScreen) ShowGameInfo() {
	if len(g.games) > 0 {
		selectedIndex := g.listComponent.GetSelectedIndex()
		gameName := g.games[selectedIndex]["name"].(string)
		output.Printf("Selected game: %s\n", gameName)
	}
}
