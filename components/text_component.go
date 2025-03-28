package components

import (
	"handheldui/helpers/sdlutils"
	"handheldui/output"
	"handheldui/vars"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type TextComponent struct {
	renderer        *sdl.Renderer
	text            string
	lines           []string
	scrollOffset    int
	maxVisibleLines int
	font            *ttf.Font
	maxWidth        int
}

func NewTextComponent(renderer *sdl.Renderer, text string, font *ttf.Font, maxVisibleLines int, maxWidth int) *TextComponent {
	visibleLines := maxVisibleLines
	if maxVisibleLines > vars.Config.Screen.MaxLines {
		visibleLines = maxVisibleLines
	}

	component := &TextComponent{
		renderer:        renderer,
		text:            text,
		maxVisibleLines: visibleLines,
		font:            font,
		maxWidth:        maxWidth,
	}
	component.splitTextToLines()
	return component
}

func (t *TextComponent) splitTextToLines() {
	lines := strings.Split(t.text, "\n")

	var wrappedLines []string
	for _, line := range lines {
		wrappedLines = append(wrappedLines, t.wrapLine(line, t.maxWidth)...)
	}

	t.lines = wrappedLines
}

func (t *TextComponent) wrapLine(line string, maxWidth int) []string {
	var wrappedLines []string
	words := strings.Split(line, " ")
	if len(words) == 0 {
		return []string{""}
	}

	currentLine := words[0]

	for _, word := range words[1:] {
		width, _, _ := t.font.SizeUTF8(currentLine + " " + word)
		if width <= maxWidth {
			currentLine += " " + word
		} else {
			wrappedLines = append(wrappedLines, currentLine)
			currentLine = word
		}
	}
	wrappedLines = append(wrappedLines, currentLine)
	return wrappedLines
}

func (t *TextComponent) ScrollDown() {
	if t.scrollOffset < len(t.lines)-t.maxVisibleLines {
		t.scrollOffset++
	}
}

func (t *TextComponent) ScrollUp() {
	if t.scrollOffset > 0 {
		t.scrollOffset--
	}
}

func (t *TextComponent) Draw(primaryColor sdl.Color) {
	startIndex := t.scrollOffset
	endIndex := startIndex + t.maxVisibleLines
	if endIndex > len(t.lines) {
		endIndex = len(t.lines)
	}
	visibleLines := t.lines[startIndex:endIndex]

	for index, line := range visibleLines {
		textSurface, err := sdlutils.RenderText(line, primaryColor, t.font)
		if err != nil {
			output.Printf("Error rendering text: %v\n", err)
			return
		}
		defer textSurface.Free()

		texture, err := t.renderer.CreateTextureFromSurface(textSurface)
		if err != nil {
			output.Printf("Error creating texture: %v\n", err)
			return
		}
		defer texture.Destroy()

		t.renderer.Copy(texture, nil, &sdl.Rect{X: 40, Y: 90 + 30*int32(index), W: textSurface.W, H: textSurface.H})
	}
}

func (t *TextComponent) GetScrollOffset() int {
	return t.scrollOffset
}

func (t *TextComponent) GetLines() []string {
	return t.lines
}
