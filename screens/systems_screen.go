package screens

import (
	"fmt"
	"handheldui/components"
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
	listComponent    *components.ListComponent
}

func NewSystemsScreen(renderer *sdl.Renderer) (*SystemsScreen, error) {
	listComponent := components.NewListComponent(renderer, "Systems List", func(index int, item map[string]interface{}) string {
		return fmt.Sprintf("%d. %s", index+1, item["name"].(string))
	})

	s := &SystemsScreen{
		detectedPlatform: vars.CurrentPlatform,
		renderer:         renderer,
		listComponent:    listComponent,
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

	// Atualize o ListComponent
	s.listComponent.SetItems(s.systems, s.selectedSystem, s.selectedSystem)
}

func (s *SystemsScreen) Draw() {
	s.InitSystems()

	s.renderer.SetDrawColor(255, 255, 255, 255)
	s.renderer.Clear()

	helpers.RenderTexture(s.renderer, "assets/textures/bg.bmp")

	// Atualize o componente da lista com os dados atuais
	s.listComponent.SetItems(s.systems, s.selectedSystem, s.selectedSystem)
	s.listComponent.Draw()

	s.renderer.Present()
}

func (s *SystemsScreen) showGames() {
	if len(s.systems) == 0 {
		return
	}

	selectedSystemKey := s.systems[s.selectedSystem]["key"].(string)
	fmt.Printf("Selecionado sistema: %s\n", selectedSystemKey)
	vars.CurrentSystem = selectedSystemKey
	vars.CurrentScreen = "games_screen"
}
