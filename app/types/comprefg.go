/*
    file:           app/types/compressor_reference.go
    description:    Reference Model dan helper UI untuk comprefg
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type Comprefg struct {
	ID        int       `db:"id" gorm:"primaryKey;autoIncrement" json:"id"`
	Ctype     string    `db:"ctype" json:"ctype"`
	Barcode   string    `db:"barcode" json:"barcode"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
