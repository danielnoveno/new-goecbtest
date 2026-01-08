/*
   file:           app/types/data_record.go
   description:    Record Model dan helper UI untuk ecbdata
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type EcbData struct {
	ID        int       `db:"id"`
	Tgl       time.Time `db:"tgl"`
	Jam       time.Time `db:"jam"`
	Wc        string    `db:"wc"`
	Prdline   string    `db:"prdline"`
	Ctgr      string    `db:"ctgr"`
	Sn        string    `db:"sn"`
	Fgtype    string    `db:"fgtype"`
	Spc       string    `db:"spc"`
	Comptype  string    `db:"comptype"`
	Compcode  string    `db:"compcode"`
	Po        string    `db:"po"`
	Status    string    `db:"status"`
	Sendsts   string    `db:"sendsts"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}