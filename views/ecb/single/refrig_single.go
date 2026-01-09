/*
   file:           views/ecb/single/refrig_single.go
   description:    Layar ECB untuk refrig single
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

func RefrigSingleScreen(w fyne.Window, db *gorp.DbMap, simoConfig configs.SimoConfig) fyne.CanvasObject {
	lineName := ecb.FirstNonEmpty(strings.Split(simoConfig.EcbLineIds, ","))
	if lineName == "" {
		lineName = "Line Utama"
	}
	if loc := strings.TrimSpace(simoConfig.EcbLocation); loc != "" {
		lineName = fmt.Sprintf("Line: %s", lineName)
	}

	line := ecb.NewRefrigLine(lineName, w, ecb.DeriveLineCardWidth(w, 1, 760, 900))
	line.SetResponsiveWidth(1, 760, 900)
	ecb.ConfigureRefrigFlow(line, db, simoConfig, 0, false)

	return container.NewBorder(
		container.NewVBox(widget.NewSeparator()),
		nil,
		nil,
		nil,
		container.NewVScroll(line.Canvas()),
	)
}
