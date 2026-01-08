/*
    file:           views/components/dialogs.go
    description:    Komponen UI umum untuk dialogs
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package components

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowError(w fyne.Window, title string, err error) {
	if w == nil || err == nil {
		return
	}

	header := strings.TrimSpace(title)
	if header == "" {
		header = "There is an error"
	}

	message := err.Error()
	if header != "" && !strings.Contains(strings.ToLower(message), strings.ToLower(header)) {
		message = fmt.Sprintf("%s: %s", header, message)
	}

	accent := color.RGBA{R: 255, G: 110, B: 110, A: 255}
	shadow := color.RGBA{R: 32, G: 34, B: 37, A: 255}

	bg := canvas.NewRectangle(shadow)
	bg.SetMinSize(fyne.NewSize(360, 180))
	bg.CornerRadius = 15

	titleText := canvas.NewText(header, accent)
	titleText.TextSize = 17
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.Alignment = fyne.TextAlignLeading

	bodyText := widget.NewLabel(message)
	bodyText.Wrapping = fyne.TextWrapWord
	bodyText.Alignment = fyne.TextAlignLeading

	content := container.NewVBox(
		titleText,
		widget.NewSeparator(),
		bodyText,
	)

	card := container.NewMax(
		bg,
		container.NewPadded(content),
	)

	dialog.NewCustom(header, "Close" , card, w).Show()
}

func ShowInfo(w fyne.Window, title, message string) {
	if w == nil {
		return
	}
	dialog.ShowInformation(title, message, w)
}

func Confirm(w fyne.Window, title, message, confirmLabel, dismissLabel string, cb func(bool)) {
	if w == nil {
		return
	}
	content := widget.NewLabel(message)
	content.Wrapping = fyne.TextWrapWord
	dialog.ShowCustomConfirm(title, confirmLabel, dismissLabel, content, func(ok bool) {
		if cb != nil {
			cb(ok)
		}
	}, w)
}
