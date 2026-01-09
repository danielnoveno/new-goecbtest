/*
   file:           services/gpio/pins.go
   description:    Driver GPIO untuk pins
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package gpio

import (
	"strings"
	"sync"
)

type PinLayout struct {
	UnderTest  string
	Fail       string
	Pass       string
	LineSelect string
	Start      string
	StartAlt   string
	Reset      string
	ResetAlt   string
}

var (
	pinLayoutMu    sync.RWMutex
	pinLayoutValue = DefaultPinLayout()
)

func DefaultPinLayout() PinLayout {
	return PinLayout{
		UnderTest:  "22",
		Fail:       "21",
		Pass:       "2",
		LineSelect: "25",
		Start:      "24",
		StartAlt:   "29",
		Reset:      "23",
		ResetAlt:   "28",
	}
}

func NormalizePinLayout(layout PinLayout) PinLayout {
	def := DefaultPinLayout()
	layout.UnderTest = normalizeOrDefault(layout.UnderTest, def.UnderTest)
	layout.Fail = normalizeOrDefault(layout.Fail, def.Fail)
	layout.Pass = normalizeOrDefault(layout.Pass, def.Pass)
	layout.LineSelect = normalizeOrDefault(layout.LineSelect, def.LineSelect)
	layout.Start = normalizeOrDefault(layout.Start, def.Start)
	layout.StartAlt = normalizeOrDefault(layout.StartAlt, def.StartAlt)
	layout.Reset = normalizeOrDefault(layout.Reset, def.Reset)
	layout.ResetAlt = normalizeOrDefault(layout.ResetAlt, def.ResetAlt)
	return layout
}

func normalizeOrDefault(value, fallback string) string {
	if trimmed := strings.TrimSpace(value); trimmed != "" {
		return trimmed
	}
	return fallback
}

func SetPinLayout(layout PinLayout) {
	pinLayoutMu.Lock()
	pinLayoutValue = NormalizePinLayout(layout)
	pinLayoutMu.Unlock()
}

func GetPinLayout() PinLayout {
	pinLayoutMu.RLock()
	defer pinLayoutMu.RUnlock()
	return pinLayoutValue
}
