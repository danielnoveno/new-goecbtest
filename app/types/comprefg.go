/*
   file:           app/types/compressor_reference.go
   description:    Reference Model dan helper UI untuk comprefg
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type Comprefg struct {
	ID        int       `db:"id"`
	Ctype     string    `db:"ctype"`
	Barcode   string    `db:"barcode"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
