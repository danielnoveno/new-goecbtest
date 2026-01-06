/*
    file:           app/types/app_state.go
    description:    Model dan helper UI untuk ui
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import (
	"fyne.io/fyne/v2"
	"github.com/go-gorp/gorp"
)

type MenuItem struct {
	Title string
	Key   string
	Icon  string
	NavID int
	Show  func() fyne.CanvasObject
}

type AppState struct {
	fyne.App
	fyne.Window
	*gorp.DbMap
	SetBody       func(fyne.CanvasObject)
	Flash         FlashNotifier
	Menu          []MenuItem
	Location      string
	Mode          string
	TryOpenWindow func(title string, setup func(fyne.Window)) bool
}
