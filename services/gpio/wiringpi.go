/*
   file:           services/gpio/wiringpi.go
   description:    Mapping WiringPi ke BCM dan Fisik untuk Raspberry Pi 3
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package gpio

type PinMapItem struct {
	WiringPi string
	BCM      string
	Physical string
}

var RPiWiringPiMap = map[string]PinMapItem{
	"0":  {WiringPi: "0", Physical: "11"},
	"1":  {WiringPi: "1", Physical: "12"},
	"2":  {WiringPi: "2", Physical: "13"},
	"3":  {WiringPi: "3", Physical: "15"},
	"4":  {WiringPi: "4", Physical: "16"},
	"5":  {WiringPi: "5", Physical: "18"},
	"6":  {WiringPi: "6", Physical: "22"},
	"7":  {WiringPi: "7", Physical: "7"},
	"8":  {WiringPi: "8", Physical: "3"},
	"9":  {WiringPi: "9", Physical: "5"},
	"10": {WiringPi: "10", Physical: "24"},
	"11": {WiringPi: "11", Physical: "26"},
	"12": {WiringPi: "12", Physical: "19"},
	"13": {WiringPi: "13", Physical: "21"},
	"14": {WiringPi: "14", Physical: "23"},
	"15": {WiringPi: "15", Physical: "8"},
	"16": {WiringPi: "16", Physical: "10"},
	"21": {WiringPi: "21", Physical: "29"},
	"22": {WiringPi: "22", Physical: "31"},
	"23": {WiringPi: "23", Physical: "33"},
	"24": {WiringPi: "24", Physical: "35"},
	"25": {WiringPi: "25", Physical: "37"},
	"26": {WiringPi: "26", Physical: "32"},
	"27": {WiringPi: "27", Physical: "36"},
	"28": {WiringPi: "28", Physical: "38"},
	"29": {WiringPi: "29", Physical: "40"},
	"30": {WiringPi: "30", Physical: "27"},
	"31": {WiringPi: "31", Physical: "28"},
}

// func WiringPiToBCM(pin string) (string, error) {
// 	if item, ok := RPiWiringPiMap[pin]; ok {
// 		return "GPIO" + item.BCM, nil
// 	}

// 	return "", fmt.Errorf("unknown WiringPi pin: %s", pin)
// }

func WiringPiToPhysical(pin string) string {
	if item, ok := RPiWiringPiMap[pin]; ok {
		return item.Physical
	}
	return "-"
}

// func IsBCMFormat(pin string) bool {
//     return len(pin) > 4 && pin[:4] == "GPIO"
// }
