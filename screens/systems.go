package screens

import (
	"fmt"
	"handheldui/components"
	"handheldui/helpers/sdlutils"
	"handheldui/input"
	"handheldui/output"
	"handheldui/services"
	"handheldui/vars"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type SystemsScreen struct {
	detectedPlatform string
	renderer         *sdl.Renderer
	initialized      bool
	listComponent    *components.ListComponent
}

func NewSystemsScreen(renderer *sdl.Renderer) (*SystemsScreen, error) {
	listComponent := components.NewListComponent(renderer, 19, func(index int, item map[string]interface{}) string {
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
		output.Printf("Error fetching platform data:", err)
		return
	}

	systems := systemsData["systems"].([]interface{})
	systemsList := make([]map[string]interface{}, len(systems))
	for i, system := range systems {
		systemsList[i] = system.(map[string]interface{})
	}
	s.listComponent.SetItems(systemsList)
	s.initialized = true
}

func (s *SystemsScreen) HandleInput(event input.InputEvent) {
	if len(s.listComponent.GetItems()) == 0 {
		return
	}

	switch event.KeyCode {
	case "DOWN":
		s.listComponent.ScrollDown()
	case "UP":
		s.listComponent.ScrollUp()
	case "A":
		s.showGames()
	case "B":
		os.Exit(0)
	}
}

func (s *SystemsScreen) Draw() {
	s.InitSystems()

	s.renderer.SetDrawColor(255, 255, 255, 255)
	s.renderer.Clear()

	sdlutils.RenderTexture(s.renderer, "assets/textures/bg.bmp", "Q2", "Q4")

	// Draw the current title
	sdlutils.DrawText(s.renderer, "Systems List", sdl.Point{X: 25, Y: 25}, vars.Colors.PRIMARY, vars.HeaderFont)

	// Draw the list component
	s.listComponent.Draw(vars.Colors.WHITE, vars.Colors.SECONDARY)

	sdlutils.RenderTexture(s.renderer, "assets/textures/ui_controls_1280_720.bmp", "Q3", "Q4")

	s.renderer.Present()
}

func (s *SystemsScreen) showGames() {
	if len(s.listComponent.GetItems()) == 0 {
		return
	}

	selectedSystem := s.listComponent.GetItems()[s.listComponent.GetSelectedIndex()]
	selectedSystemKey := selectedSystem["key"].(string)
	output.Printf("Selected system: %s\n", selectedSystemKey)
	vars.CurrentSystem = selectedSystemKey
	vars.CurrentScreen = "games_screen"
}
