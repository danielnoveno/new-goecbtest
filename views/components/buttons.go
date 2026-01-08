/*
   file:           views/components/buttons.go
   description:    Komponen UI umum untuk buttons
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	defaultButtonWidth float32 = 140
)

func PrimaryButton(label string, tapped func()) *widget.Button {
	btn := widget.NewButton(label, tapped)
	btn.Importance = widget.HighImportance
	return btn
}

func SecondaryButton(label string, tapped func()) *widget.Button {
	btn := widget.NewButton(label, tapped)
	btn.Importance = widget.MediumImportance
	return btn
}

func TextButton(label string, tapped func()) *widget.Button {
	btn := widget.NewButton(label, tapped)
	btn.Importance = widget.LowImportance
	return btn
}

func DangerButton(label string, tapped func()) *widget.Button {
	btn := widget.NewButtonWithIcon(label, theme.CancelIcon(), tapped)
	btn.Importance = widget.HighImportance
	return btn
}

func ButtonGroup(buttons ...*widget.Button) fyne.CanvasObject {
	if len(buttons) == 0 {
		return widget.NewLabel("")
	}

	cells := make([]fyne.CanvasObject, 0, len(buttons))
	for _, btn := range buttons {
		if btn == nil {
			continue
		}
		min := btn.MinSize()
		if min.Width < defaultButtonWidth {
			min.Width = defaultButtonWidth
		}
		cell := container.NewGridWrap(min, container.NewCenter(btn))
		cells = append(cells, cell)
	}

	row := container.NewHBox()
	for _, cell := range cells {
		row.Add(container.NewPadded(cell))
	}
	return container.NewHBox(layout.NewSpacer(), row)
}

func ActionButtonRow(buttons ...*widget.Button) fyne.CanvasObject {
	row := container.NewHBox()
	for _, btn := range buttons {
		if btn == nil {
			continue
		}
		row.Add(container.NewPadded(btn))
	}
	row.Add(layout.NewSpacer())
	return row
}
