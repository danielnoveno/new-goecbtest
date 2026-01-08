/*
   file:           app/types/navigation_record.go
   description:    Record Model dan helper UI untuk navigation
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import (
	"database/sql"
	"time"
)

type Navigation struct {
	ID          int           `db:"id"`
	ParentId    sql.NullInt64 `db:"parent_id"`
	Icon        string        `db:"icon"`
	Title       string        `db:"title"`
	Description string        `db:"description"`
	Url         string        `db:"url"`
	Route       string        `db:"route"`
	Mode        int           `db:"mode"`
	Urutan      int           `db:"urutan"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
}
