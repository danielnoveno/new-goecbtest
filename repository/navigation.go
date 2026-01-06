/*
    file:           repository/navigation.go
    description:    Repositori data untuk navigation
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package repository

import (
	"database/sql"

	"go-ecb/app/types"
)

type NavigationStore interface {
	Insert(navigation *types.Navigation) error
	Update(navigation *types.Navigation) error
	Delete(id int) error
	GetAll() ([]types.Navigation, error)
	FindNavigationByID(id int) (*types.Navigation, error)
	FindNavigationByParentID(parentID int) ([]*types.Navigation, error)
	FindRootNavigations() ([]*types.Navigation, error)
	ListAll() ([]*types.Navigation, error)
	FindNavigationByUrutan(urutan int, parentID *int) (*types.Navigation, error)
}

type navigationRepo struct {
	db *sql.DB
}

// NewNavigationRepository adalah fungsi untuk baru navigasi repository.
func NewNavigationRepository(db *sql.DB) NavigationStore {
	return &navigationRepo{db: db}
}

// Insert adalah fungsi untuk insert.
func (r *navigationRepo) Insert(navigation *types.Navigation) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO navigations (parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		navigation.ParentId,
		navigation.Icon,
		navigation.Title,
		navigation.Description,
		navigation.Url,
		navigation.Route,
		navigation.Mode,
		navigation.Urutan,
		navigation.CreatedAt,
		navigation.UpdatedAt,
	)
	return err
}

// Update adalah fungsi untuk memperbarui.
func (r *navigationRepo) Update(navigation *types.Navigation) error {
	stmt, err := r.db.Prepare(`
		UPDATE navigations
		SET parent_id = ?, icon = ?, title = ?, description = ?, url = ?, route = ?, mode = ?, urutan = ?, updated_at = ?
		WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		navigation.ParentId,
		navigation.Icon,
		navigation.Title,
		navigation.Description,
		navigation.Url,
		navigation.Route,
		navigation.Mode,
		navigation.Urutan,
		navigation.UpdatedAt,
		navigation.ID,
	)
	return err
}

// Delete adalah fungsi untuk menghapus.
func (r *navigationRepo) Delete(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM navigations WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// GetAll adalah fungsi untuk mengambil all.
func (r *navigationRepo) GetAll() ([]types.Navigation, error) {
	rows, err := r.db.Query(`
		SELECT id, parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at
		FROM navigations
		ORDER BY COALESCE(parent_id, 0) ASC, urutan ASC, title ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var navigations []types.Navigation
	for rows.Next() {
		var nav types.Navigation
		if err := rows.Scan(
			&nav.ID,
			&nav.ParentId,
			&nav.Icon,
			&nav.Title,
			&nav.Description,
			&nav.Url,
			&nav.Route,
			&nav.Mode,
			&nav.Urutan,
			&nav.CreatedAt,
			&nav.UpdatedAt,
		); err != nil {
			return nil, err
		}
		navigations = append(navigations, nav)
	}
	return navigations, rows.Err()
}

// FindNavigationByID adalah fungsi untuk menemukan navigasi by id.
func (r *navigationRepo) FindNavigationByID(id int) (*types.Navigation, error) {
	row := r.db.QueryRow(`
		SELECT id, parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at
		FROM navigations
		WHERE id = ?`, id)
	return scanNavigation(row)
}

// FindNavigationByParentID adalah fungsi untuk menemukan navigasi by parent id.
func (r *navigationRepo) FindNavigationByParentID(parentID int) ([]*types.Navigation, error) {
	return r.fetchList(`
		SELECT id, parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at
		FROM navigations
		WHERE parent_id = ?
		ORDER BY urutan ASC, title ASC`, parentID)
}

// FindRootNavigations adalah fungsi untuk menemukan root navigations.
func (r *navigationRepo) FindRootNavigations() ([]*types.Navigation, error) {
	return r.fetchList(`
		SELECT id, parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at
		FROM navigations
		WHERE parent_id IS NULL
		ORDER BY urutan ASC, title ASC`)
}

// ListAll adalah fungsi untuk daftar all.
func (r *navigationRepo) ListAll() ([]*types.Navigation, error) {
	return r.fetchList(`
		SELECT id, parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at
		FROM navigations
		ORDER BY COALESCE(parent_id, 0) ASC, urutan ASC, title ASC`)
}

// FindNavigationByUrutan adalah fungsi untuk menemukan navigasi by urutan.
func (r *navigationRepo) FindNavigationByUrutan(urutan int, parentID *int) (*types.Navigation, error) {
	query := `
		SELECT id, parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at
		FROM navigations
		WHERE urutan = ?`
	var row *sql.Row
	if parentID == nil {
		row = r.db.QueryRow(query+` AND parent_id IS NULL`, urutan)
	} else {
		row = r.db.QueryRow(query+` AND parent_id = ?`, urutan, *parentID)
	}
	return scanNavigation(row)
}

// fetchList adalah fungsi untuk fetch daftar.
func (r *navigationRepo) fetchList(query string, args ...any) ([]*types.Navigation, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var navigations []*types.Navigation
	for rows.Next() {
		nav, err := scanNavigation(rows)
		if err != nil {
			return nil, err
		}
		navigations = append(navigations, nav)
	}
	return navigations, rows.Err()
}

type navigationRowScanner interface {
	Scan(dest ...any) error
}

// scanNavigation adalah fungsi untuk scan navigasi.
func scanNavigation(scanner navigationRowScanner) (*types.Navigation, error) {
	var nav types.Navigation
	if err := scanner.Scan(
		&nav.ID,
		&nav.ParentId,
		&nav.Icon,
		&nav.Title,
		&nav.Description,
		&nav.Url,
		&nav.Route,
		&nav.Mode,
		&nav.Urutan,
		&nav.CreatedAt,
		&nav.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &nav, nil
}
