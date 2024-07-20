package components

import (
	"fmt"
	"handheldui/helpers"
	"handheldui/vars"

	"github.com/veandco/go-sdl2/sdl"
)

type ListComponent struct {
	renderer      *sdl.Renderer
	items         []map[string]interface{}
	selectedIndex int
	scrollOffset  int
	title         string
	itemFormatter func(index int, item map[string]interface{}) string
}

func NewListComponent(renderer *sdl.Renderer, title string, itemFormatter func(index int, item map[string]interface{}) string) *ListComponent {
	return &ListComponent{
		renderer:      renderer,
		title:         title,
		itemFormatter: itemFormatter,
	}
}

func (l *ListComponent) SetItems(items []map[string]interface{}, selectedIndex, scrollOffset int) {
	l.items = items
	l.selectedIndex = selectedIndex
	l.scrollOffset = scrollOffset
}

func (l *ListComponent) Draw() {
	// Desenhe o título
	titleColor := sdl.Color{R: 0, G: 0, B: 255, A: 255} // Azul para o título
	titleSurface, err := helpers.RenderText(l.title, titleColor, vars.HeaderFont)
	if err != nil {
		fmt.Printf("Erro ao renderizar texto do título: %v\n", err)
		return
	}
	defer titleSurface.Free()

	titleTexture, err := l.renderer.CreateTextureFromSurface(titleSurface)
	if err != nil {
		fmt.Printf("Erro ao criar textura do título: %v\n", err)
		return
	}
	defer titleTexture.Destroy()

	l.renderer.Copy(titleTexture, nil, &sdl.Rect{X: 40, Y: 20, W: int32(titleSurface.W), H: int32(titleSurface.H)})

	// Desenhe os itens
	startIndex := l.scrollOffset
	endIndex := startIndex + int(vars.ScreenHeight/30)
	if endIndex > len(l.items) {
		endIndex = len(l.items)
	}
	visibleItems := l.items[startIndex:endIndex]

	for index, item := range visibleItems {
		color := sdl.Color{R: 255, G: 255, B: 255, A: 255}
		if index == l.selectedIndex-startIndex {
			color = sdl.Color{R: 128, G: 0, B: 128, A: 255}
		}
		itemText := l.itemFormatter(index, item)
		textSurface, err := helpers.RenderText(itemText, color, vars.BodyFont)
		if err != nil {
			fmt.Printf("Erro ao renderizar texto: %v\n", err)
			return
		}
		defer textSurface.Free()

		texture, err := l.renderer.CreateTextureFromSurface(textSurface)
		if err != nil {
			fmt.Printf("Erro ao criar textura: %v\n", err)
			return
		}
		defer texture.Destroy()

		l.renderer.Copy(texture, nil, &sdl.Rect{X: 40, Y: 60 + 30*int32(index), W: textSurface.W, H: textSurface.H})
	}
}
