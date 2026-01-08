/*
    file:           services/gpio/startup.go
    description:    Driver GPIO untuk startup
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package gpio

import (
	"strings"
	"time"

	"go-ecb/configs"
)

func InitializeHardware() {
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
