package geometry

import (
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// WrapText splits a long text into multiple lines based on the specified maximum width.
func WrapText(text string, font *ttf.Font, maxWidth int) []string {
	words := strings.Fields(text)
	var lines []string
	var currentLine string

	for _, word := range words {
		lineWithWord := currentLine + word + " "
		lineWidth := textWidth(font, lineWithWord)

		if lineWidth > maxWidth {
			if len(currentLine) > 0 {
				lines = append(lines, strings.TrimSpace(currentLine))
			}
			currentLine = word + " "
		} else {
			currentLine = lineWithWord
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, strings.TrimSpace(currentLine))
	}

	return lines
}

// textWidth calculates the width of a string of text based on the provided font.
func textWidth(font *ttf.Font, text string) int {
	surface, err := font.RenderUTF8Blended(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return 0
	}
	defer surface.Free()

	return int(surface.W)
}
