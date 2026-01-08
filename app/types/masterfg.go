/*
    file:           app/types/master_fg.go
    description:    Model dan helper UI untuk masterfg
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type Masterfg struct {
	ID           int       `db:"id"`
	Mattype      string    `db:"mattype"`
	Matdesc      string    `db:"matdesc"`
	Fgtype       string    `db:"fgtype"`
	AgingTypesId int       `db:"aging_tipes_id"`
	Kdbar        string    `db:"kdbar"`
	Warna        string    `db:"warna"`
	Lotinv       string    `db:"lotinv"`
	Attrib       string    `db:"attrib"`
	Category     string    `db:"category"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
