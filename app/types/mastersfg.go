/*
   file:           app/types/mastersfg.go
   description:    Model dan helper UI untuk mastersfg
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type Mastersfg struct {
	ID        int       `db:"id"`
	Plant     string    `db:"plant"`
	Mattype   string    `db:"mattype"`
	Matdesc   string    `db:"matdesc"`
	Sfgtype   string    `db:"sfgtype"`
	Sfgdesc   string    `db:"sfgdesc"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
