/*
    file:           app/types/theme_record.go
    description:    Record Model dan helper UI untuk theme
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type Theme struct {
	ID                   int       `db:"id" json:"id"`
	Nama                 string    `db:"nama" json:"nama"`
	Keterangan           string    `db:"keterangan" json:"keterangan"`
	ColorBackground      string    `db:"color_background" json:"colorBackground"`
	ColorForeground      string    `db:"color_foreground" json:"colorForeground"`
	ColorText            string    `db:"color_text" json:"colorText"`
	ColorButton          string    `db:"color_button" json:"colorButton"`
	ColorDisabled        string    `db:"color_disabled" json:"colorDisabled"`
	ColorError           string    `db:"color_error" json:"colorError"`
	ColorFocus           string    `db:"color_focus" json:"colorFocus"`
	ColorHover           string    `db:"color_hover" json:"colorHover"`
	HeaderStart          string    `db:"header_start" json:"headerStart"`
	HeaderEnd            string    `db:"header_end" json:"headerEnd"`
	Accent               string    `db:"accent" json:"accent"`
	ColorInputBackground string    `db:"color_input_background" json:"colorInputBackground"`
	ColorPlaceholder     string    `db:"color_placeholder" json:"colorPlaceholder"`
	ColorPrimary         string    `db:"color_primary" json:"colorPrimary"`
	ColorScrollbar       string    `db:"color_scrollbar" json:"colorScrollbar"`
	ColorSelection       string    `db:"color_selection" json:"colorSelection"`
	ColorNavbar          string    `db:"color_navbar" json:"colorNavbar"`
	ColorFooter          string    `db:"color_footer" json:"colorFooter"`
	CreatedAt            time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt            time.Time `db:"updated_at" json:"updatedAt"`
}
