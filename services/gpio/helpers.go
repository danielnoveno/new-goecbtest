/*
   file:           services/gpio/helpers.go
   description:    Driver GPIO untuk helpers
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package gpio

import "go-ecb/pkg/logging"

func readLevel(pin string) string {
	value, err := ReadPin(pin)
	if err != nil {
		logging.Logger().Warnf("gpio read %s failed: %v", pin, err)
		return "0"
	}
	return value.String()
}

func writeLevel(pin string, level Level) {
	if err := WritePin(pin, level); err != nil {
		logging.Logger().Warnf("gpio write %s=%s failed: %v", pin, level.String(), err)
	}
}

func configureMode(pin string, mode PinMode) {
	if err := SetPinMode(pin, mode); err != nil {
		logging.Logger().Warnf("gpio mode %s %s failed: %v", pin, mode, err)
	}
}
