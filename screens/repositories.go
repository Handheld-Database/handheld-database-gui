package screens

import (
	"handheldui/components"
	"handheldui/helpers/sdlutils"
	"handheldui/input"
	"handheldui/vars"

	"github.com/veandco/go-sdl2/sdl"
)

type RepositoriesScreen struct {
	initialized   bool
	renderer      *sdl.Renderer
	listComponent *components.ListComponent
}

func NewRepositoriesScreen(renderer *sdl.Renderer) (*RepositoriesScreen, error) {
	listComponent := components.NewListComponent(renderer, 10, func(index int, item map[string]interface{}) string {
		return item["name"].(string)
	})

	return &RepositoriesScreen{
		renderer:      renderer,
		listComponent: listComponent,
	}, nil
}

func (r *RepositoriesScreen) InitRepositories() {
	if r.initialized {
		return
	}

	repositories := vars.Config.Repositories

	var items []map[string]interface{}

	for innerKey, repo := range repositories {
		items = append(items, map[string]interface{}{
			"name":  repo.Name,
			"value": innerKey,
		})
	}

	r.listComponent.SetItems(items)

	r.initialized = true
}

func (r *RepositoriesScreen) HandleInput(event input.InputEvent) {
	if len(r.listComponent.GetItems()) == 0 {
		return
	}

	switch event.KeyCode {
	case "DOWN":
		r.listComponent.ScrollDown()
	case "UP":
		r.listComponent.ScrollUp()
	case "L1":
		r.listComponent.PageUp()
	case "R1":
		r.listComponent.PageDown()
	case "A":
		selectedItem := r.listComponent.GetItems()[r.listComponent.GetSelectedIndex()]
		vars.CurrentRepo = selectedItem["value"].(string)
		vars.CurrentScreen = "files_screen"
	case "B":
		vars.CurrentScreen = "home_screen"
	}
}

func (r *RepositoriesScreen) Draw() {
	r.InitRepositories()

	r.renderer.SetDrawColor(255, 255, 255, 255)
	r.renderer.Clear()

	sdlutils.RenderTexture(r.renderer, "assets/textures/bg.bmp", "Q2", "Q4")

	// Draw the current title
	sdlutils.DrawText(r.renderer, "Repositories List", sdl.Point{X: 25, Y: 25}, vars.Colors.WHITE, vars.HeaderFont)

	// Draw the list component
	r.listComponent.Draw(vars.Colors.SECONDARY, vars.Colors.WHITE)

	sdlutils.RenderTexture(r.renderer, "assets/textures/ui_controls_1280_720.bmp", "Q3", "Q4")

	r.renderer.Present()
}
