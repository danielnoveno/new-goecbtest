/*
    file:           app/types/compressor_record.go
    description:    Record Model dan helper UI untuk compressor
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type Compressor struct {
	ID         int       `db:"id" gorm:"primaryKey;autoIncrement" json:"id"`
	Ctype      string    `db:"ctype" json:"ctype"`
	Merk       string    `db:"merk" json:"merk"`
	Type       string    `db:"type" json:"type"`
	Itemcode   string    `db:"itemcode" json:"itemcode"`
	ForceScan  int       `db:"force_scan" json:"forceScan"`
	FamilyCode string    `db:"familycode" json:"familyCode"`
	Status     string    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
}
