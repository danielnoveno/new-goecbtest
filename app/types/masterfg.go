/*
    file:           app/types/master_fg.go
    description:    Model dan helper UI untuk masterfg
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type Masterfg struct {
	ID           int       `db:"id" gorm:"primaryKey;autoIncrement" json:"id"`
	Mattype      string    `db:"mattype" json:"mattype"`
	Matdesc      string    `db:"matdesc" json:"matdesc"`
	Fgtype       string    `db:"fgtype" json:"fgtype"`
	AgingTypesId int       `db:"aging_tipes_id" json:"agingTypesID"`
	Kdbar        string    `db:"kdbar" json:"kdbar"`
	Warna        string    `db:"warna" json:"warna"`
	Lotinv       string    `db:"lotinv" json:"lotinv"`
	Attrib       string    `db:"attrib" json:"attrib"`
	Category     string    `db:"category" json:"category"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}
