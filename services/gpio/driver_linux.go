//go:build linux
// +build linux

/*
    file:           services/gpio/driver_linux.go
    description:    Driver GPIO untuk driver linux
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package gpio

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

type periphDriver struct {
	mu   sync.Mutex
	pins map[string]gpio.PinIO
}

// init adalah fungsi untuk inisialisasi.
func init() {
	if _, err := host.Init(); err != nil {
		log.Printf("periph host initialization failed (%v), falling back to noop driver", err)
		SetDriver(newNoopDriver())
		return
	}
	SetDriver(&periphDriver{pins: make(map[string]gpio.PinIO)})
}

// pin adalah fungsi untuk pin.
func (d *periphDriver) pin(pin string) (gpio.PinIO, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if p, ok := d.pins[pin]; ok {
		return p, nil
	}

	name := normalizePinName(pin)
	p := gpioreg.ByName(name)
	if p == nil {
		return nil, fmt.Errorf("gpio pin %s (%s) not found", pin, name)
	}
	d.pins[pin] = p
	return p, nil
}

// normalizePinName adalah fungsi untuk normalize pin name.
func normalizePinName(pin string) string {
	trimmed := strings.TrimSpace(pin)
	if trimmed == "" {
		return trimmed
	}
	if strings.HasPrefix(strings.ToUpper(trimmed), "GPIO") {
		return trimmed
	}
	if _, err := strconv.Atoi(trimmed); err == nil {
		return "GPIO" + trimmed
	}
	return trimmed
}

// Write adalah fungsi untuk menulis.
func (d *periphDriver) Write(pin string, level Level) error {
	p, err := d.pin(pin)
	if err != nil {
		return err
	}
	return p.Out(levelToPeriph(level))
}

// Read adalah fungsi untuk membaca.
func (d *periphDriver) Read(pin string) (Level, error) {
	p, err := d.pin(pin)
	if err != nil {
		return LevelLow, err
	}
	if err := p.In(gpio.PullNoChange, gpio.NoEdge); err != nil {
		return LevelLow, err
	}
	return p.Read() == gpio.High, nil
}

// SetMode adalah fungsi untuk mengatur mode.
func (d *periphDriver) SetMode(pin string, mode PinMode) error {
	p, err := d.pin(pin)
	if err != nil {
		return err
	}
	if mode == ModeOutput {
		return p.Out(gpio.Low)
	}
	return p.In(gpio.PullNoChange, gpio.NoEdge)
}

// levelToPeriph adalah fungsi untuk level to periph.
func levelToPeriph(level Level) gpio.Level {
	if level {
		return gpio.High
	}
	return gpio.Low
}
