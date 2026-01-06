/*
    file:           views/ecb/about.go
    description:    Layar ECB untuk about
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package ecb

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// AboutPage adalah fungsi untuk tentang page.
func AboutPage(iconName, appName, version, author string) fyne.CanvasObject {
	var iconWidget fyne.CanvasObject
	if iconName != "" {
		iconWidget = widget.NewIcon(theme.InfoIcon())
	}

	title := widget.NewRichTextFromMarkdown("**About**")
	body := widget.NewRichTextFromMarkdown(
		"App Name : " + appName +
			"\n\nApp version: " + version +
			"\n\nProgrammed by " + author,
	)

	headerContent := []fyne.CanvasObject{title}
	if iconWidget != nil {
		headerContent = append([]fyne.CanvasObject{iconWidget}, headerContent...)
	}
	header := container.NewHBox(headerContent...)


	card := container.NewVBox(header, body)
	scroll := container.NewVScroll(card)
	scroll.SetMinSize(fyne.NewSize(480, 360))
	return scroll
}
