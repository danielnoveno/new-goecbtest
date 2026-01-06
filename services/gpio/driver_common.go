/*
    file:           services/gpio/driver_common.go
    description:    Driver GPIO untuk driver common
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package gpio

import (
	"strings"
	"sync"
)

type Level bool

const (
	LevelLow  Level = false
	LevelHigh Level = true
)

// String adalah fungsi untuk string.
func (l Level) String() string {
	if l {
		return "1"
	}
	return "0"
}

// LevelFromString adalah fungsi untuk level from string.
func LevelFromString(value string) Level {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "high", "true", "on":
		return LevelHigh
	default:
		return LevelLow
	}
}

type PinMode string

const (
	ModeInput  PinMode = "in"
	ModeOutput PinMode = "out"
)

type driver interface {
	Write(pin string, level Level) error
	Read(pin string) (Level, error)
	SetMode(pin string, mode PinMode) error
}

var (
	driverMu sync.RWMutex
	hwDriver driver = newNoopDriver()
)

// SetDriver adalah fungsi untuk mengatur driver.
func SetDriver(d driver) {
	driverMu.Lock()
	defer driverMu.Unlock()
	if d != nil {
		hwDriver = d
	}
}

// WritePin adalah fungsi untuk menulis pin.
func WritePin(pin string, level Level) error {
	driverMu.RLock()
	defer driverMu.RUnlock()
	if hwDriver == nil {
		return nil
	}
	return hwDriver.Write(pin, level)
}

// ReadPin adalah fungsi untuk membaca pin.
func ReadPin(pin string) (Level, error) {
	driverMu.RLock()
	defer driverMu.RUnlock()
	if hwDriver == nil {
		return LevelLow, nil
	}
	return hwDriver.Read(pin)
}

// SetPinMode adalah fungsi untuk mengatur pin mode.
func SetPinMode(pin string, mode PinMode) error {
	driverMu.RLock()
	defer driverMu.RUnlock()
	if hwDriver == nil {
		return nil
	}
	return hwDriver.SetMode(pin, mode)
}

type noopDriver struct {
	mu     sync.Mutex
	states map[string]Level
}

// newNoopDriver adalah fungsi untuk baru noop driver.
func newNoopDriver() *noopDriver {
	return &noopDriver{states: make(map[string]Level)}
}

// Write adalah fungsi untuk menulis.
func (d *noopDriver) Write(pin string, level Level) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.states[pin] = level
	return nil
}

// Read adalah fungsi untuk membaca.
func (d *noopDriver) Read(pin string) (Level, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if val, ok := d.states[pin]; ok {
		return val, nil
	}
	return LevelLow, nil
}

// SetMode adalah fungsi untuk mengatur mode.
func (d *noopDriver) SetMode(pin string, mode PinMode) error {
	return nil
}
