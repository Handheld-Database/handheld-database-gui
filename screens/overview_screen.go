package screens

import (
	"fmt"
	"handheldui/helpers"
	"handheldui/input"
	"handheldui/services"
	"handheldui/vars"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type OverviewScreen struct {
	detectedPlatform string
	systems          []map[string]interface{}
	renderer         *sdl.Renderer
	currentView      string
	gameImage        *sdl.Texture
}

func NewOverviewScreen(renderer *sdl.Renderer) (*OverviewScreen, error) {
	// Inicializar SDL_ttf
	if err := ttf.Init(); err != nil {
		return nil, fmt.Errorf("erro ao inicializar SDL_ttf: %w", err)
	}

	s := &OverviewScreen{
		renderer:    renderer,
		currentView: "Overview",
	}

	return s, nil
}

func (s *OverviewScreen) HandleInput(event input.InputEvent) {
	switch event.KeyCode {
	case "B":
		vars.CurrentScreen = "games_screen"
		vars.CurrentGame = ""
	}
}

func (s *OverviewScreen) Draw() {
	systemsData, err := services.FetchPlatform(s.detectedPlatform)
	if err != nil {
		fmt.Println("Error fetching platform data:", err)
		return
	}

	if systemsData["systems"] != nil {
		systems := systemsData["systems"].([]interface{})
		s.systems = make([]map[string]interface{}, len(systems))
		for i, system := range systems {
			s.systems[i] = system.(map[string]interface{})
		}
	}

	s.renderer.SetDrawColor(255, 255, 255, 255)
	s.renderer.Clear()

	// Desenhar o título atual
	titleColor := vars.Colors.WHITE
	textSurface, err := helpers.RenderText(s.currentView, titleColor, vars.BodyFont)
	if err != nil {
		fmt.Printf("Erro ao renderizar texto: %v\n", err)
		return
	}
	defer textSurface.Free()

	titleTexture, err := s.renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		fmt.Printf("Erro ao criar textura: %v\n", err)
		return
	}
	defer titleTexture.Destroy()

	s.renderer.Copy(titleTexture, nil, &sdl.Rect{X: 500, Y: 10, W: int32(textSurface.W), H: int32(textSurface.H)})

	// Desenhar a imagem do jogo no canto superior direito
	if s.gameImage != nil {
		s.renderer.Copy(s.gameImage, nil, &sdl.Rect{X: 500, Y: 60, W: 100, H: 100})
	}

	// Desenhar o conteúdo de Overview ou Details
	if s.currentView == "Overview" {
		s.drawOverview()
	} else {
		s.drawDetails()
	}

	s.renderer.Present()
}

func (s *OverviewScreen) drawOverview() {
	// Lógica para desenhar a visão Overview
	overviewText, err := services.FetchGameOverview(vars.CurrentGame)
	if err != nil {
		fmt.Println("Error fetching overview:", err)
		return
	}

	textSurface, err := helpers.RenderText(overviewText, sdl.Color{R: 0, G: 0, B: 0, A: 255}, vars.BodyFont)
	if err != nil {
		fmt.Printf("Erro ao renderizar texto: %v\n", err)
		return
	}
	defer textSurface.Free()

	overviewTexture, err := s.renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		fmt.Printf("Erro ao criar textura: %v\n", err)
		return
	}
	defer overviewTexture.Destroy()

	s.renderer.Copy(overviewTexture, nil, &sdl.Rect{X: 40, Y: 200, W: int32(textSurface.W), H: int32(textSurface.H)})
}

func (s *OverviewScreen) drawDetails() {
	// Lógica para desenhar a visão Details
	detailsText, err := services.FetchGameMarkdown(s.detectedPlatform, vars.CurrentSystem, vars.CurrentGame)
	if err != nil {
		fmt.Println("Error fetching details:", err)
		return
	}

	textSurface, err := helpers.RenderText(detailsText, sdl.Color{R: 0, G: 0, B: 0, A: 255}, vars.HeaderFont)
	if err != nil {
		fmt.Printf("Erro ao renderizar texto: %v\n", err)
		return
	}
	defer textSurface.Free()

	detailsTexture, err := s.renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		fmt.Printf("Erro ao criar textura: %v\n", err)
		return
	}
	defer detailsTexture.Destroy()

	s.renderer.Copy(detailsTexture, nil, &sdl.Rect{X: 40, Y: 200, W: int32(textSurface.W), H: int32(textSurface.H)})
}

func (s *OverviewScreen) LoadGameImage(imagePath string) error {
	imageSurface, err := sdl.LoadBMP(imagePath)
	if err != nil {
		return fmt.Errorf("erro ao carregar imagem: %w", err)
	}
	defer imageSurface.Free()

	s.gameImage, err = s.renderer.CreateTextureFromSurface(imageSurface)
	if err != nil {
		return fmt.Errorf("erro ao criar textura de imagem: %w", err)
	}

	return nil
}
