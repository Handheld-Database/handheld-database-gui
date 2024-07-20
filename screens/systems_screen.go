package screens

import (
	"fmt"
	"handheldui/helpers"
	"handheldui/input"
	"handheldui/services"
	"handheldui/vars"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type SystemsScreen struct {
	detectedPlatform string
	selectedSystem   int
	systems          []map[string]interface{}
	renderer         *sdl.Renderer
	initialized      bool
}

func NewSystemsScreen(platform string, renderer *sdl.Renderer) (*SystemsScreen, error) {
	s := &SystemsScreen{
		detectedPlatform: platform,
		renderer:         renderer,
	}

	return s, nil
}

func (s *SystemsScreen) InitSystems() {
	if s.initialized {
		return
	}

	systemsData, err := services.FetchPlatform(s.detectedPlatform)
	if err != nil {
		fmt.Println("Error fetching platform data:", err)
		return
	}

	systems := systemsData["systems"].([]interface{})
	s.systems = make([]map[string]interface{}, len(systems))
	for i, system := range systems {
		s.systems[i] = system.(map[string]interface{})
	}
	s.initialized = true
}

func (s *SystemsScreen) HandleInput(event input.InputEvent) {
	go func() {
		if len(s.systems) == 0 {
			return
		}

		systemsCount := len(s.systems)

		switch event.KeyCode {
		case sdl.SCANCODE_DOWN:
			s.selectedSystem = (s.selectedSystem + 1) % systemsCount
		case sdl.SCANCODE_UP:
			s.selectedSystem = (s.selectedSystem - 1 + systemsCount) % systemsCount
		case sdl.SCANCODE_A:
			s.showGames()
		case sdl.SCANCODE_B:
			os.Exit(0)
		}
	}()
}

func (s *SystemsScreen) Draw() {
	s.InitSystems()

	s.renderer.SetDrawColor(255, 255, 255, 255)
	s.renderer.Clear()

	helpers.RenderTexture(s.renderer, "assets/textures/bg.bmp")

	// Desenhe o título
	titleColor := sdl.Color{R: 0, G: 0, B: 255, A: 255} // Azul para o título
	titleSurface, err := helpers.RenderText("Systems List", titleColor, vars.HeaderFont)
	if err != nil {
		fmt.Printf("Erro ao renderizar texto do título: %v\n", err)
		return
	}
	defer titleSurface.Free()

	titleTexture, err := s.renderer.CreateTextureFromSurface(titleSurface)
	if err != nil {
		fmt.Printf("Erro ao criar textura do título: %v\n", err)
		return
	}
	defer titleTexture.Destroy()

	s.renderer.Copy(titleTexture, nil, &sdl.Rect{X: 40, Y: 20, W: int32(titleSurface.W), H: int32(titleSurface.H)})

	// Desenhe os sistemas
	for index, system := range s.systems {
		name := system["name"].(string)
		color := sdl.Color{R: 0, G: 0, B: 0, A: 255}
		if index == s.selectedSystem {
			color = sdl.Color{R: 128, G: 0, B: 128, A: 255}
		}
		textSurface, err := helpers.RenderText(fmt.Sprintf("%d. %s", index+1, name), color, vars.BodyFont)
		if err != nil {
			fmt.Printf("Erro ao renderizar texto: %v\n", err)
			return
		}
		defer textSurface.Free()

		texture, err := s.renderer.CreateTextureFromSurface(textSurface)
		if err != nil {
			fmt.Printf("Erro ao criar textura: %v\n", err)
			return
		}
		defer texture.Destroy()

		s.renderer.Copy(texture, nil, &sdl.Rect{X: 40, Y: 60 + 30*int32(index), W: int32(textSurface.W), H: int32(textSurface.H)})
	}

	s.renderer.Present()
}

func (s *SystemsScreen) showGames() {
	if len(s.systems) == 0 {
		return
	}

	vars.CurrentScreen = "games_screen"
	selectedSystemKey := s.systems[s.selectedSystem]["key"].(string)
	fmt.Printf("Selecionado sistema: %s\n", selectedSystemKey)
	vars.CurrentSystem = selectedSystemKey
}
