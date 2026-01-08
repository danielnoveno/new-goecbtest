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

func StartTest(line *int) {
	if line != nil {
		SetLineActive(clampLine(*line))
	}
	if !shouldInteractWithHardware() {
		return
	}
	layout := GetPinLayout()
	// Turn OFF Reset LED, Turn ON Start LED
	writeLevel(layout.Reset, LevelLow)
	writeLevel(layout.ResetAlt, LevelHigh)
	
	writeLevel(layout.Start, LevelHigh)
	writeLevel(layout.StartAlt, LevelLow)
}

func ResetTest(line *int) {
	if line != nil {
		SetLineActive(clampLine(*line))
	}
	if !shouldInteractWithHardware() {
		return
	}
	layout := GetPinLayout()
	// Turn OFF Start LED, Turn ON Reset LED
	writeLevel(layout.Start, LevelLow)
	writeLevel(layout.StartAlt, LevelHigh)

	writeLevel(layout.Reset, LevelHigh)
	writeLevel(layout.ResetAlt, LevelLow)
}

func LineToggle() int {
	return ToggleLineActive()
}

func LineSet(line int) int {
	SetLineActive(line)
	return GetLineActive()
}

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

func readHardwareState() string {
	layout := GetPinLayout()
	pass := readLevel(layout.Pass)
	fail := readLevel(layout.Fail)
	undertest := readLevel(layout.UnderTest)
	line := readLevel(layout.LineSelect)
	return fmt.Sprintf("%s.%s.%s.%s", undertest, pass, fail, line)
}
