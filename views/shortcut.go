/*
   file:           views/dialog_shortcuts.go
   description:    Antarmuka Fyne untuk dialog shortcuts
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package views

import (
	"strings"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func activateDialogPrimaryAction(w fyne.Window) bool {
	pop := dialogOverlay(w)
	if pop == nil {
		return false
	}
	return tapDialogPrimaryButton(pop)
}

func activateDialogActionByLabel(w fyne.Window, label string) bool {
	pop := dialogOverlay(w)
	if pop == nil {
		return false
	}
	return tapDialogButtonByLabel(pop, label)
}

func tapDialogPrimaryButton(pop *widget.PopUp) bool {
	for _, btn := range getDialogButtons(pop) {
		if btn == nil || btn.Disabled() {
			continue
		}
		if btn.Importance == widget.HighImportance {
			btn.Tapped(nil)
			return true
		}
	}

	for _, btn := range getDialogButtons(pop) {
		if btn == nil || btn.Disabled() {
			continue
		}
		btn.Tapped(nil)
		return true
	}

	return false
}

func tapDialogButtonByLabel(pop *widget.PopUp, label string) bool {
	if strings.TrimSpace(label) == "" {
		return false
	}

	for _, btn := range getDialogButtons(pop) {
		if btn == nil || btn.Disabled() || btn.Importance != widget.HighImportance {
			continue
		}

		if matchesLabel(btn.Text, label) {
			btn.Tapped(nil)
			return true
		}
	}

	return false
}

func dialogOverlay(w fyne.Window) *widget.PopUp {
	if w == nil {
		return nil
	}
	canvas := w.Canvas()
	if canvas == nil {
		return nil
	}
	topOverlay := canvas.Overlays().Top()
	pop, _ := topOverlay.(*widget.PopUp)
	return pop
}

func getDialogButtons(pop *widget.PopUp) []*widget.Button {
	if pop == nil {
		return nil
	}

	dialogContent, ok := pop.Content.(*fyne.Container)
	if !ok || len(dialogContent.Objects) <= 3 {
		return nil
	}

	buttonsContainer, ok := dialogContent.Objects[3].(*fyne.Container)
	if !ok {
		return nil
	}

	var buttons []*widget.Button
	for _, obj := range buttonsContainer.Objects {
		if btn, ok := obj.(*widget.Button); ok {
			buttons = append(buttons, btn)
		}
	}
	return buttons
}

func matchesLabel(text, label string) bool {
	text = strings.ToLower(strings.TrimSpace(text))
	label = strings.ToLower(strings.TrimSpace(label))
	if text == "" || label == "" {
		return false
	}
	if text == label {
		return true
	}
	if strings.HasPrefix(text, label) {
		return true
	}
	return firstRune(text) == firstRune(label)
}

func firstRune(s string) rune {
	if s == "" {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(s)
	return r
}
