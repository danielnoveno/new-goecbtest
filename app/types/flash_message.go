/*
    file:           app/types/flash_message.go
    description:    Model dan helper UI untuk flash
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

type FlashLevel string

const (
	FlashLevelInfo    FlashLevel = "info"
	FlashLevelSuccess FlashLevel = "success"
	FlashLevelWarning FlashLevel = "warning"
	FlashLevelError   FlashLevel = "error"
)

type FlashMessage struct {
	Level FlashLevel
	Title string
	Body  string
}

type FlashNotifier interface {
	Notify(msg FlashMessage)
	Subscribe(handler func(FlashMessage)) func()
}
