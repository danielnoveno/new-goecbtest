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
	ID          int           `gorm:"primaryKey;autoIncrement" json:"id"`
	ParentId    sql.NullInt64 `gorm:"parent_id" json:"parentID"`
	Icon        string        `gorm:"icon" json:"icon"`
	Title       string        `gorm:"title" json:"title"`
	Description string        `gorm:"description" json:"description"`
	Url         string        `gorm:"url" json:"url"`
	Route       string        `gorm:"route" json:"route"`
	Mode        int           `gorm:"mode" json:"mode"`
	Urutan      int           `gorm:"urutan" json:"urutan"`
	CreatedAt   time.Time     `gorm:"created_at" json:"createdAt"`
	UpdatedAt   time.Time     `gorm:"updated_at" json:"updatedAt"`
}
