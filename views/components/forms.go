/*
    file:           views/components/forms.go
    description:    Komponen UI umum untuk forms
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
	defaultLabelWidth float32 = 160
)

type FormField struct {
	Label string
	Input fyne.CanvasObject
	Hint  string
}

// NewForm adalah fungsi untuk baru form.
func NewForm(fields ...FormField) fyne.CanvasObject {
	pairs := make([]fyne.CanvasObject, 0, len(fields)*2)
	for _, field := range fields {
		if field.Input == nil {
			continue
		}

		label := widget.NewLabelWithStyle(field.Label, fyne.TextAlignLeading, fyne.TextStyle{})
		label.Wrapping = fyne.TextWrapWord

		min := label.MinSize()
		if min.Width < defaultLabelWidth {
			min.Width = defaultLabelWidth
		}
		labelContainer := container.NewGridWrap(min, label)

		input := field.Input
		if field.Hint != "" {
			hint := widget.NewLabelWithStyle(field.Hint, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
			hint.Wrapping = fyne.TextWrapWord
			input = container.NewVBox(field.Input, hint)
		}

		pairs = append(pairs, labelContainer, input)
	}

	return container.New(layout.NewFormLayout(), pairs...)
}

// StatusBanner adalah fungsi untuk status banner.
func StatusBanner(message string) fyne.CanvasObject {
	if message == "" {
		return widget.NewLabel("")
	}

	lbl := widget.NewLabel(message)
	lbl.Wrapping = fyne.TextWrapWord
	card := widget.NewCard("", "", lbl)
	card.SetSubTitle("")
	return card
}
