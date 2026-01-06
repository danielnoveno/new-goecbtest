/*
    file:           services/maintenance/service.go
    description:    Layanan pemeliharaan untuk service
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package maintenance

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"go-ecb/app/types"
	"go-ecb/repository"
	"go-ecb/services/gpio"
)

const (
	pinsSection     = "maintenance_pins"
	pinVarUnderTest = "pin_under_test"
	pinVarFail      = "pin_fail"
	pinVarPass      = "pin_pass"
	pinVarLine      = "pin_line"
	pinVarStart     = "pin_start"
	pinVarStartAlt  = "pin_start_alt"
	pinVarReset     = "pin_reset"
	pinVarResetAlt  = "pin_reset_alt"
)

// PinController menjabarkan operasi yang dibutuhkan UI untuk membaca dan menyimpan
// mapping wiringPi tanpa mengandaikan cara persistensinya.
type PinController interface {
	GetPins() gpio.PinLayout
	SavePins(gpio.PinLayout) error
}

// PinConfigService mengelola entri mapping pin maintenance di tabel konfigurasi ECB
// dan menyinkronkannya dengan helper GPIO.
type PinConfigService struct {
	store repository.EcbConfigStore
	mu    sync.RWMutex
	pins  gpio.PinLayout
}

// NewPinConfigService membuat layanan dan me-load layout pin default. Panggil Refresh
// sebelum menggunakan pin yang dimuat.
func NewPinConfigService(store repository.EcbConfigStore) *PinConfigService {
	if store == nil {
		return nil
	}
	return &PinConfigService{
		store: store,
		pins:  gpio.DefaultPinLayout(),
	}
}

// Refresh memuat ulang mapping pin dari database dan memperbaruinya di helper GPIO
// agar perintah memakai nilai terbaru.
func (s *PinConfigService) Refresh() error {
	if s == nil || s.store == nil {
		return fmt.Errorf("pin config service unavailable")
	}
	records, err := s.store.FindEcbConfigsBySection(pinsSection)
	if err != nil {
		s.apply(gpio.DefaultPinLayout())
		return fmt.Errorf("failed to load maintenance pins: %w", err)
	}
	layout := gpio.DefaultPinLayout()
	for _, record := range records {
		value := strings.TrimSpace(record.Value)
		if value == "" {
			continue
		}
		switch record.Variable {
		case pinVarUnderTest:
			layout.UnderTest = value
		case pinVarFail:
			layout.Fail = value
		case pinVarPass:
			layout.Pass = value
		case pinVarLine:
			layout.LineSelect = value
		case pinVarStart:
			layout.Start = value
		case pinVarStartAlt:
			layout.StartAlt = value
		case pinVarReset:
			layout.Reset = value
		case pinVarResetAlt:
			layout.ResetAlt = value
		}
	}
	s.apply(layout)
	return nil
}

// GetPins mengembalikan layout pin terakhir yang dimuat.
func (s *PinConfigService) GetPins() gpio.PinLayout {
	if s == nil {
		return gpio.DefaultPinLayout()
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.pins
}

// SavePins menyimpan layout pin ke database dan memperbarui helper GPIO bila berhasil.
func (s *PinConfigService) SavePins(layout gpio.PinLayout) error {
	if s == nil || s.store == nil {
		return fmt.Errorf("pin config service unavailable")
	}
	normalized := gpio.NormalizePinLayout(layout)
	variables := map[string]string{
		pinVarUnderTest: normalized.UnderTest,
		pinVarFail:      normalized.Fail,
		pinVarPass:      normalized.Pass,
		pinVarLine:      normalized.LineSelect,
		pinVarStart:     normalized.Start,
		pinVarStartAlt:  normalized.StartAlt,
		pinVarReset:     normalized.Reset,
		pinVarResetAlt:  normalized.ResetAlt,
	}
	now := time.Now()
	for name, value := range variables {
		if err := s.persist(name, value, now); err != nil {
			return err
		}
	}
	s.apply(normalized)
	return nil
}

// apply adalah fungsi untuk apply.
func (s *PinConfigService) apply(layout gpio.PinLayout) {
	s.mu.Lock()
	normalized := gpio.NormalizePinLayout(layout)
	s.pins = normalized
	s.mu.Unlock()
	gpio.SetPinLayout(normalized)
}

// persist adalah fungsi untuk persist.
func (s *PinConfigService) persist(variable, value string, now time.Time) error {
	config, err := s.store.FindEcbConfigBySectionAndVariable(pinsSection, variable)
	if err != nil {
		if err.Error() == "ecbconfig not found" {
			newConfig := &types.EcbConfig{
				Section:   pinsSection,
				Variable:  variable,
				Value:     value,
				Ordering:  "0",
				CreatedAt: now,
				UpdatedAt: now,
			}
			return s.store.CreateEcbConfig(newConfig)
		}
		return err
	}
	config.Value = value
	config.UpdatedAt = now
	return s.store.UpdateEcbConfig(config)
}
