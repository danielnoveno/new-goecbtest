/*
    file:           app/types/mastersfg.go
    description:    Model dan helper UI untuk mastersfg
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type Mastersfg struct {
	ID        int       `db:"id" gorm:"primaryKey;autoIncrement" json:"id"`
	Plant     string    `db:"plant" json:"plant"`
	Mattype   string    `db:"mattype" json:"mattype"`
	Matdesc   string    `db:"matdesc" json:"matdesc"`
	Sfgtype   string    `db:"sfgtype" json:"sfgtype"`
	Sfgdesc   string    `db:"sfgdesc" json:"sfgdesc"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
