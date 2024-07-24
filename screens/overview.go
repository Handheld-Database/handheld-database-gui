package screens

import (
	"handheldui/components"
	"handheldui/helpers"
	"handheldui/input"
	"handheldui/services"
	"handheldui/vars"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

type OverviewScreen struct {
	renderer      *sdl.Renderer
	textComponent *components.TextComponent
	textContent   string
	initialized   bool
}

func NewOverviewScreen(renderer *sdl.Renderer) (*OverviewScreen, error) {
	return &OverviewScreen{
		renderer: renderer,
	}, nil
}

func (o *OverviewScreen) InitOverview() {
	if o.initialized {
		return
	}

	overview, err := services.FetchGameOverview(vars.CurrentGame)
	if err != nil {
		overview = "Help us to find an overview!"
	}

	review, err := services.FetchGameMarkdown(vars.CurrentPlatform, vars.CurrentSystem, vars.CurrentGame)
	if err != nil {
		overview = "Oops, game description not found!"
	}

	o.textContent = helpers.MarkdownToPlaintext(strings.ReplaceAll(review, "%game_overview%", overview))
	o.textComponent = components.NewTextComponent(o.renderer, o.textContent, vars.LongTextFont, 18, 1200)

	o.initialized = true
}

func (o *OverviewScreen) HandleInput(event input.InputEvent) {
	switch event.KeyCode {
	case "DOWN":
		o.textComponent.ScrollDown()
	case "UP":
		o.textComponent.ScrollUp()
	case "B":
		vars.CurrentGame = ""
		vars.CurrentScreen = "games_screen"
		o.initialized = false
	}
}

func (o *OverviewScreen) Draw() {
	o.InitOverview()

	o.renderer.SetDrawColor(0, 0, 0, 255) // Background color
	o.renderer.Clear()

	helpers.RenderTexture(o.renderer, "assets/textures/bg.bmp", "Q2", "Q4")

	helpers.RenderTexture(o.renderer, "assets/textures/bg_overlay.bmp", "Q2", "Q4")

	// Draw the title
	helpers.DrawText(o.renderer, "Overview", sdl.Point{X: 25, Y: 25}, vars.Colors.PRIMARY, vars.HeaderFont)

	// Draw the text component with scrolling
	o.textComponent.Draw(vars.Colors.WHITE)

	helpers.RenderTexture(o.renderer, "assets/textures/ui_controls_1280_720.bmp", "Q3", "Q4")

	o.renderer.Present()
}