/*
    file:           app/types/purchase_order.go
    description:    Model dan helper UI untuk ecbpo
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type EcbPo struct {
	ID         int       `db:"id" gorm:"primaryKey;autoIncrement" json:"id"`
	WorkCenter string    `db:"workcenter" json:"workCenter"`
	Po         string    `db:"po" json:"po"`
	Sn         string    `db:"sn" json:"sn"`
	Ctype      string    `db:"ctype" json:"ctype"`
	UpdatedBy  int       `db:"updated_by" json:"updatedBy"`
	Status     string    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
}
