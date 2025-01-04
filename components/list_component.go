package components

import (
	"handheldui/helpers/sdlutils"
	"handheldui/output"
	"handheldui/vars"

	"github.com/veandco/go-sdl2/sdl"
)

type ListComponent struct {
	renderer        *sdl.Renderer
	items           []map[string]interface{}
	selectedIndex   int
	scrollOffset    int
	itemFormatter   func(index int, item map[string]interface{}) string
	maxVisibleItems int
}

func NewListComponent(renderer *sdl.Renderer, maxVisibleItems int, itemFormatter func(index int, item map[string]interface{}) string) *ListComponent {
	return &ListComponent{
		renderer:        renderer,
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

func (l *ListComponent) PageDown() {
	if l.selectedIndex < len(l.items)-1 {
		l.selectedIndex += l.maxVisibleItems
		if l.selectedIndex >= len(l.items) {
			l.selectedIndex = len(l.items) - 1
		}
		l.scrollOffset = l.selectedIndex - (l.selectedIndex % l.maxVisibleItems)
		if l.scrollOffset+l.maxVisibleItems > len(l.items) {
			l.scrollOffset = len(l.items) - l.maxVisibleItems
			if l.scrollOffset < 0 {
				l.scrollOffset = 0
			}
		}
	}
}

func (l *ListComponent) PageUp() {
	if l.selectedIndex > 0 {
		l.selectedIndex -= l.maxVisibleItems
		if l.selectedIndex < 0 {
			l.selectedIndex = 0
		}
		l.scrollOffset = l.selectedIndex - (l.selectedIndex % l.maxVisibleItems)
		if l.scrollOffset < 0 {
			l.scrollOffset = 0
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
		textSurface, err := sdlutils.RenderText(itemText, color, vars.BodyFont)
		if err != nil {
			output.Errorf("Error rendering text: %v\n", err)
			return
		}
		defer textSurface.Free()

		texture, err := l.renderer.CreateTextureFromSurface(textSurface)
		if err != nil {
			output.Errorf("Error creating texture: %v\n", err)
			return
		}
		defer texture.Destroy()

		l.renderer.Copy(texture, nil, &sdl.Rect{X: 40, Y: 90 + 48*int32(index), W: textSurface.W, H: textSurface.H})
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
