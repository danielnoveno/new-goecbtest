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
	ID        int 		`db:"id"`
	Section   string	`db:"section"`
	Variable  string	`db:"variable"`
	Value     string	`db:"value"`
	Ordering  string	`db:"ordering"`
	CreatedAt time.Time	`db:"created_at"`
	UpdatedAt time.Time	`db:"updated_at"`
}
