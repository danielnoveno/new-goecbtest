/*
    file:           app/types/state_snapshot.go
    description:    Model dan helper UI untuk ecbstate
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type EcbState struct {
	ID		  int       `db:"id"`	
	Tgl       time.Time `db:"tgl"`
	EcbState  string    `db:"ecbstate"`
	ReadState string    `db:"readstate"`
}
