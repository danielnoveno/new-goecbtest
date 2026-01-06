/*
   file:           repository/theme.go
   description:    Repositori data untuk theme
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package repository

import (
	"go-ecb/app/types"
	"go-ecb/db"
	"log"
)

// GetThemes adalah fungsi untuk mengambil tema.
func GetThemes() ([]types.Theme, error) {
	rows, err := db.Query(`
		SELECT id, nama, keterangan,
		       color_background, color_foreground, color_text,
		       color_button, color_disabled, color_error,
		       color_focus, color_hover,
		       color_input_background, color_placeholder, color_primary,
		       color_scrollbar, color_selection, color_navbar, color_footer,
		       header_start, header_end, accent,
		       created_at, updated_at
		FROM themes
		ORDER BY id ASC
	`)
	if err != nil {
		log.Printf("Error querying themes: %v", err)
		return nil, err
	}
	defer rows.Close()

	var themes []types.Theme
	for rows.Next() {
		var theme types.Theme
		if err := rows.Scan(
			&theme.ID,
			&theme.Nama,
			&theme.Keterangan,
			&theme.ColorBackground,
			&theme.ColorForeground,
			&theme.ColorText,
			&theme.ColorButton,
			&theme.ColorDisabled,
			&theme.ColorError,
			&theme.ColorFocus,
			&theme.ColorHover,
			&theme.ColorInputBackground,
			&theme.ColorPlaceholder,
			&theme.ColorPrimary,
			&theme.ColorScrollbar,
			&theme.ColorSelection,
			&theme.ColorNavbar,
			&theme.ColorFooter,
			&theme.HeaderStart,
			&theme.HeaderEnd,
			&theme.Accent,
			&theme.CreatedAt,
			&theme.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning theme row: %v", err)
			return nil, err
		}
		themes = append(themes, theme)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return themes, nil
}

// GetThemeByName adalah fungsi untuk mengambil tema by name.
func GetThemeByName(name string) (*types.Theme, error) {
	row := db.QueryRow(`
		SELECT id, nama, keterangan,
		       color_background, color_foreground, color_text,
		       color_button, color_disabled, color_error,
		       color_focus, color_hover,
		       color_input_background, color_placeholder, color_primary,
		       color_scrollbar, color_selection, color_navbar, color_footer,
		       header_start, header_end, accent,
		       created_at, updated_at
		FROM themes
		WHERE nama = ?
		LIMIT 1
	`, name)

	var theme types.Theme
	if err := row.Scan(
		&theme.ID,
		&theme.Nama,
		&theme.Keterangan,
		&theme.ColorBackground,
		&theme.ColorForeground,
		&theme.ColorText,
		&theme.ColorButton,
		&theme.ColorDisabled,
		&theme.ColorError,
		&theme.ColorFocus,
		&theme.ColorHover,
		&theme.ColorInputBackground,
		&theme.ColorPlaceholder,
		&theme.ColorPrimary,
		&theme.ColorScrollbar,
		&theme.ColorSelection,
		&theme.ColorNavbar,
		&theme.ColorFooter,
		&theme.HeaderStart,
		&theme.HeaderEnd,
		&theme.Accent,
		&theme.CreatedAt,
		&theme.UpdatedAt,
	); err != nil {
		log.Printf("Error finding theme %s: %v", name, err)
		return nil, err
	}
	return &theme, nil
}
