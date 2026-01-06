/*
   file:           views/ecb/double/refrig_double.go
   description:    Layar ECB untuk refrig double
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package double

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

// RefrigDoubleScreen adalah fungsi untuk refrigerator double screen.
func RefrigDoubleScreen(w fyne.Window, db *gorp.DbMap, simoConfig configs.SimoConfig) fyne.CanvasObject {
	lineIDs := strings.Split(simoConfig.EcbLineIds, ",")
	lineA := ecb.FirstNonEmpty(lineIDs)
	lineB := ""
	if len(lineIDs) > 1 {
		lineB = ecb.FirstNonEmpty(lineIDs[1:])
	}
	if lineA == "" {
		lineA = "Line A"
	}
	if lineB == "" {
		lineB = "Line B"
	}
	if loc := strings.TrimSpace(simoConfig.EcbLocation); loc != "" {
		lineA = fmt.Sprintf("Line: %s", lineA)
		lineB = fmt.Sprintf("Line: %s", lineB)
	}

	cardWidth := ecb.DeriveLineCardWidth(w, 2, 460, 560)
	left := ecb.NewRefrigLine(lineA, w, cardWidth)
	right := ecb.NewRefrigLineWithAutoFocus(lineB, w, false, cardWidth)

	left.SetResponsiveWidth(2, 460, 560)
	right.SetResponsiveWidth(2, 460, 560)

	ecb.ConfigureRefrigFlow(left, db, simoConfig, 0, false)
	ecb.ConfigureRefrigFlow(right, db, simoConfig, 1, false)

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
