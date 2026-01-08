/*
    file:           app/types/theme_record.go
    description:    Record Model dan helper UI untuk theme
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package types

import "time"

type Theme struct {
	ID                   int       `db:"id"`
	Nama                 string    `db:"nama"`
	Keterangan           string    `db:"keterangan"`
	ColorBackground      string    `db:"color_background"`
	ColorForeground      string    `db:"color_foreground"`
	ColorText            string    `db:"color_text"`
	ColorButton          string    `db:"color_button"`
	ColorDisabled        string    `db:"color_disabled"`
	ColorError           string    `db:"color_error"`
	ColorFocus           string    `db:"color_focus"`
	ColorHover           string    `db:"color_hover"`
	HeaderStart          string    `db:"header_start"`
	HeaderEnd            string    `db:"header_end"`
	Accent               string    `db:"accent"`
	ColorInputBackground string    `db:"color_input_background"`
	ColorPlaceholder     string    `db:"color_placeholder"`
	ColorPrimary         string    `db:"color_primary"`
	ColorScrollbar       string    `db:"color_scrollbar"`
	ColorSelection       string    `db:"color_selection"`
	ColorNavbar          string    `db:"color_navbar"`
	ColorFooter          string    `db:"color_footer"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
}
