/*
    file:           views/errors/forbidden.go
    description:    Antarmuka Fyne untuk forbidden
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package errors

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ForbiddenPage menampilkan widget yang menjelaskan error 403 berpusat di layar.
func ForbiddenPage(onHome func()) fyne.CanvasObject {
	title := canvas.NewText("Error 403", color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	title.TextSize = 72
	title.TextStyle = fyne.TextStyle{Bold: true}

	message := widget.NewLabel("Maaf, Anda tidak punya akses ke laman ini.")
	message.Alignment = fyne.TextAlignCenter
	message.TextStyle = fyne.TextStyle{Bold: true}

	btn := widget.NewButtonWithIcon("Ke Beranda", theme.HomeIcon(), func() {
		if onHome != nil {
			onHome()
		}
	})
	btn.Importance = widget.HighImportance

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		message,
		layout.NewSpacer(),
		btn,
	)

	bg := canvas.NewRectangle(color.NRGBA{R: 0x22, G: 0x2b, B: 0x3b, A: 0xff})

	return container.NewMax(
		bg,
		container.NewCenter(container.NewVBox(container.NewCenter(content))),
	)
}
