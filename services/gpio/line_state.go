/*
    file:           services/gpio/line_state.go
    description:    Driver GPIO untuk line state
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package gpio

import (
	"strconv"
	"strings"
	"sync"

	"go-ecb/configs"
)

var (
	lineActiveMu    sync.RWMutex
	lineActiveValue int
	lineActiveInit  bool
)

// getLineActive adalah fungsi untuk mengambil jalur active.
func getLineActive() int {
	lineActiveMu.RLock()
	defer lineActiveMu.RUnlock()
	return lineActiveValue
}

// setLineActive adalah fungsi untuk mengatur jalur active.
func setLineActive(line int) {
	normalized := clampLine(line)
	lineActiveMu.Lock()
	defer lineActiveMu.Unlock()
	if lineActiveInit && lineActiveValue == normalized {
		return
	}
	lineActiveInit = true
	lineActiveValue = normalized
	if shouldInteractWithHardware() {
		layout := GetPinLayout()
		writeLevel(layout.LineSelect, levelFromLine(normalized))
	}
}

// toggleLineValue adalah fungsi untuk mengubah jalur value.
func toggleLineValue(line int) int {
	if line == 1 {
		return 0
	}
	return 1
}

// parseLineValue adalah fungsi untuk parse jalur value.
func parseLineValue(value string, fallback int) int {
	if strings.TrimSpace(value) == "" {
		return clampLine(fallback)
	}
	line, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return clampLine(fallback)
	}
	return clampLine(line)
}

// clampLine adalah fungsi untuk clamp jalur.
func clampLine(line int) int {
	if line <= 0 {
		return 0
	}
	if line >= 1 {
		return 1
	}
	return line
}

// levelFromLine adalah fungsi untuk level from jalur.
func levelFromLine(line int) Level {
	if line == 1 {
		return LevelHigh
	}
	return LevelLow
}

// shouldInteractWithHardware adalah fungsi untuk should interact with hardware.
func shouldInteractWithHardware() bool {
	mode := configs.LoadSimoConfig().EcbMode
	return mode != "simulateAll" && mode != "simulateHW"
}

// GetLineActive adalah fungsi untuk mengambil jalur active.
func GetLineActive() int {
	return getLineActive()
}

// SetLineActive adalah fungsi untuk mengatur jalur active.
func SetLineActive(line int) {
	setLineActive(line)
}

// ToggleLineActive adalah fungsi untuk mengubah jalur active.
func ToggleLineActive() int {
	next := toggleLineValue(getLineActive())
	setLineActive(next)
	return next
}
