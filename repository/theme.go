/*
   file:           repository/theme.go
   description:    Repositori data untuk theme
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package repository



import (
	"database/sql"
	"go-ecb/app/types"
	"go-ecb/pkg/logging"
	"github.com/go-gorp/gorp"
)

type ThemeRepository struct {
	dbmap *gorp.DbMap
}

func NewThemeRepository(dbmap *gorp.DbMap) *ThemeRepository {
	return &ThemeRepository{dbmap: dbmap}
}

func (r *ThemeRepository) GetThemes() ([]types.Theme, error) {
	var themes []types.Theme
	_, err := r.dbmap.Select(&themes, `
		SELECT * FROM themes
		ORDER BY id ASC
	`)
	if err != nil {
		logging.Logger().Warnf("Error querying themes: %v", err)
		return nil, err
	}
	return themes, nil
}

func (r *ThemeRepository) GetThemeByName(name string) (*types.Theme, error) {
	var theme types.Theme
	err := r.dbmap.SelectOne(&theme, `
		SELECT * FROM themes
		WHERE nama = ?
		LIMIT 1
	`, name)

	if err != nil {
		if err == sql.ErrNoRows {
			// logging.Logger().Warnf("Error finding theme %s: %v", name, err) // Optional: maybe distinct warning for not found vs error
			return nil, err
		}
		logging.Logger().Warnf("Error finding theme %s: %v", name, err)
		return nil, err
	}
	return &theme, nil
}
