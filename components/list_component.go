package components

import (
	"handheldui/helpers"
	"handheldui/output"
	"handheldui/vars"

	"github.com/veandco/go-sdl2/sdl"
)

type ListComponent struct {
	renderer        *sdl.Renderer
	items           []map[string]interface{}
	selectedIndex   int
	scrollOffset    int
	title           string
	itemFormatter   func(index int, item map[string]interface{}) string
	maxVisibleItems int
}

func NewListComponent(renderer *sdl.Renderer, title string, maxVisibleItems int, itemFormatter func(index int, item map[string]interface{}) string) *ListComponent {
	return &ListComponent{
		renderer:        renderer,
		title:           title,
		itemFormatter:   itemFormatter,
		maxVisibleItems: maxVisibleItems,
		items:           []map[string]interface{}{},
	}
}

func (l *ListComponent) SetItems(items []map[string]interface{}) {
	l.items = items
	l.selectedIndex = 0
	l.scrollOffset = 0
}

func (l *ListComponent) ScrollDown() {
	if l.selectedIndex < len(l.items)-1 {
		l.selectedIndex++
		if l.selectedIndex >= l.scrollOffset+l.maxVisibleItems {
			l.scrollOffset++
		}
	}
}

func (l *ListComponent) ScrollUp() {
	if l.selectedIndex > 0 {
		l.selectedIndex--
		if l.selectedIndex < l.scrollOffset {
			l.scrollOffset--
		}
	}
}

func (l *ListComponent) Draw(primaryColor sdl.Color, selectedColor sdl.Color) {
	// Draw the items
	startIndex := l.scrollOffset
	endIndex := startIndex + l.maxVisibleItems
	if endIndex > len(l.items) {
		endIndex = len(l.items)
	}
	visibleItems := l.items[startIndex:endIndex]

	for index, item := range visibleItems {
		color := primaryColor
		if index+startIndex == l.selectedIndex {
			color = selectedColor
		}
		itemText := l.itemFormatter(index+startIndex, item)
		textSurface, err := helpers.RenderText(itemText, color, vars.BodyFont)
		if err != nil {
			output.Printf("Error rendering text: %v\n", err)
			return
		}
		defer textSurface.Free()

		texture, err := l.renderer.CreateTextureFromSurface(textSurface)
		if err != nil {
			output.Printf("Error creating texture: %v\n", err)
			return
		}
		defer texture.Destroy()

		l.renderer.Copy(texture, nil, &sdl.Rect{X: 40, Y: 90 + 30*int32(index), W: textSurface.W, H: textSurface.H})
	}
}

func (l *ListComponent) GetSelectedIndex() int {
	return l.selectedIndex
}

func (l *ListComponent) GetScrollOffset() int {
	return l.scrollOffset
}

func (l *ListComponent) GetItems() []map[string]interface{} {
	return l.items
}
