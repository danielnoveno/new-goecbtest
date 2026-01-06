/*
    file:           views/components/page.go
    description:    Komponen UI umum untuk page
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	defaultFormWidth float32 = 540
)

type FormPageConfig struct {
	Title       string
	Description string
	Content     fyne.CanvasObject
	Actions     []fyne.CanvasObject
	HeaderExtra fyne.CanvasObject
	Width       float32
	Fluid       bool
}

// FormPage adalah fungsi untuk form page.
func FormPage(cfg FormPageConfig) fyne.CanvasObject {
	if cfg.Content == nil {
		cfg.Content = widget.NewLabel("")
	}

	title := widget.NewLabel(cfg.Title)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignLeading

	description := widget.NewLabel(cfg.Description)
	description.Wrapping = fyne.TextWrapWord
	description.Alignment = fyne.TextAlignLeading

	header := container.NewVBox(
		title,
		description,
	)
	if cfg.HeaderExtra != nil {
		header.Add(cfg.HeaderExtra)
	}

	body := container.NewVBox(
		header,
		widget.NewSeparator(),
		cfg.Content,
	)

	if len(cfg.Actions) > 0 {
		body.Add(widget.NewSeparator())
		body.Add(container.NewVBox(cfg.Actions...))
	}

	width := cfg.Width
	if width <= 0 {
		width = defaultFormWidth
	}

	var grid fyne.CanvasObject
	if cfg.Fluid {
		grid = container.NewVBox(container.NewPadded(body))
	} else {
		grid = container.NewGridWrap(fyne.NewSize(width, 0), container.NewPadded(body))
	}

	return container.NewVBox(
		grid,
		layout.NewSpacer(),
	)
}
