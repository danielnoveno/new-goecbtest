/*
   file:           views/ecb/single/refrig_po.go
   description:    Layar ECB untuk refrig PO single
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package single

import (
	"fmt"
	"go-ecb/configs"
	"go-ecb/views/ecb"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/go-gorp/gorp"
)

func RefrigPoSingleScreen(w fyne.Window, db *gorp.DbMap, simoConfig configs.SimoConfig) fyne.CanvasObject {
	lineIDs := strings.Split(simoConfig.EcbLineIds, ",")
	lineA := ecb.FirstNonEmpty(lineIDs)
	if lineA == "" {
		lineA = "Line A"
	}
	if loc := strings.TrimSpace(simoConfig.EcbLocation); loc != "" {
		lineA = fmt.Sprintf("Line: %s", lineA)
	}

	line := ecb.NewRefrigLine(lineA, w, ecb.DeriveLineCardWidth(w, 1, 760, 900))
	line.SetResponsiveWidth(1, 760, 900)
	ecb.ConfigureRefrigFlow(line, db, simoConfig, 0, true)

	return container.NewBorder(
		container.NewVBox(widget.NewSeparator()),
		nil,
		nil,
		nil,
		container.NewVScroll(line.Canvas()),
	)
}
