/*
   file:           services/gpio/wiringpi.go
   description:    Mapping WiringPi ke BCM dan Fisik untuk Raspberry Pi 3
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package gpio

import (
	"fmt"
)

type PinMapItem struct {
	WiringPi string
	BCM      string
	Physical string
}

var RPiWiringPiMap = map[string]PinMapItem{
	"0":  {WiringPi: "0", BCM: "17", Physical: "11"},
	"1":  {WiringPi: "1", BCM: "18", Physical: "12"},
	"2":  {WiringPi: "2", BCM: "27", Physical: "13"},
	"3":  {WiringPi: "3", BCM: "22", Physical: "15"},
	"4":  {WiringPi: "4", BCM: "23", Physical: "16"},
	"5":  {WiringPi: "5", BCM: "24", Physical: "18"},
	"6":  {WiringPi: "6", BCM: "25", Physical: "22"},
	"7":  {WiringPi: "7", BCM: "4", Physical: "7"},
	"8":  {WiringPi: "8", BCM: "2", Physical: "3"},
	"9":  {WiringPi: "9", BCM: "3", Physical: "5"},
	"10": {WiringPi: "10", BCM: "8", Physical: "24"},
	"11": {WiringPi: "11", BCM: "7", Physical: "26"},
	"12": {WiringPi: "12", BCM: "10", Physical: "19"},
	"13": {WiringPi: "13", BCM: "9", Physical: "21"},
	"14": {WiringPi: "14", BCM: "11", Physical: "23"},
	"15": {WiringPi: "15", BCM: "14", Physical: "8"},
	"16": {WiringPi: "16", BCM: "15", Physical: "10"},
	"21": {WiringPi: "21", BCM: "5", Physical: "29"},
	"22": {WiringPi: "22", BCM: "6", Physical: "31"},
	"23": {WiringPi: "23", BCM: "13", Physical: "33"},
	"24": {WiringPi: "24", BCM: "19", Physical: "35"},
	"25": {WiringPi: "25", BCM: "26", Physical: "37"},
	"26": {WiringPi: "26", BCM: "12", Physical: "32"},
	"27": {WiringPi: "27", BCM: "16", Physical: "36"},
	"28": {WiringPi: "28", BCM: "20", Physical: "38"},
	"29": {WiringPi: "29", BCM: "21", Physical: "40"},
	"30": {WiringPi: "30", BCM: "0", Physical: "27"},
	"31": {WiringPi: "31", BCM: "1", Physical: "28"},
}

func WiringPiToBCM(pin string) (string, error) {
	if item, ok := RPiWiringPiMap[pin]; ok {
		return "GPIO" + item.BCM, nil
	}
	return "", fmt.Errorf("unknown WiringPi pin: %s", pin)
}

func WiringPiToPhysical(pin string) string {
	if item, ok := RPiWiringPiMap[pin]; ok {
		return item.Physical
	}
	return "-"
}

func ValidatePinAccess(pin string) error {
	if _, ok := RPiWiringPiMap[pin]; ok {
		return nil
	}
	return fmt.Errorf("pin %s is not a valid or accessible GPIO pin on this device", pin)
}
