package screens

import (
	"handheldui/components"
	"handheldui/helpers/markdown"
	"handheldui/helpers/sdlutils"
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

	review, err := services.FetchGameMarkdown(vars.CurrentPlatform, vars.CurrentSystem, vars.CurrentGame, vars.CurrentTester)
	if err != nil {
		review = "Oops, game description not found!"
	}

	plainReview := markdown.MarkdownToPlaintext(review)
	plainOverview := markdown.MarkdownToPlaintext(overview)

	o.textContent = strings.ReplaceAll(plainReview, "%game_overview%", plainOverview)
	o.textComponent = components.NewTextComponent(o.renderer, o.textContent, vars.LongTextFont, vars.Config.Screen.MaxLines, int(vars.Config.Screen.Width)-20)

	o.initialized = true
}

func (o *OverviewScreen) HandleInput(event input.InputEvent) {
	switch event.KeyCode {
	case "DOWN":
		o.textComponent.ScrollDown()
	case "UP":
		o.textComponent.ScrollUp()
	case "B":
		vars.CurrentTester = ""
		vars.CurrentScreen = "reviews_screen"
		o.initialized = false
	}
}

func (o *OverviewScreen) Draw() {
	o.InitOverview()

	o.renderer.SetDrawColor(0, 0, 0, 255) // Background color
	o.renderer.Clear()

	sdlutils.RenderTextureCartesian(o.renderer, "assets/textures/bg.bmp", "Q2", "Q4")

	sdlutils.RenderTextureCartesian(o.renderer, "assets/textures/bg_overlay.bmp", "Q2", "Q4")

	// Draw the title
	sdlutils.DrawText(o.renderer, "Overview", sdl.Point{X: 25, Y: 25}, vars.Colors.WHITE, vars.HeaderFont)

	// Draw the text component with scrolling
	o.textComponent.Draw(vars.Colors.WHITE)

	sdlutils.RenderTextureCartesian(o.renderer, "assets/textures/$aspect_ratio/ui_controls.bmp", "Q3", "Q4")

	o.renderer.Present()
}
