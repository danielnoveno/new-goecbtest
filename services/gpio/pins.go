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

// PinLayout menjelaskan penugasan pin wiringPi yang digunakan UI maintenance
// ECB dan driver GPIO. Struktur ini memungkinkan nilai di-load ulang atau
// dioverride saat aplikasi berjalan.
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

// DefaultPinLayout mengembalikan nilai default wiringPi sesuai tabel referensi
// agar penyebaran baru memiliki perilaku yang sama.
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

// NormalizePinLayout ensures every field has a non-empty value by falling back
// to the defaults and trimming whitespace.
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

// normalizeOrDefault adalah fungsi untuk normalize or default.
func normalizeOrDefault(value, fallback string) string {
	if trimmed := strings.TrimSpace(value); trimmed != "" {
		return trimmed
	}
	return fallback
}

// SetPinLayout menimpa layout saat ini dan menyimpan versi ter-normalisasi agar
// helper GPIO lain bisa membacanya dengan aman.
func SetPinLayout(layout PinLayout) {
	pinLayoutMu.Lock()
	pinLayoutValue = NormalizePinLayout(layout)
	pinLayoutMu.Unlock()
}

// GetPinLayout mengembalikan snapshot dari layout aktif.
func GetPinLayout() PinLayout {
	pinLayoutMu.RLock()
	defer pinLayoutMu.RUnlock()
	return pinLayoutValue
}
