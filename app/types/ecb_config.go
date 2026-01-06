/*
    file:           app/types/configuration_entry.go
    description:    Model dan helper UI untuk ecbconfig
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types


import (
	"time"
)

type EcbConfig struct {
	ID        int 		`gorm:"primaryKey;autoIncrement" json:"id"`
	Section   string	`gorm:"section" json:"section"`
	Variable  string	`gorm:"variable" json:"variable"`
	Value     string	`gorm:"value" json:"value"`
	Ordering  string	`gorm:"ordering" json:"ordering"`
	CreatedAt time.Time	`gorm:"created_at" json:"createdAt"`
	UpdatedAt time.Time	`gorm:"updated_at" json:"updatedAt"`
}
