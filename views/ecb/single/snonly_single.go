/*
   file:           views/ecb/single/snonly_single.go
   description:    Layar ECB untuk snonly single
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package single

import (
	"fmt"
	"strings"

	"go-ecb/configs"
	"go-ecb/views/ecb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/go-gorp/gorp"
)

func SnOnlySingleScreen(w fyne.Window, db *gorp.DbMap, simoConfig configs.SimoConfig) fyne.CanvasObject {
	lineName := ecb.FirstNonEmpty(strings.Split(simoConfig.EcbLineIds, ","))
	if lineName == "" {
		lineName = "x (Line Kosong)"
	}
	if loc := strings.TrimSpace(simoConfig.EcbLocation); loc != "" {
		lineName = fmt.Sprintf("Line - %s", lineName)
	}

	line := ecb.NewSnOnlyLine(lineName, w, ecb.DeriveLineCardWidth(w, 1, 900, 1100))
	line.SetResponsiveWidth(1, 900, 1100)
	ecb.ConfigureSnOnlyLine(line, db, simoConfig, 0)

	return container.NewBorder(
		container.NewVBox(widget.NewSeparator()),
		nil,
		nil,
		nil,
		line.Canvas(),
	)
}
