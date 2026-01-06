/*
   file:           views/ecb/line_components.go
   description:    Layar ECB untuk line components
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package ecb

import (
	"fmt"
	"image/color"
	"log"
	"strings"
	"sync"
	"time"

	"runtime/debug"

	"go-ecb/services/system"
	navigation "go-ecb/views/ecb/navigation"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xlayout "fyne.io/x/fyne/layout"
)

type lineKind int

const (
	lineRefrig lineKind = iota
	lineSnOnly
)

type lineStepKind int

const (
	stepSPC lineStepKind = iota
	stepSerial
	stepCompressorType
	stepCompressorCode
)

// StepKind adalah alias yang dapat digunakan di luar paket untuk merujuk langkah input.
type StepKind = lineStepKind

const (
	StepSPC            StepKind = stepSPC
	StepSerial         StepKind = stepSerial
	StepCompressorType StepKind = stepCompressorType
	StepCompressorCode StepKind = stepCompressorCode
)

// FirstNonEmpty returns the first non-empty string from the supplied slice.
func FirstNonEmpty(items []string) string {
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

type lineStep struct {
	kind             lineStepKind
	label            string
	requiresHardware bool
	completesFlow    bool
}

type historyRow struct {
	Date   string
	Time   string
	Serial string
	FgType string
}

type historyTable struct {
	headers []string
	rows    []historyRow
	maxRows int
	table   *widget.Table
	mu      sync.Mutex
}

type sizedButton struct {
	widget.Button
	minSize fyne.Size
}

// newSizedButton adalah fungsi untuk baru sized button.
func newSizedButton(label string, icon fyne.Resource, tapped func(), minSize fyne.Size) *sizedButton {
	btn := &sizedButton{
		Button:  *widget.NewButtonWithIcon(label, icon, tapped),
		minSize: minSize,
	}
	btn.ExtendBaseWidget(btn)
	return btn
}

// MinSize adalah fungsi untuk min size.
func (b *sizedButton) MinSize() fyne.Size {
	base := b.Button.MinSize()
	if b.minSize.Width > base.Width {
		base.Width = b.minSize.Width
	}
	if b.minSize.Height > base.Height {
		base.Height = b.minSize.Height
	}
	return base
}

type responsiveWidthSpec struct {
	columns  int
	fallback float32
	max      float32
}

// availableCanvasWidth adalah fungsi untuk available canvas width.
func availableCanvasWidth(w fyne.Window) float32 {
	if w == nil || w.Canvas() == nil {
		return 0
	}
	width := w.Canvas().Size().Width - 64
	if width < 0 {
		return 0
	}
	return width
}

// DeriveLineCardWidth adalah fungsi untuk derive jalur card width.
func DeriveLineCardWidth(w fyne.Window, columns int, fallback, maxWidth float32) float32 {
	if fallback <= 0 {
		fallback = 640
	}
	if columns <= 0 {
		columns = 1
	}
	width := fallback
	if available := availableCanvasWidth(w); available > 0 {
		perColumn := available / float32(columns)
		if perColumn > 0 {
			width = perColumn
		}
	}
	if width <= 0 {
		width = fallback
	}
	if maxWidth > 0 && width > maxWidth {
		width = maxWidth
	}
	return width
}

// newHistoryTable adalah fungsi untuk baru history tabel.
func newHistoryTable(maxRows int) *historyTable {
	h := &historyTable{
		headers: []string{"Date", "Hour", "S/N", "Type"},
		maxRows: maxRows,
		rows:    make([]historyRow, 0, maxRows),
	}

	h.table = widget.NewTable(
		func() (int, int) {
			return len(h.rows) + 1, len(h.headers)
		},
		func() fyne.CanvasObject {
			lbl := widget.NewLabel("")
			lbl.Alignment = fyne.TextAlignCenter
			return lbl
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			lbl := obj.(*widget.Label)
			lbl.Alignment = fyne.TextAlignCenter

			if id.Row == 0 {
				lbl.SetText(h.headers[id.Col])
				lbl.TextStyle = fyne.TextStyle{Bold: true}
				return
			}

			if id.Row-1 >= len(h.rows) {
				lbl.SetText("")
				return
			}

			row := h.rows[id.Row-1]

			switch id.Col {
			case 0:
				lbl.SetText(row.Date)
			case 1:
				lbl.SetText(row.Time)
			case 2:
				lbl.SetText(row.Serial)
			case 3:
				lbl.SetText(row.FgType)
			}

			lbl.TextStyle = fyne.TextStyle{}
		},
	)

	h.adjustColumnWidths(0)

	return h
}

// adjustColumnWidths adalah fungsi untuk adjust column widths.
func (h *historyTable) adjustColumnWidths(total float32) {
	if len(h.headers) == 0 || h.table == nil {
		return
	}
	if total <= 0 {
		total = 420
	}
	perColumn := total / float32(len(h.headers))
	if perColumn < 60 {
		perColumn = 60
	}
	for i := range h.headers {
		h.table.SetColumnWidth(i, perColumn)
	}
}

// append adalah fungsi untuk append.
func (h *historyTable) append(row historyRow) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.maxRows <= 0 {
		h.maxRows = 5
	}

	h.rows = append([]historyRow{row}, h.rows...)
	if len(h.rows) > h.maxRows {
		h.rows = h.rows[:h.maxRows]
	}
	h.table.Refresh()
}

type LineState struct {
	name           string
	kind           lineKind
	window         fyne.Window
	steps          []lineStep
	currentStep    int
	testRunning    bool
	mu             sync.Mutex
	otherFocus     func()
	history        *historyTable
	root           fyne.CanvasObject
	inputEntry     *widget.Entry
	jumboText      *fyne.Container
	jumboBg        *canvas.Rectangle
	lastInput      *widget.Label
	fieldRows      []fyne.CanvasObject
	spcLabel       *widget.Label
	snLabel        *widget.Label
	fgLabel        *widget.Label
	compTypeLabel  *widget.Label
	compCodeLabel  *widget.Label
	cardWidth      float32
	responsiveSpec *responsiveWidthSpec
	jumboFontSize  float32
	focusRetries   int
	headerFontSize float32

	spcValue       string
	serialValue    string
	fgTypeValue    string
	compTypeValue  string
	compCodeValue  string
	stepValidators map[lineStepKind]func(string) error

	customSerialValidator  func(string) (string, error)
	customSaveFlow         func() error
	successMessage         string
	focusPeerAfterComplete bool
	autoFocus              bool
}

// newRefrigLine adalah fungsi untuk baru refrigerator jalur.
func NewRefrigLine(name string, w fyne.Window, width ...float32) *LineState {
	return NewRefrigLineWithAutoFocus(name, w, true, width...)
}

// newRefrigLineWithAutoFocus adalah fungsi untuk baru refrigerator jalur with auto focus.
func NewRefrigLineWithAutoFocus(name string, w fyne.Window, autoFocus bool, width ...float32) *LineState {
	cardWidth := float32(0)
	if len(width) > 0 {
		cardWidth = width[0]
	}
	line := &LineState{
		name:           name,
		window:         w,
		kind:           lineRefrig,
		cardWidth:      cardWidth,
		autoFocus:      autoFocus,
		stepValidators: make(map[lineStepKind]func(string) error),
		steps: []lineStep{
			{kind: stepSPC, label: "SCAN SPC", requiresHardware: true},
			{kind: stepSerial, label: "SCAN NOMOR SERI"},
			{kind: stepCompressorType, label: "SCAN TIPE KOMPRESOR"},
			{kind: stepCompressorCode, label: "SCAN KODE KOMPRESOR", completesFlow: true},
		},
		history:        newHistoryTable(8),
		successMessage: "Data tersimpan",
	}
	line.buildUI()
	return line
}

// newSnOnlyLine adalah fungsi untuk baru sn only jalur.
func NewSnOnlyLine(name string, w fyne.Window, width ...float32) *LineState {
	return NewSnOnlyLineWithAutoFocus(name, w, true, width...)
}

// NewSnOnlyLineWithAutoFocus returns a sn-only line builder that can opt out of initial focus.
func NewSnOnlyLineWithAutoFocus(name string, w fyne.Window, autoFocus bool, width ...float32) *LineState {
	cardWidth := float32(0)
	if len(width) > 0 {
		cardWidth = width[0]
	}
	line := &LineState{
		name:           name,
		window:         w,
		kind:           lineSnOnly,
		cardWidth:      cardWidth,
		headerFontSize: 25,
		stepValidators: make(map[lineStepKind]func(string) error),
		steps: []lineStep{
			{kind: stepSerial, label: "SCAN NOMOR SERI", requiresHardware: true, completesFlow: true},
		},
		history:        newHistoryTable(8),
		successMessage: "Data tersimpan",
		autoFocus:      autoFocus,
	}
	line.buildUI()
	return line
}

// buildUI adalah fungsi untuk menyusun ui.
func (l *LineState) buildUI() {
	l.jumboBg = canvas.NewRectangle(color.Black)
	l.jumboBg.SetMinSize(fyne.NewSize(0, 140))
	l.jumboFontSize = 50
	initialLine := canvas.NewText("READY", color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF})
	initialLine.TextSize = l.jumboFontSize
	initialLine.TextStyle = fyne.TextStyle{Monospace: true}
	initialLine.Alignment = fyne.TextAlignCenter
	l.jumboText = container.NewVBox(initialLine)
	jumbo := container.NewStack(
		l.jumboBg,
		container.NewCenter(l.jumboText),
	)
	l.refreshJumboSizes()

	makeRow := func(title string) (*widget.Label, fyne.CanvasObject) {
		value := widget.NewLabel("-")
		value.Alignment = fyne.TextAlignLeading
		row := container.NewGridWithColumns(2,
			widget.NewLabel(title),
			value,
		)
		return value, row
	}

	l.fieldRows = []fyne.CanvasObject{}
	switch l.kind {
	case lineRefrig:
		if value, row := makeRow("S/N SPC"); row != nil {
			l.spcLabel = value
			l.fieldRows = append(l.fieldRows, row)
		}
		if value, row := makeRow("S/N Produk"); row != nil {
			l.snLabel = value
			l.fieldRows = append(l.fieldRows, row)
		}
		if value, row := makeRow("Tipe Produk S/N"); row != nil {
			l.fgLabel = value
			l.fieldRows = append(l.fieldRows, row)
		}
		if value, row := makeRow("Tipe Kompresor"); row != nil {
			l.compTypeLabel = value
			l.fieldRows = append(l.fieldRows, row)
		}
		if value, row := makeRow("S/N Kompresor"); row != nil {
			l.compCodeLabel = value
			l.fieldRows = append(l.fieldRows, row)
		}
	case lineSnOnly:
		if value, row := makeRow("S/N Produk"); row != nil {
			l.snLabel = value
			l.fieldRows = append(l.fieldRows, row)
		}
		if value, row := makeRow("Tipe Produk S/N"); row != nil {
			l.fgLabel = value
			l.fieldRows = append(l.fieldRows, row)
		}
	}

	l.inputEntry = widget.NewEntry()
	l.inputEntry.SetPlaceHolder("Scan / tipe lalu tekan Enter")
	l.inputEntry.OnSubmitted = func(text string) {
		log.Printf("[%s] OnSubmitted text=%q", l.name, text)
		l.handleInput(text)
	}

	// lastInputTitle := widget.NewLabel("Last Input")
	// l.lastInput = widget.NewLabel("-")
	// lastInputRow := container.NewGridWithColumns(2, lastInputTitle, l.lastInput)

	resetBtn := newSizedButton("Reset", theme.ViewRefreshIcon(), func() {
		l.Reset(true)
	}, fyne.NewSize(170, 30))
	resetBtn.Importance = widget.MediumImportance
	l.inputEntry.Resize(fyne.NewSize(260, 32))

	resetContainer := container.New(layout.NewCenterLayout(), resetBtn)

	label := widget.NewLabelWithStyle("INPUT SCAN", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	form := container.NewBorder(
		nil,
		nil,
		label,
		resetContainer,
		l.inputEntry,
	)

	historyTitle := widget.NewLabelWithStyle("History", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	historyStack := container.NewStack()
	historyScroller := container.NewVScroll(historyStack)
	historyWrap := container.NewBorder(nil, nil, nil, nil,
		container.NewVBox(
			historyTitle,
			historyScroller,
		),
	)

	body := container.NewVBox(
		jumbo,
		widget.NewSeparator(),
		container.NewVBox(l.fieldRows...),
		widget.NewSeparator(),
		// lastInputRow,
		form,
		widget.NewSeparator(),
		historyWrap,
	)

	var headerText fyne.CanvasObject
	var headerTextCanvas *canvas.Text
	if l.headerFontSize > 0 {
		headerTextCanvas = canvas.NewText(l.name, theme.Color(theme.ColorNameForeground))
		headerTextCanvas.Alignment = fyne.TextAlignLeading
		headerTextCanvas.TextStyle = fyne.TextStyle{Bold: true}
		headerTextCanvas.TextSize = l.headerFontSize
		headerText = headerTextCanvas
	} else {
		headerText = widget.NewLabelWithStyle(l.name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	}
	card := widget.NewCard("", "", container.NewVBox(headerText, body))

	if headerTextCanvas != nil {
		if app := fyne.CurrentApp(); app != nil {
			app.Settings().AddListener(func(fyne.Settings) {
				headerTextCanvas.Color = theme.Color(theme.ColorNameForeground)
				headerTextCanvas.Refresh()
			})
		}
	}

	historyWidth := l.responsiveWidthLimit()
	if historyWidth <= 0 {
		historyWidth = availableCanvasWidth(l.window)
	}
	minHistoryHeight := float32(240)
	if historyWidth > 0 {
		padding := float32(24)
		if historyWidth > padding*2 {
			historyWidth = historyWidth - padding*2
		}
		l.history.adjustColumnWidths(historyWidth)
		bg := canvas.NewRectangle(color.Transparent)
		bg.SetMinSize(fyne.NewSize(historyWidth, minHistoryHeight))
		historyStack.Objects = []fyne.CanvasObject{bg, l.history.table}
		historyScroller.SetMinSize(fyne.NewSize(historyWidth, minHistoryHeight))
	} else {
		bg := canvas.NewRectangle(color.Transparent)
		bg.SetMinSize(fyne.NewSize(0, minHistoryHeight))
		historyStack.Objects = []fyne.CanvasObject{bg, l.history.table}
		historyScroller.SetMinSize(fyne.NewSize(0, minHistoryHeight))
	}

	l.root = container.NewStack(card)
	l.showInstruction()
	if l.autoFocus {
		l.refocusInput()
	}
}

// SetResponsiveWidth allows the card to recompute its ideal width per breakpoint.
func (l *LineState) SetResponsiveWidth(columns int, fallback, max float32) {
	if columns <= 0 {
		columns = 1
	}
	if fallback <= 0 {
		fallback = l.cardWidth
	}
	if fallback <= 0 {
		fallback = 640
	}
	l.responsiveSpec = &responsiveWidthSpec{
		columns:  columns,
		fallback: fallback,
		max:      max,
	}
	l.refreshJumboSizes()
}

func (l *LineState) responsiveWidthLimit() float32 {
	if spec := l.responsiveSpec; spec != nil {
		columns := spec.columns
		if columns <= 0 {
			columns = 1
		}
		fallback := spec.fallback
		if fallback <= 0 {
			fallback = 640
		}
		return DeriveLineCardWidth(l.window, columns, fallback, spec.max)
	}
	if l.cardWidth > 0 {
		return l.cardWidth
	}
	return DeriveLineCardWidth(l.window, 1, 640, 0)
}

func (l *LineState) refreshJumboSizes() {
	if l.jumboBg == nil {
		return
	}
	width := l.responsiveWidthLimit()
	if width <= 0 {
		width = 0
	}
	l.jumboBg.SetMinSize(fyne.NewSize(width, 140))
}

func (l *LineState) updateJumboText(lines []string, clr color.Color) {
	if l.jumboText == nil {
		return
	}
	if len(lines) == 0 {
		lines = []string{""}
	}
	objs := make([]fyne.CanvasObject, 0, len(lines))
	for _, lineText := range lines {
		line := canvas.NewText(lineText, clr)
		fontSize := l.jumboFontSize
		if fontSize <= 0 {
			fontSize = 50
		}
		line.TextSize = fontSize
		line.TextStyle = fyne.TextStyle{Monospace: true}
		line.Alignment = fyne.TextAlignCenter
		objs = append(objs, line)
	}
	l.jumboText.Objects = objs
	l.jumboText.Refresh()
}

func (l *LineState) buildJumboLines(msg string) []string {
	if msg == "" {
		return []string{""}
	}

	maxWidth := l.responsiveWidthLimit()
	if maxWidth <= 0 {
		maxWidth = 640
	}

	fontSize := l.jumboFontSize
	if fontSize <= 0 {
		fontSize = 50
	}
	style := fyne.TextStyle{Monospace: true}

	var result []string
	for _, segment := range strings.Split(msg, "\n") {
		trimmed := strings.TrimSpace(segment)
		if trimmed == "" {
			result = append(result, "")
			continue
		}
		result = append(result, l.wrapSegment(trimmed, maxWidth, fontSize, style)...)
	}

	if len(result) == 0 {
		return []string{""}
	}
	return result
}

func (l *LineState) wrapSegment(segment string, maxWidth, fontSize float32, style fyne.TextStyle) []string {
	if maxWidth <= 0 {
		return []string{segment}
	}
	if measureTextWidth(segment, fontSize, style) <= maxWidth {
		return []string{segment}
	}

	var lines []string
	var current strings.Builder
	words := strings.Fields(segment)
	for _, word := range words {
		candidate := current.String()
		if candidate != "" {
			candidate += " "
		}
		candidate += word

		if measureTextWidth(candidate, fontSize, style) <= maxWidth {
			current.Reset()
			current.WriteString(candidate)
			continue
		}

		if current.Len() > 0 {
			lines = append(lines, current.String())
			current.Reset()
		}

		if measureTextWidth(word, fontSize, style) <= maxWidth {
			current.WriteString(word)
			continue
		}

		for _, part := range l.breakWord(word, maxWidth, fontSize, style) {
			lines = append(lines, part)
		}
	}

	if current.Len() > 0 {
		lines = append(lines, current.String())
	}

	if len(lines) == 0 {
		return []string{segment}
	}
	return lines
}

func (l *LineState) breakWord(word string, maxWidth, fontSize float32, style fyne.TextStyle) []string {
	runes := []rune(word)
	var parts []string
	for len(runes) > 0 {
		end := len(runes)
		for end > 0 && measureTextWidth(string(runes[:end]), fontSize, style) > maxWidth {
			end--
		}
		if end == 0 {
			end = 1
		}
		parts = append(parts, string(runes[:end]))
		runes = runes[end:]
	}
	return parts
}

func measureTextWidth(text string, size float32, style fyne.TextStyle) float32 {
	return fyne.MeasureText(text, size, style).Width
}

// canvas adalah fungsi untuk canvas.
func (l *LineState) Canvas() fyne.CanvasObject {
	return l.root
}

// handleInput adalah fungsi untuk menangani input.
func (l *LineState) handleInput(raw string) {
	log.Printf("[%s] handleInput start raw=%q currentStep=%d testRunning=%v", l.name, raw, l.currentStep, l.testRunning)
	skipRefocus := false
	defer func() {
		if !skipRefocus {
			l.refocusInput()
		}
	}()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[%s] panic in handleInput raw=%q: %v\n%s", l.name, raw, r, debug.Stack())
			l.flashStatus("Terjadi error saat proses input", color.RGBA{R: 0xFF, G: 0x45, B: 0x00, A: 0xFF}, 3000*time.Millisecond)
		}
	}()

	text := strings.TrimSpace(raw)
	l.inputEntry.SetText("")
	if text == "" {
		l.flashStatus("Input Kosong", color.RGBA{R: 0xFF, G: 0xA5, B: 0x00, A: 0xFF}, 1500*time.Millisecond)
		return
	}
	l.setLastInput(text)

	if l.handleKeyword(text, &skipRefocus) {
		log.Printf("[%s] handleInput handled as keyword text=%q", l.name, text)
		return
	}

	l.mu.Lock()
	step := l.steps[l.currentStep]
	log.Printf("[%s] step=%d kind=%v input=%q", l.name, l.currentStep, step.kind, text)
	requiresHardware := step.requiresHardware
	completesFlow := step.completesFlow
	autoSwitch := l.focusPeerAfterComplete
	l.mu.Unlock()

	if ok := l.applyInput(step.kind, text); !ok {
		return
	}

	if requiresHardware {
		l.startHardware(step)
		return
	}

	if completesFlow {
		if autoSwitch {
			skipRefocus = true
		}
		l.completeFlow()
		return
	}

	l.advanceStep()
}

// applyInput adalah fungsi untuk apply input.
func (l *LineState) applyInput(kind lineStepKind, input string) bool {
	log.Printf("[%s] applyInput kind=%v input=%q", l.name, kind, input)
	trimmed := strings.TrimSpace(input)
	upper := strings.ToUpper(trimmed)
	if validator, exists := l.stepValidators[kind]; exists && validator != nil {
		if err := validator(trimmed); err != nil {
			l.flashStatus(err.Error(), color.RGBA{R: 0xFF, G: 0x45, B: 0x00, A: 0xFF}, 2000*time.Millisecond)
			return false
		}
	}
	switch kind {
	case stepSPC:
		l.mu.Lock()
		l.spcValue = upper
		if l.spcLabel != nil {
			l.spcLabel.SetText(l.spcValue)
		}
		l.mu.Unlock()
	case stepSerial:
		if l.customSerialValidator != nil {
			fg, err := l.customSerialValidator(trimmed)
			if err != nil {
				l.flashStatus(err.Error(), color.RGBA{R: 0xFF, G: 0xA5, B: 0x00, A: 0xFF}, 1800*time.Millisecond)
				return false
			}
			upper = strings.ToUpper(trimmed)
			l.mu.Lock()
			l.serialValue = upper
			l.fgTypeValue = fg
			if l.snLabel != nil {
				l.snLabel.SetText(l.serialValue)
			}
			if l.fgLabel != nil {
				l.fgLabel.SetText(l.fgTypeValue)
			}
			l.mu.Unlock()
			l.flashStatus("sn ok", color.RGBA{R: 0x00, G: 0xBF, B: 0x8F, A: 0xFF}, 800*time.Millisecond)
			return true
		}
		l.mu.Lock()
		l.serialValue = upper
		if l.snLabel != nil {
			l.snLabel.SetText(l.serialValue)
		}
		l.fgTypeValue = DeriveFgType(l.serialValue)
		if l.fgLabel != nil {
			l.fgLabel.SetText(l.fgTypeValue)
		}
		l.mu.Unlock()
	case stepCompressorType:
		l.mu.Lock()
		l.compTypeValue = upper
		if l.compTypeLabel != nil {
			l.compTypeLabel.SetText(deriveCompressorTitle(l.compTypeValue))
		}
		l.mu.Unlock()
	case stepCompressorCode:
		l.mu.Lock()
		l.compCodeValue = upper
		if l.compCodeLabel != nil {
			l.compCodeLabel.SetText(l.compCodeValue)
		}
		l.mu.Unlock()
	}
	return true
}

// startHardware adalah fungsi untuk menjalankan hardware.
func (l *LineState) startHardware(step lineStep) {
	l.mu.Lock()
	if l.testRunning {
		l.mu.Unlock()
		l.flashStatus("WAIT FOR TEST", color.RGBA{R: 0xFF, G: 0xA5, B: 0x00, A: 0xFF}, 1500*time.Millisecond)
		return
	}
	stepIdx := l.currentStep
	l.testRunning = true
	l.mu.Unlock()

	log.Printf("[%s] startHardware step=%d completes=%v", l.name, stepIdx, step.completesFlow)
	l.flashStatus("WAIT FOR TEST", color.RGBA{R: 0xFF, G: 0xD7, B: 0x00, A: 0xFF}, 0)

	go func(stepIdx int, completes bool) {
		// Hardware simulator: always pass to avoid random false-fail during scans.
		pass := true
		runOnMain(func() {
			l.mu.Lock()
			if l.currentStep != stepIdx {
				l.testRunning = false
				l.mu.Unlock()
				return
			}
			l.testRunning = false
			l.mu.Unlock()
			if pass {
				log.Printf("[%s] hardware PASS step=%d", l.name, stepIdx)
				l.flashStatus("PASS", color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}, 1200*time.Millisecond)
				if completes {
					l.completeFlow()
				} else {
					l.advanceStep()
				}
			} else {
				log.Printf("[%s] hardware FAIL step=%d", l.name, stepIdx)
				l.flashStatus("FAIL", color.RGBA{R: 0xFF, G: 0x45, B: 0x00, A: 0xFF}, 0)
			}
		})
	}(stepIdx, step.completesFlow)
}

// advanceStep adalah fungsi untuk advance step.
func (l *LineState) advanceStep() {
	log.Printf("[%s] advanceStep from=%d", l.name, l.currentStep)
	l.mu.Lock()
	if l.currentStep < len(l.steps)-1 {
		l.currentStep++
	} else {
		l.currentStep = 0
	}
	l.mu.Unlock()
	log.Printf("[%s] advanceStep to=%d", l.name, l.currentStep)
	l.showInstruction()
}

// completeFlow adalah fungsi untuk complete flow.
func (l *LineState) completeFlow() {
	l.mu.Lock()
	serial := l.serialValue
	fg := l.fgTypeValue
	spc := l.spcValue
	compType := l.compTypeValue
	compCode := l.compCodeValue
	l.mu.Unlock()

	log.Printf("[%s] completeFlow start serial=%q fg=%q spc=%q compType=%q compCode=%q", l.name, serial, fg, spc, compType, compCode)
	if serial == "" {
		l.flashStatus("S/N Input Empty", color.RGBA{R: 0xFF, G: 0x45, B: 0x00, A: 0xFF}, 1500*time.Millisecond)
		return
	}

	if l.customSaveFlow != nil {
		if err := l.customSaveFlow(); err != nil {
			l.flashStatus(err.Error(), color.RGBA{R: 0xFF, G: 0x45, B: 0x00, A: 0xFF}, 2000*time.Millisecond)
			return
		}
	}

	log.Printf("[%s] completeFlow serial=%s fg=%s spc=%s compType=%s compCode=%s", l.name, serial, fg, spc, compType, compCode)
	l.history.append(historyRow{
		Date:   time.Now().Format("02-01-2006"),
		Time:   time.Now().Format("15:04:05"),
		Serial: serial,
		FgType: fg,
	})
	l.mu.Lock()
	focusSelf := !l.focusPeerAfterComplete
	l.resetLocked(false, focusSelf)
	l.mu.Unlock()
	msg := l.successMessage
	if msg == "" {
		msg = "Data tersimpan"
	}
	l.flashStatus(msg, color.RGBA{R: 0x00, G: 0xBF, B: 0x8F, A: 0xFF}, 1200*time.Millisecond)
	if !focusSelf && l.otherFocus != nil {
		runOnMain(l.otherFocus)
	}
}

// showInstruction adalah fungsi untuk show instruction.
func (l *LineState) showInstruction() {
	log.Printf("[%s] showInstruction step=%d label=%q", l.name, l.currentStep, l.steps[l.currentStep].label)
	step := l.steps[l.currentStep]
	l.updateJumboText(l.buildJumboLines(step.label), color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF})
}

// flashStatus adalah fungsi untuk flash status.
func (l *LineState) flashStatus(msg string, clr color.Color, revertDuration time.Duration) {
	log.Printf("[%s] flashStatus msg=%q revert=%s", l.name, msg, revertDuration)
	l.updateJumboText(l.buildJumboLines(msg), clr)
	if revertDuration > 0 {
		currentStep := l.currentStep
		time.AfterFunc(revertDuration, func() {
			runOnMain(func() {
				l.mu.Lock()
				defer l.mu.Unlock()
				if currentStep == l.currentStep {
					l.showInstruction()
				}
			})
		})
	}
}

// Reset adalah fungsi untuk reset.
func (l *LineState) Reset(showInfo bool) {
	log.Printf("[%s] Reset called showInfo=%v", l.name, showInfo)
	l.mu.Lock()
	defer l.mu.Unlock()
	l.resetLocked(showInfo, true)
}

// resetLocked adalah fungsi untuk reset locked.
func (l *LineState) resetLocked(showInfo bool, focusSelf bool) {
	l.currentStep = 0
	l.testRunning = false
	l.spcValue = ""
	l.serialValue = ""
	l.fgTypeValue = ""
	l.compTypeValue = ""
	l.compCodeValue = ""

	if l.spcLabel != nil {
		l.spcLabel.SetText("-")
	}
	if l.snLabel != nil {
		l.snLabel.SetText("-")
	}
	if l.fgLabel != nil {
		l.fgLabel.SetText("-")
	}
	if l.compTypeLabel != nil {
		l.compTypeLabel.SetText("-")
	}
	if l.compCodeLabel != nil {
		l.compCodeLabel.SetText("-")
	}
	l.showInstruction()
	if showInfo {
		l.flashStatus("Sequence direset", color.RGBA{R: 0x1E, G: 0x90, B: 0xFF, A: 0xFF}, 1200*time.Millisecond)
	}
	if focusSelf {
		l.refocusInput()
	}
}

// handleKeyword adalah fungsi untuk menangani keyword.
func (l *LineState) handleKeyword(text string, skipRefocus *bool) bool {
	log.Printf("[%s] handleKeyword text=%q", l.name, text)
	cmd := strings.ToUpper(text)
	switch cmd {
	case "REBOOT1234":
		dialog.ShowInformation("Reboot", "Perintah REBOOT diterima. Silakan konfirmasi di panel utama.", l.window)
		runSystemCommandAsync(l.window, system.Reboot)
		log.Printf("[%s] keyword handled REBOOT", l.name)
		return true
	case "DOWN123456":
		dialog.ShowInformation("Shutdown", "Perintah shutdown diterima.", l.window)
		runSystemCommandAsync(l.window, system.Shutdown)
		log.Printf("[%s] keyword handled SHUTDOWN", l.name)
		return true
	case "LEFT123456":
		if l.otherFocus != nil {
			l.otherFocus()
			// l.flashStatus("Pindah ke line kiri", color.RGBA{R: 0x1E, G: 0x90, B: 0xFF, A: 0xFF}, 1000*time.Millisecond)
		}
		log.Printf("[%s] keyword handled LEFT", l.name)
		if skipRefocus != nil {
			*skipRefocus = true
		}
		return true
	case "RIGHT12345":
		if l.otherFocus != nil {
			l.otherFocus()
			// l.flashStatus("Pindah ke line kanan", color.RGBA{R: 0x1E, G: 0x90, B: 0xFF, A: 0xFF}, 1000*time.Millisecond)
		}
		log.Printf("[%s] keyword handled RIGHT", l.name)
		if skipRefocus != nil {
			*skipRefocus = true
		}
		return true
	case "RESET12345":
		l.Reset(true)
		log.Printf("[%s] keyword handled RESET", l.name)
		return true
	case "MAINTENANCE":
		if navigation.NavigateToRoute("maintenance") {
			// l.flashStatus("Halaman Maintenance dibuka", color.RGBA{R: 0x00, G: 0xBF, B: 0x8F, A: 0xFF}, 1200*time.Millisecond)
		} else {
			dialog.ShowInformation("Maintenance", "Membuka halaman Maintenance (simulasi).", l.window)
		}
		log.Printf("[%s] keyword handled MAINTENANCE", l.name)
		return true
	case "SIMULATEALL":
		l.applySimulationMode("simulateAll")
		return true
	case "SIMULATEDB":
		l.applySimulationMode("simulateDB")
		return true
	case "SIMULATEHW":
		l.applySimulationMode("simulateHW")
		return true
	case "SIMULATELIVE":
		l.applySimulationMode("LIVE")
		return true
	default:
		return false
	}
}

// applySimulationMode adalah fungsi untuk apply simulation mode.
func (l *LineState) applySimulationMode(mode string) {
	show := getMaintenanceState().renderMode(mode)
	applyMode(mode, l.window)
	// l.flashStatus("Mode diubah ke "+show, color.RGBA{R: 0x00, G: 0xBF, B: 0x8F, A: 0xFF}, 1200*time.Millisecond)
	log.Printf("[%s] keyword handled MODE=%s", l.name, show)
}

// setPeerFocus adalah fungsi untuk mengatur peer focus.
func (l *LineState) SetPeerFocus(fn func(), autoSwitch bool) {
	l.otherFocus = fn
	l.focusPeerAfterComplete = autoSwitch
}

// focus adalah fungsi untuk focus.
func (l *LineState) Focus() {
	if l.inputEntry == nil || l.window == nil {
		log.Printf("[%s] focus skipped inputEntry or window nil", l.name)
		return
	}
	canv := l.window.Canvas()
	if canv == nil {
		log.Printf("[%s] focus skipped canvas nil", l.name)
		return
	}
	if canv.Content() == nil {
		log.Printf("[%s] focus skipped canvas content nil", l.name)
		return
	}
	if drv := fyne.CurrentApp().Driver(); drv != nil {
		if drv.CanvasForObject(l.inputEntry) == nil {
			log.Printf("[%s] focus skipped inputEntry not on canvas (retry %d)", l.name, l.focusRetries)
			if l.focusRetries < 5 {
				l.focusRetries++
				time.AfterFunc(200*time.Millisecond, func() {
					l.refocusInput()
				})
			}
			return
		}
	}
	l.focusRetries = 0
	log.Printf("[%s] focus attempt", l.name)
	l.window.Canvas().Focus(l.inputEntry)
}

// refocusInput adalah fungsi untuk refocus input.
func (l *LineState) refocusInput() {
	log.Printf("[%s] refocusInput request", l.name)
	runOnMain(func() {
		l.Focus()
	})
}

// setLastInput adalah fungsi untuk mengatur last input.
func (l *LineState) setLastInput(val string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.lastInput != nil {
		l.lastInput.SetText(val)
	}
}

// setSerialValidator adalah fungsi untuk mengatur serial validator.
func (l *LineState) SetSerialValidator(fn func(string) (string, error)) {
	l.customSerialValidator = fn
}

// SetStepValidator menambahkan validator sebelum nilai diset untuk langkah tertentu.
func (l *LineState) SetStepValidator(kind StepKind, fn func(string) error) {
	if l.stepValidators == nil {
		l.stepValidators = make(map[lineStepKind]func(string) error)
	}
	l.stepValidators[lineStepKind(kind)] = fn
}

// setSaveHandler adalah fungsi untuk mengatur menyimpan handler.
func (l *LineState) SetSaveHandler(fn func() error) {
	l.customSaveFlow = fn
}

// setSuccessMessage adalah fungsi untuk mengatur success message.
func (l *LineState) SetSuccessMessage(msg string) {
	l.successMessage = msg
}

// values adalah fungsi untuk values.
func (l *LineState) Values() (serial, fg, spc, compType, compCode string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.serialValue, l.fgTypeValue, l.spcValue, l.compTypeValue, l.compCodeValue
}

// DeriveFgType adalah fungsi untuk derive fg type.
func DeriveFgType(sn string) string {
	if sn == "" {
		return "-"
	}
	if len(sn) >= 4 {
		return fmt.Sprintf("FG-%s", sn[:4])
	}
	return fmt.Sprintf("FG-%s", sn)
}

// deriveCompressorTitle adalah fungsi untuk derive kompresor title.
func deriveCompressorTitle(code string) string {
	if code == "" {
		return "-"
	}
	if len(code) >= 2 {
		return fmt.Sprintf("COMP %s", code[:2])
	}
	return fmt.Sprintf("COMP %s", code)
}

// runOnMain adalah fungsi untuk menjalankan on utama.
func runOnMain(fn func()) {
	if fn == nil {
		return
	}
	if fyne.CurrentApp() != nil {
		fyne.Do(fn)
		return
	}
	fn()
}

// runSystemCommandAsync adalah fungsi untuk menjalankan system command async.
func runSystemCommandAsync(w fyne.Window, action func() error) {
	if action == nil {
		return
	}
	go func() {
		if err := action(); err != nil {
			fyne.Do(func() {
				dialog.ShowError(err, w)
			})
		}
	}()
}

// NewResponsiveLineLayout arranges line cards so they wrap responsively on small screens.
func NewResponsiveLineLayout(objects ...fyne.CanvasObject) fyne.CanvasObject {
	config := make([]fyne.CanvasObject, 0, len(objects))
	for _, obj := range objects {
		if obj == nil {
			continue
		}
		config = append(config, xlayout.Responsive(container.NewPadded(obj), 1, 0.5, 0.5, 0.5))
	}
	return xlayout.NewResponsiveLayout(config...)
}
