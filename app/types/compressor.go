/*
    file:           app/types/compressor_record.go
    description:    Record Model dan helper UI untuk compressor
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type Compressor struct {
	ID         int       `db:"id"`
	Ctype      string    `db:"ctype"`
	Merk       string    `db:"merk"`
	Type       string    `db:"type"`
	Itemcode   string    `db:"itemcode"`
	ForceScan  int       `db:"force_scan"`
	FamilyCode string    `db:"familycode"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
