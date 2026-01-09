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

func init() {
	if _, err := host.Init(); err != nil {
		log.Printf("periph host initialization failed (%v), falling back to noop driver", err)
		SetDriver(newNoopDriver())
		return
	}
	log.Printf("GPIO hardware driver (periph) initialized successfully")
	SetDriver(&periphDriver{pins: make(map[string]gpio.PinIO)})
}

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

func normalizePinName(pin string) string {
	trimmed := strings.TrimSpace(pin)
	if trimmed == "" {
		return trimmed
	}

	if bcm, err := WiringPiToBCM(trimmed); err == nil {
		return bcm
	}
	
	if strings.HasPrefix(strings.ToUpper(trimmed), "GPIO") {
		return trimmed
	}

	if _, err := strconv.Atoi(trimmed); err == nil {
		return "GPIO" + trimmed
	}
	return trimmed
}

func (d *periphDriver) Write(pin string, level Level) error {
	if err := ValidatePinAccess(pin); err != nil {
		return err
	}
	p, err := d.pin(pin)
	if err != nil {
		return err
	}
	return p.Out(levelToPeriph(level))
}

func (d *periphDriver) Read(pin string) (Level, error) {
	if err := ValidatePinAccess(pin); err != nil {
		return LevelLow, err
	}
	p, err := d.pin(pin)
	if err != nil {
		return LevelLow, err
	}
	if err := p.In(gpio.PullUp, gpio.NoEdge); err != nil {
		return LevelLow, err
	}
	return p.Read() == gpio.High, nil
}

func (d *periphDriver) SetMode(pin string, mode PinMode) error {
	if err := ValidatePinAccess(pin); err != nil {
		return err
	}
	p, err := d.pin(pin)
	if err != nil {
		return err
	}
	if mode == ModeOutput {
		return p.Out(gpio.Low)
	}
	return p.In(gpio.PullUp, gpio.NoEdge)
}

func levelToPeriph(level Level) gpio.Level {
	if level {
		return gpio.High
	}
	return gpio.Low
}
