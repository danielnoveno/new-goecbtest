/*
    file:           app/types/data_record.go
    description:    Record Model dan helper UI untuk ecbdata
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type EcbData struct {
	ID        int       `db:"id" gorm:"primaryKey;autoIncrement" json:"id"`
	Tgl       time.Time `db:"tgl" gorm:"column:tgl" json:"tgl"`
	Jam       time.Time `db:"jam" gorm:"column:jam" json:"jam"`
	Wc        string    `db:"wc" gorm:"column:wc" json:"wc"`
	Prdline   string    `db:"prdline" gorm:"column:prdline" json:"prdLine"`
	Ctgr      string    `db:"ctgr" gorm:"column:ctgr" json:"ctgr"`
	Sn        string    `db:"sn" gorm:"column:sn" json:"sn"`
	Fgtype    string    `db:"fgtype" gorm:"column:fgtype" json:"fgtype"`
	Spc       string    `db:"spc" gorm:"column:spc" json:"spc"`
	Comptype  string    `db:"comptype" gorm:"column:comptype" json:"compType"`
	Compcode  string    `db:"compcode" gorm:"column:compcode" json:"compCode"`
	Po        string    `db:"po" gorm:"column:po" json:"po"`
	Status    string    `db:"status" gorm:"column:status" json:"status"`
	Sendsts   string    `db:"sendsts" gorm:"column:sendsts" json:"sendSts"`
	CreatedAt time.Time `db:"created_at" gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" gorm:"column:updated_at" json:"updatedAt"`
}

// func (EcbData) TableName() string {
// 	return "ecbdatas"
// }
