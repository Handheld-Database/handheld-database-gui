package geometry

import (
	"handheldui/vars"

	"github.com/veandco/go-sdl2/sdl"
)

// Element represents a UI element with padding and margin
type Element struct {
	Width, Height int32
	Padding       int32
	Margin        int32
	Texture       *sdl.Texture
	Position      string // "top-left", "top-right", "bottom-left", "bottom-right"
}

// NewElement creates a new element with the given dimensions and texture
func NewElement(width, height, padding, margin int32, position string) *Element {
	return &Element{
		Width:    width,
		Height:   height,
		Padding:  padding,
		Margin:   margin,
		Position: position,
	}
}

// GetPaddedRect returns the rect considering padding
func (e *Element) GetPaddedRect(x, y int32) sdl.Rect {
	return sdl.Rect{
		X: x + e.Padding,
		Y: y + e.Padding,
		W: e.Width - 2*e.Padding,
		H: e.Height - 2*e.Padding,
	}
}

// GetPosition calculates the element's position relative to the window borders
func (e *Element) GetPosition() sdl.Rect {
	var x, y int32
	switch e.Position {
	case "top-left":
		x = e.Margin
		y = e.Margin
	case "top-right":
		x = vars.ScreenWidth - e.Width - e.Margin
		y = e.Margin
	case "bottom-left":
		x = e.Margin
		y = vars.ScreenHeight - e.Height - e.Margin
	case "bottom-right":
		x = vars.ScreenWidth - e.Width - e.Margin
		y = vars.ScreenHeight - e.Height - e.Margin
	default:
		x = e.Margin
		y = e.Margin
	}
	return e.GetPaddedRect(x, y)
}
