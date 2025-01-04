package components

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ProgressBarComponent struct {
	renderer      *sdl.Renderer
	progress      float64
	width         int32
	height        int32
	x, y          int32
	borderColor   sdl.Color
	progressColor sdl.Color
}

func (p *ProgressBarComponent) GetProgress() float64 {
	return p.progress
}

func NewProgressBarComponent(renderer *sdl.Renderer, width, height, x, y int32, borderColor, progressColor sdl.Color) *ProgressBarComponent {
	return &ProgressBarComponent{
		renderer:      renderer,
		width:         width,
		height:        height,
		x:             x,
		y:             y,
		borderColor:   borderColor,
		progressColor: progressColor,
	}
}

func (p *ProgressBarComponent) SetProgress(progress float64) {
	if progress < 0 {
		p.progress = 0
	} else if progress > 100 {
		p.progress = 100
	} else {
		p.progress = progress
	}
}

func (p *ProgressBarComponent) Draw() {
	// Draws the border of the progress bar
	p.renderer.SetDrawColor(p.borderColor.R, p.borderColor.G, p.borderColor.B, p.borderColor.A)
	p.renderer.FillRect(&sdl.Rect{
		X: p.x - 8,
		Y: p.y - 8,
		W: p.width + 16,
		H: p.height + 16,
	})

	// Draws the filled part of the progress bar
	p.renderer.SetDrawColor(p.progressColor.R, p.progressColor.G, p.progressColor.B, p.progressColor.A)
	p.renderer.FillRect(&sdl.Rect{
		X: p.x,
		Y: p.y,
		W: int32(float64(p.width) * (p.progress / 100)),
		H: p.height,
	})
}
