/*
   file:           app/types/purchase_order.go
   description:    Model dan helper UI untuk ecbpo
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type EcbPo struct {
	ID         int       `db:"id"`
	WorkCenter string    `db:"workcenter"`
	Po         string    `db:"po"`
	Sn         string    `db:"sn"`
	Ctype      string    `db:"ctype"`
	UpdatedBy  int       `db:"updated_by"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
