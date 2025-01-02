package screens

import (
	"handheldui/components"
	"handheldui/helpers/sdlutils"
	"handheldui/input"
	"handheldui/vars"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type HomeScreen struct {
	initialized   bool
	renderer      *sdl.Renderer
	listComponent *components.ListComponent
}

func NewHomeScreen(renderer *sdl.Renderer) (*HomeScreen, error) {
	listComponent := components.NewListComponent(renderer, 15, func(index int, item map[string]interface{}) string {
		return item["label"].(string)
	})

	return &HomeScreen{
		renderer:      renderer,
		listComponent: listComponent,
	}, nil
}

func (h *HomeScreen) InitHome() {
	if h.initialized {
		return
	}

	buttons := []map[string]interface{}{
		{"label": "Reviews", "action": func() { vars.CurrentScreen = "systems_screen" }},
		{"label": "Repositories", "action": func() { vars.CurrentScreen = "repositories_screen" }},
	}

	h.listComponent.SetItems(buttons)

	h.initialized = true
}

func (h *HomeScreen) HandleInput(event input.InputEvent) {
	if len(h.listComponent.GetItems()) == 0 {
		return
	}

	switch event.KeyCode {
	case "DOWN":
		h.listComponent.ScrollDown()
	case "UP":
		h.listComponent.ScrollUp()
	case "A":
		selectedItem := h.listComponent.GetItems()[h.listComponent.GetSelectedIndex()]
		if action, ok := selectedItem["action"].(func()); ok {
			action()
		}
	case "B":
		os.Exit(0)
	}
}

func (h *HomeScreen) Draw() {
	h.InitHome()

	h.renderer.SetDrawColor(255, 255, 255, 255)
	h.renderer.Clear()

	sdlutils.RenderTexture(h.renderer, "assets/textures/bg.bmp", "Q2", "Q4")

	// Draw the current title
	sdlutils.DrawText(h.renderer, "Home", sdl.Point{X: 25, Y: 25}, vars.Colors.PRIMARY, vars.HeaderFont)

	// Draw the list component
	h.listComponent.Draw(vars.Colors.WHITE, vars.Colors.SECONDARY)

	sdlutils.RenderTexture(h.renderer, "assets/textures/ui_controls_1280_720.bmp", "Q3", "Q4")

	h.renderer.Present()
}
