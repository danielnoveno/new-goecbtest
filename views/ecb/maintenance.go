/*
   file:           views/ecb/maintenance.go
   description:    Layar ECB untuk maintenance
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package ecb

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-ecb/configs"
	"go-ecb/services/gpio"
	maintenanceSvc "go-ecb/services/maintenance"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type maintenanceState struct {
	mode       binding.String
	pesan      binding.String
	hasil      binding.String
	pass       binding.String
	fail       binding.String
	undertest  binding.String
	lineSelect binding.String
	connStatus binding.String

	lineCount  int
	lineActive int

	mu     sync.Mutex
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

var (
	maintenanceOnce   sync.Once
	sharedState       *maintenanceState
	modeChangeHandler func(string)
)

func RegisterModeChangeHandler(fn func(string)) {
	modeChangeHandler = fn
}

func getMaintenanceState() *maintenanceState {
	maintenanceOnce.Do(func() {
		simoCfg := configs.LoadSimoConfig()
		sharedState = &maintenanceState{
			mode:       binding.NewString(),
			pesan:      binding.NewString(),
			hasil:      binding.NewString(),
			pass:       binding.NewString(),
			fail:       binding.NewString(),
			undertest:  binding.NewString(),
			lineSelect: binding.NewString(),
			connStatus: binding.NewString(),
			lineCount:  deriveLineCount(simoCfg.EcbLineType),
		}
		_ = sharedState.mode.Set("LIVE")
		_ = sharedState.pesan.Set("Menunggu perintah...")
		_ = sharedState.hasil.Set("IDLE")
		_ = sharedState.pass.Set("-")
		_ = sharedState.fail.Set("-")
		_ = sharedState.undertest.Set("-")
		_ = sharedState.lineSelect.Set("-")
		_ = sharedState.connStatus.Set("Polling: idle")
	})
	return sharedState
}

func deriveLineCount(lineType string) int {
	switch strings.ToLower(strings.TrimSpace(lineType)) {
	case "sn-only-single", "refrig-single", "refrig-po-single":
		return 1
	default:
		return 2
	}
}

// func (s *maintenanceState) adminerURL() (string, error) {
// 	return adminer.EnsureRunning()
// }

func (s *maintenanceState) setString(target binding.String, value string) {
	if target == nil {
		return
	}
	if err := target.Set(value); err != nil {
		log.Printf("[maintenance] failed to set binding: %v", err)
	}
}

func (s *maintenanceState) setMode(mode string, onModeChange func(string), w fyne.Window) {
	mode = strings.TrimSpace(mode)
	if mode == "" {
		mode = "LIVE"
	}
	_ = os.Setenv("ECB_MODE", mode)
	_ = os.Setenv("SIMO_ECBMODE", mode)
	s.mu.Lock()
	s.lineActive = 0
	s.mu.Unlock()

	display := s.renderMode(mode)
	s.setString(s.mode, display)
	s.setString(s.pesan, "Mode diubah ke "+display)
	s.setString(s.connStatus, "Polling: menyambung ulang...")
	s.restartWatcher()
	if onModeChange != nil {
		onModeChange(mode)
	}
	if w != nil {
		w.Canvas().Refresh(w.Content())
	}
}

func (s *maintenanceState) renderMode(mode string) string {
	switch strings.ToLower(mode) {
	case "live":
		return "LIVE"
	case "simulatehw":
		return "Simulate H/W"
	case "simulatedb":
		return "Simulate DB"
	case "simulateall":
		return "Simulate H/W+DB"
	default:
		return strings.ToUpper(mode)
	}
}

func (s *maintenanceState) restartWatcher() {
	s.mu.Lock()
	if s.cancel != nil {
		s.cancel()
	}
	s.mu.Unlock()
	s.wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	s.cancel = cancel
	s.mu.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.watchLocal(ctx)
	}()
}

func (s *maintenanceState) watchLocal(ctx context.Context) {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	s.setString(s.connStatus, "Local polling aktif")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			state := gpio.ReadLocalEcbState()
			s.applyStateString(state)
			s.setString(s.pesan, "status: OK")
		}
	}
}

func (s *maintenanceState) applyStateString(value string) {
	parts := strings.Split(value, ".")
	if len(parts) < 4 {
		return
	}
	s.setString(s.undertest, parts[0])
	s.setString(s.pass, parts[1])
	s.setString(s.fail, parts[2])
	s.setString(s.lineSelect, parts[3])
	if line, err := strconv.Atoi(parts[3]); err == nil {
		s.mu.Lock()
		s.lineActive = clampLine(line)
		s.mu.Unlock()
	}
	result := deriveHasil(value)
	s.setString(s.hasil, result)
	if result == "RUSAK" {
		s.setString(s.pesan, fmt.Sprintf("Logika Error: Kombinasi Pin %s tidak valid", value))
	} else {
		s.setString(s.pesan, "Status: OK")
	}
}

func deriveHasil(state string) string {
	if strings.HasPrefix(state, "1.1.1") {
		return "IDLE"
	}
	if strings.HasPrefix(state, "1.0.1") {
		return "PASS"
	}
	if strings.HasPrefix(state, "1.1.0") {
		return "FAIL"
	}
	if strings.HasPrefix(state, "0.1.1") {
		return "UNDERTEST"
	}
	return "RUSAK"
}

func clampLine(line int) int {
	if line <= 0 {
		return 0
	}
	if line >= 1 {
		return 1
	}
	return line
}

func (s *maintenanceState) runCommand(action func(), successMessage string) {
	go func() {
		action()
		if successMessage != "" {
			s.setString(s.pesan, successMessage)
		}
	}()
}

func applyMode(mode string, w fyne.Window) {
	getMaintenanceState().setMode(mode, modeChangeHandler, w)
}

func MaintenanceUI(onModeChange func(string), currentMode string, w fyne.Window, pinController maintenanceSvc.PinController) fyne.CanvasObject {
	state := getMaintenanceState()

	if strings.TrimSpace(currentMode) != "" {
		state.setMode(currentMode, onModeChange, w)
	} else {
		state.restartWatcher()
	}

	pinLayout := gpio.DefaultPinLayout()
	if pinController != nil {
		pinLayout = pinController.GetPins()
	}
	pinForm := newPinConfigForm(pinLayout)

	modeLabel := widget.NewLabelWithData(state.mode)
	modeLabel.TextStyle = fyne.TextStyle{Bold: true}

	statusLabel := widget.NewLabelWithData(state.connStatus)
	statusLabel.Wrapping = fyne.TextWrapWord

	pesanLabel := widget.NewLabelWithData(state.pesan)
	pesanLabel.Wrapping = fyne.TextWrapWord
	hasilLabel := widget.NewLabelWithData(state.hasil)
	hasilLabel.TextStyle = fyne.TextStyle{Bold: true}

	passLabel := widget.NewLabelWithData(state.pass)
	failLabel := widget.NewLabelWithData(state.fail)
	undertestLabel := widget.NewLabelWithData(state.undertest)
	lineLabel := widget.NewLabelWithData(state.lineSelect)

	btnLive := widget.NewButton("LIVE", func() {
		state.setMode("LIVE", onModeChange, w)
	})
	btnSimHW := widget.NewButton("Simulate H/W", func() {
		state.setMode("simulateHW", onModeChange, w)
	})
	btnSimDB := widget.NewButton("Simulate DB", func() {
		state.setMode("simulateDB", onModeChange, w)
	})
	btnSimAll := widget.NewButton("Simulate H/W+DB", func() {
		state.setMode("simulateAll", onModeChange, w)
	})

	btnInit := widget.NewButton("RE-INIT", func() {
		state.runCommand(gpio.InitializeControl, "Sistem diinisialisasi ulang")
	})
	btnStart := widget.NewButton("START", func() {
		state.runCommand(func() { gpio.StartTest(nil) }, "Proses pengujian dimulai")
		state.setString(state.hasil, "UNDERTEST")
	})
	btnReset := widget.NewButton("RESET", func() {
		state.runCommand(func() { gpio.ResetTest(nil) }, "Sistem direset")
	})

	btnLineSelect := widget.NewButton("Line Select", func() {
		state.runCommand(func() {
			line := gpio.LineToggle()
			state.setString(state.lineSelect, fmt.Sprintf("%d", line))
		}, "Line berganti")
	})
	if state.lineCount < 2 {
		btnLineSelect.Disable()
	}

	// btnDBManager := widget.NewButton("Database Manager", func() {
	// 	go func() {
	// 		target, startErr := state.adminerURL()
	// 		if startErr != nil {
	// 			log.Printf("gagal menjalankan Adminer: %v", startErr)
	// 			if w != nil {
	// 				dialog.ShowError(fmt.Errorf("Adminer gagal dijalankan: %w", startErr), w)
	// 			}
	// 			return
	// 		}
	// 		if _, err := url.Parse(target); err != nil {
	// 			dialog.ShowError(fmt.Errorf("URL Adminer tidak valid: %w", err), w)
	// 			return
	// 		}
	// 		if err := browser.OpenURL(target); err != nil {
	// 			log.Printf("gagal membuka Adminer: %v", err)
	// 			if w != nil {
	// 				dialog.ShowError(err, w)
	// 			}
	// 		}
	// 	}()
	// })

	modeSection := container.NewVBox(
		widget.NewLabelWithStyle("Application Mode", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		modeLabel,
		container.NewGridWithColumns(4, btnLive, btnSimHW, btnSimDB, btnSimAll),
		statusLabel,
	)

	commandButtons := []fyne.CanvasObject{btnInit, btnStart, btnReset}
	if state.lineCount > 1 {
		commandButtons = append(commandButtons, btnLineSelect)
	}

	// commandButtons = append(commandButtons, btnDBManager)

	cmdSection := container.NewVBox(
		widget.NewLabelWithStyle("ECB Commands", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(len(commandButtons), commandButtons...),
	)

	statusSection := container.NewVBox(
		widget.NewLabelWithStyle("Pesan", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		pesanLabel,
		widget.NewLabelWithStyle("Hasil", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		hasilLabel,
	)

	underTestPinValue := widget.NewLabel(pinLayout.UnderTest)
	passPinValue := widget.NewLabel(pinLayout.Pass)
	failPinValue := widget.NewLabel(pinLayout.Fail)
	linePinValue := widget.NewLabel(pinLayout.LineSelect)

	headerRow := container.NewGridWithColumns(4,
		widget.NewLabelWithStyle("Nama", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Pin Fisik (ref)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("WiringPi", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Nilai", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	underTestRefLabel := widget.NewLabel(physicalReference(pinLayout.UnderTest))
	passRefLabel := widget.NewLabel(physicalReference(pinLayout.Pass))
	failRefLabel := widget.NewLabel(physicalReference(pinLayout.Fail))
	lineRefLabel := widget.NewLabel(physicalReference(pinLayout.LineSelect))

	pinGrid := container.NewVBox(
		headerRow,
		container.NewGridWithColumns(4,
			widget.NewLabel("UNDERTEST"), underTestRefLabel, underTestPinValue, undertestLabel),
		container.NewGridWithColumns(4,
			widget.NewLabel("PASS"), passRefLabel, passPinValue, passLabel),
		container.NewGridWithColumns(4,
			widget.NewLabel("FAIL"), failRefLabel, failPinValue, failLabel),
		container.NewGridWithColumns(4,
			widget.NewLabel("LINESELECT"), lineRefLabel, linePinValue, lineLabel),
	)

	configInfo := "Ubah wiringPi sesuai kebutuhan; referensi pin fisik ada di tabel di atas."
	if pinController == nil {
		configInfo = "Konfigurasi pin belum tersedia karena layanan tidak diinisialisasi."
	}
	makeFormRow := func(label string, entry *widget.Entry) fyne.CanvasObject {
		return container.NewGridWithColumns(2,
			widget.NewLabel(label),
			entry,
		)
	}

	saveButton := widget.NewButton("Simpan konfigurasi pin", func() {
		if pinController == nil {
			return
		}
		newLayout := pinForm.toLayout()
		if err := pinController.SavePins(newLayout); err != nil {
			dialog.ShowError(err, w)
			return
		}
		dialog.ShowInformation("Konfigurasi pin tersimpan", "Perubahan tersimpan.", w)
		updated := pinController.GetPins()
		underTestPinValue.SetText(updated.UnderTest)
		passPinValue.SetText(updated.Pass)
		failPinValue.SetText(updated.Fail)
		linePinValue.SetText(updated.LineSelect)
		
		underTestRefLabel.SetText(physicalReference(updated.UnderTest))
		passRefLabel.SetText(physicalReference(updated.Pass))
		failRefLabel.SetText(physicalReference(updated.Fail))
		lineRefLabel.SetText(physicalReference(updated.LineSelect))
		
		pinForm.update(updated)
	})

	resetButton := widget.NewButton("Kembalikan default", func() {
		defaultLayout := gpio.DefaultPinLayout()
		pinForm.update(defaultLayout)
	})
	if pinController == nil {
		saveButton.Disable()
		resetButton.Disable()
	}

	pinConfigSection := widget.NewCard(
		"Konfigurasi Pin WiringPi",
		"Ubah wiringPi agar sesuai wiring actual. Tombol RE-INIT akan mengirim ulang konfigurasi ke hardware.",
		container.NewVBox(
			widget.NewLabel(configInfo),
			container.NewVBox(
				makeFormRow("UNDERTEST (WiringPi)", pinForm.underTestEntry),
				makeFormRow("FAIL (WiringPi)", pinForm.failEntry),
				makeFormRow("PASS (WiringPi)", pinForm.passEntry),
				makeFormRow("LINE SELECT (WiringPi)", pinForm.lineSelectEntry),
				// makeFormRow("RESET (WiringPi)", pinForm.resetEntry),
				// makeFormRow("RESET (alt) (WiringPi)", pinForm.resetAltEntry),
				// makeFormRow("START (WiringPi)", pinForm.startEntry),
				// makeFormRow("START (alt) (WiringPi)", pinForm.startAltEntry),
			),
			container.NewHBox(saveButton, resetButton),
		),
	)

	content := container.NewVBox(
		widget.NewLabelWithStyle("ECB Maintenance Panel", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		modeSection,
		widget.NewSeparator(),
		cmdSection,
		widget.NewSeparator(),
		statusSection,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Status Pin", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		pinGrid,
		widget.NewSeparator(),
		pinConfigSection,
	)

	return container.NewVScroll(container.NewPadded(content))
}

func physicalReference(wiringPiPin string) string {
	return gpio.WiringPiToPhysical(wiringPiPin)
}

type pinConfigForm struct {
	underTestEntry  *widget.Entry
	failEntry       *widget.Entry
	passEntry       *widget.Entry
	lineSelectEntry *widget.Entry
	resetEntry      *widget.Entry
	resetAltEntry   *widget.Entry
	startEntry      *widget.Entry
	startAltEntry   *widget.Entry
}

func newPinConfigForm(layout gpio.PinLayout) *pinConfigForm {
	return &pinConfigForm{
		underTestEntry:  entryWithValue(layout.UnderTest),
		failEntry:       entryWithValue(layout.Fail),
		passEntry:       entryWithValue(layout.Pass),
		lineSelectEntry: entryWithValue(layout.LineSelect),
		resetEntry:      entryWithValue(layout.Reset),
		resetAltEntry:   entryWithValue(layout.ResetAlt),
		startEntry:      entryWithValue(layout.Start),
		startAltEntry:   entryWithValue(layout.StartAlt),
	}
}

func entryWithValue(value string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText(value)
	return entry
}

func (f *pinConfigForm) update(layout gpio.PinLayout) {
	if f == nil {
		return
	}
	f.underTestEntry.SetText(layout.UnderTest)
	f.failEntry.SetText(layout.Fail)
	f.passEntry.SetText(layout.Pass)
	f.lineSelectEntry.SetText(layout.LineSelect)
	f.resetEntry.SetText(layout.Reset)
	f.resetAltEntry.SetText(layout.ResetAlt)
	f.startEntry.SetText(layout.Start)
	f.startAltEntry.SetText(layout.StartAlt)
}

func (f *pinConfigForm) toLayout() gpio.PinLayout {
	if f == nil {
		return gpio.DefaultPinLayout()
	}
	return gpio.PinLayout{
		UnderTest:  strings.TrimSpace(f.underTestEntry.Text),
		Fail:       strings.TrimSpace(f.failEntry.Text),
		Pass:       strings.TrimSpace(f.passEntry.Text),
		LineSelect: strings.TrimSpace(f.lineSelectEntry.Text),
		Reset:      strings.TrimSpace(f.resetEntry.Text),
		ResetAlt:   strings.TrimSpace(f.resetAltEntry.Text),
		Start:      strings.TrimSpace(f.startEntry.Text),
		StartAlt:   strings.TrimSpace(f.startAltEntry.Text),
	}
}
