/*
   file:           views/ecb/double/refrig_po.go
   description:    Layar ECB untuk refrig PO double
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package double

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

// RefrigPoDoubleScreen adalah fungsi untuk refrigerator PO double screen.
func RefrigPoDoubleScreen(w fyne.Window, db *gorp.DbMap, simoConfig configs.SimoConfig) fyne.CanvasObject {
	lineIDs := strings.Split(simoConfig.EcbLineIds, ",")
	lineA := ecb.FirstNonEmpty(lineIDs)
	lineB := ""
	if len(lineIDs) > 1 {
		lineB = ecb.FirstNonEmpty(lineIDs[1:])
	}
	if lineA == "" {
		lineA = "Line A"
	}
	if loc := strings.TrimSpace(simoConfig.EcbLocation); loc != "" {
		lineA = fmt.Sprintf("Line: %s", lineA)
		lineB = fmt.Sprintf("Line: %s", lineB)
	}

	cardWidth := ecb.DeriveLineCardWidth(w, 2, 560, 720)
	left := ecb.NewRefrigLine(lineA, w, cardWidth)
	right := ecb.NewRefrigLineWithAutoFocus(lineB, w, false, cardWidth)

	left.SetResponsiveWidth(2, 560, 720)
	right.SetResponsiveWidth(2, 560, 720)

	ecb.ConfigureRefrigFlow(left, db, simoConfig, 0, true)
	ecb.ConfigureRefrigFlow(right, db, simoConfig, 1, true)

	left.SetPeerFocus(right.Focus, true)
	right.SetPeerFocus(left.Focus, true)

	content := ecb.NewResponsiveLineLayout(left.Canvas(), right.Canvas())

	return container.NewBorder(
		container.NewVBox(widget.NewSeparator()),
		nil,
		nil,
		nil,
		container.NewVScroll(content),
	)
}
