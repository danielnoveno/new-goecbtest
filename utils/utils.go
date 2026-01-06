/*
    file:           utils/utils.go
    description:    Utilitas pendukung untuk utils
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"go-ecb/services/gpio"
)

var Validate = validator.New()

// ExecCommand adalah fungsi untuk exec command.
func ExecCommand(name string, arg ...string) {
	fmt.Printf("Executing command: %s %v\n", name, arg)
	if strings.ToLower(name) != "gpio" || len(arg) == 0 {
		return
	}

	switch strings.ToLower(arg[0]) {
	case "write":
		if len(arg) >= 3 {
			if err := gpio.WritePin(arg[1], gpio.LevelFromString(arg[2])); err != nil {
				fmt.Printf("gpio write failed: %v\n", err)
			}
		}
	case "mode":
		if len(arg) >= 3 {
			if err := gpio.SetPinMode(arg[1], gpio.PinMode(strings.ToLower(arg[2]))); err != nil {
				fmt.Printf("gpio mode failed: %v\n", err)
			}
		}
	case "read":
		if len(arg) >= 2 {
			if val, err := gpio.ReadPin(arg[1]); err == nil {
				fmt.Printf("gpio read %s -> %s\n", arg[1], val.String())
			} else {
				fmt.Printf("gpio read failed: %v\n", err)
			}
		}
	}
}
