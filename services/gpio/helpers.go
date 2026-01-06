/*
    file:           services/gpio/helpers.go
    description:    Driver GPIO untuk helpers
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package gpio

import "log"

// readLevel adalah fungsi untuk membaca level.
func readLevel(pin string) string {
	value, err := ReadPin(pin)
	if err != nil {
		log.Printf("gpio read %s failed: %v", pin, err)
		return "0"
	}
	return value.String()
}

// writeLevel adalah fungsi untuk menulis level.
func writeLevel(pin string, level Level) {
	if err := WritePin(pin, level); err != nil {
		log.Printf("gpio write %s=%s failed: %v", pin, level.String(), err)
	}
}

// configureMode adalah fungsi untuk configure mode.
func configureMode(pin string, mode PinMode) {
	if err := SetPinMode(pin, mode); err != nil {
		log.Printf("gpio mode %s %s failed: %v", pin, mode, err)
	}
}
