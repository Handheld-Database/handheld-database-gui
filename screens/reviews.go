package screens

import (
	"fmt"
	"handheldui/components"
	"handheldui/helpers/sdlutils"
	"handheldui/input"
	"handheldui/output"
	"handheldui/services"
	"handheldui/vars"

	"github.com/veandco/go-sdl2/sdl"
)

type ReviewsScreen struct {
	renderer      *sdl.Renderer
	testers       []map[string]interface{}
	initialized   bool
	listComponent *components.ListComponent
}

func NewReviewsScreen(renderer *sdl.Renderer) (*ReviewsScreen, error) {
	listComponent := components.NewListComponent(
		renderer,
		vars.Config.Screen.MaxListItens,
		vars.Config.Screen.MaxListItemWidth,
		func(index int, item map[string]interface{}) string {
			return fmt.Sprintf("%d. %s", index+1, item["name"].(string))
		})

	s := &ReviewsScreen{
		renderer:      renderer,
		listComponent: listComponent,
	}

	return s, nil
}

func (r *ReviewsScreen) InitReviews() {
	if r.initialized {
		return
	}

	testers, err := services.FetchTesters(vars.CurrentPlatform, vars.CurrentSystem, vars.CurrentGame)
	if err != nil {
		output.Errorf("Error fetching games: %v\n", err)
		return
	}

	testersList := make([]map[string]interface{}, len(testers))
	for i, tester := range testers {
		testersList[i] = map[string]interface{}{"name": tester}
	}
	r.testers = testersList
	r.listComponent.SetItems(r.testers)
	r.initialized = true
}

func (r *ReviewsScreen) HandleInput(event input.InputEvent) {
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
		r.showReview()
	case "B":
		vars.CurrentScreen = "games_screen"
		vars.CurrentGame = ""
		r.initialized = false
	}
}

func (r *ReviewsScreen) Draw() {
	r.InitReviews()

	r.renderer.SetDrawColor(255, 255, 255, 255)
	r.renderer.Clear()

	sdlutils.RenderTextureCartesian(r.renderer, "assets/textures/bg.bmp", "Q2", "Q4")

	// Draw the current title
	sdlutils.DrawText(r.renderer, "Reviewers List", sdl.Point{X: 25, Y: 25}, vars.Colors.WHITE, vars.HeaderFont)

	// Draw the list component
	r.listComponent.Draw(vars.Colors.SECONDARY, vars.Colors.WHITE)

	sdlutils.RenderTextureCartesian(r.renderer, "assets/textures/$aspect_ratio/ui_controls.bmp", "Q3", "Q4")

	r.renderer.Present()
}

func (r *ReviewsScreen) showReview() {
	if len(r.listComponent.GetItems()) == 0 {
		return
	}

	selectedTester := r.listComponent.GetItems()[r.listComponent.GetSelectedIndex()]
	vars.CurrentTester = selectedTester["name"].(string)
	vars.CurrentScreen = "overview_screen"
}
