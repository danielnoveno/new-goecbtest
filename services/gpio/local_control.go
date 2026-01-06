/*
    file:           services/gpio/local_control.go
    description:    Driver GPIO untuk local control
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package gpio

import (
	"fmt"
	"strings"
	"time"

	"go-ecb/configs"
)

// InitializeControl adalah fungsi untuk initialize control.
func InitializeControl() {
	simoCfg := configs.LoadSimoConfig()
	ecbMode := simoCfg.EcbMode
	if ecbMode == "simulateAll" || ecbMode == "simulateHW" {
		return
	}

	layout := GetPinLayout()
	for _, pin := range []struct {
		name string
		mode PinMode
	}{
		{layout.Reset, ModeOutput},
		{layout.ResetAlt, ModeOutput},
		{layout.Start, ModeOutput},
		{layout.StartAlt, ModeOutput},
		{layout.LineSelect, ModeOutput},
		{layout.Pass, ModeInput},
		{layout.Fail, ModeInput},
		{layout.UnderTest, ModeInput},
	} {
		if strings.TrimSpace(pin.name) == "" {
			continue
		}
		configureMode(pin.name, pin.mode)
	}

	writeLevel(layout.Pass, LevelHigh)
	writeLevel(layout.Fail, LevelHigh)
	writeLevel(layout.UnderTest, LevelHigh)
	writeLevel(layout.ResetAlt, LevelHigh)
	writeLevel(layout.StartAlt, LevelHigh)

	for i := 0; i < 3; i++ {
		writeLevel(layout.Reset, LevelHigh)
		writeLevel(layout.ResetAlt, LevelLow)
		time.Sleep(20 * time.Millisecond)
		writeLevel(layout.Reset, LevelLow)
		writeLevel(layout.ResetAlt, LevelHigh)
		time.Sleep(500 * time.Millisecond)
	}
}

// StartTest adalah fungsi untuk menjalankan test.
func StartTest(line *int) {
	if line != nil {
		SetLineActive(clampLine(*line))
	}
	if !shouldInteractWithHardware() {
		return
	}
	layout := GetPinLayout()
	writeLevel(layout.Reset, LevelHigh)
	writeLevel(layout.ResetAlt, LevelLow)
	time.Sleep(20 * time.Millisecond)
	writeLevel(layout.Reset, LevelLow)
	writeLevel(layout.ResetAlt, LevelHigh)
	time.Sleep(500 * time.Millisecond)
	writeLevel(layout.Start, LevelHigh)
	writeLevel(layout.StartAlt, LevelLow)
	time.Sleep(100 * time.Millisecond)
	writeLevel(layout.Start, LevelLow)
	writeLevel(layout.StartAlt, LevelHigh)
}

// ResetTest adalah fungsi untuk reset test.
func ResetTest(line *int) {
	if line != nil {
		SetLineActive(clampLine(*line))
	}
	if !shouldInteractWithHardware() {
		return
	}
	layout := GetPinLayout()
	for i := 0; i < 2; i++ {
		writeLevel(layout.Reset, LevelHigh)
		writeLevel(layout.ResetAlt, LevelLow)
		time.Sleep(20 * time.Millisecond)
		writeLevel(layout.Reset, LevelLow)
		writeLevel(layout.ResetAlt, LevelHigh)
		time.Sleep(500 * time.Millisecond)
	}
}

// LineToggle adalah fungsi untuk jalur mengubah.
func LineToggle() int {
	return ToggleLineActive()
}

// LineSet adalah fungsi untuk jalur mengatur.
func LineSet(line int) int {
	SetLineActive(line)
	return GetLineActive()
}

// ReadLocalEcbState adalah fungsi untuk membaca local ecb status.
func ReadLocalEcbState() string {
	simoCfg := configs.LoadSimoConfig()
	mode := simoCfg.EcbMode
	switch mode {
	case "simulateAll", "simulateHW":
		return simoCfg.EcbStateDefault
	default:
		return readHardwareState()
	}
}

// readHardwareState adalah fungsi untuk membaca hardware status.
func readHardwareState() string {
	layout := GetPinLayout()
	pass := readLevel(layout.Pass)
	fail := readLevel(layout.Fail)
	undertest := readLevel(layout.UnderTest)
	line := readLevel(layout.LineSelect)
	return fmt.Sprintf("%s.%s.%s.%s", pass, fail, undertest, line)
}
