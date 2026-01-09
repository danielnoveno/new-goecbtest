/*
    file:           views/components/surface.go
    description:    Komponen UI umum untuk surface
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	defaultSurfaceBackground = color.RGBA{R: 0x12, G: 0x1B, B: 0x2A, A: 0xFF}
	defaultSurfaceBorder     = color.RGBA{R: 0x2A, G: 0x3D, B: 0x5B, A: 0x55}
)

type SurfaceConfig struct {
	Background   color.Color
	BorderColor  color.Color
	CornerRadius float32
	MinSize      fyne.Size
	Content      fyne.CanvasObject
}

func Surface(cfg SurfaceConfig) fyne.CanvasObject {
	background := cfg.Background
	if background == nil {
		background = defaultSurfaceBackground
	}

	border := cfg.BorderColor
	if border == nil {
		border = defaultSurfaceBorder
	}

	rect := canvas.NewRectangle(background)
	rect.StrokeColor = border
	rect.StrokeWidth = 1
	if cfg.CornerRadius > 0 {
		rect.CornerRadius = cfg.CornerRadius
	}
	if cfg.MinSize.Width > 0 || cfg.MinSize.Height > 0 {
		rect.SetMinSize(cfg.MinSize)
	}

	content := cfg.Content
	if content == nil {
		content = widget.NewLabel("")
	}

	card := container.New(layout.NewMaxLayout(),
		rect,
		container.New(layout.NewVBoxLayout(), container.NewPadded(content)),
	)

	return card
}
