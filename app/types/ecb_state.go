/*
    file:           app/types/state_snapshot.go
    description:    Model dan helper UI untuk ecbstate
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type EcbState struct {
	ID		  int       `gorm:"primaryKey;autoIncrement" json:"id"`	
	Tgl       time.Time `gorm:"tgl" json:"tgl"`
	EcbState  string    `gorm:"ecbstate" json:"ecbState"`
	ReadState string    `gorm:"readstate" json:"readState"`
}
