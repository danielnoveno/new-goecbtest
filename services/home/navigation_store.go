/*
    file:           services/home/navigation_store.go
    description:    Layanan beranda untuk navigation store
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package home

import (
	"database/sql"
	"fmt"
	"go-ecb/app/types"
)

type NavigationStore struct {
	db *sql.DB
}

// NewNavigationStore adalah fungsi untuk baru navigasi store.
func NewNavigationStore(db *sql.DB) *NavigationStore {
	return &NavigationStore{db: db}
}

// FindNavigationByID adalah fungsi untuk menemukan navigasi by id.
func (s *NavigationStore) FindNavigationByID(id int) (*types.Navigation, error) {
	query := `SELECT id, parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at FROM navigations WHERE id = ?`
	row := s.db.QueryRow(query, id)

	nav := &types.Navigation{}
	err := row.Scan(
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
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("navigation not found")
		}
		return nil, err
	}
	return nav, nil
}

// FindNavigationByParentID adalah fungsi untuk menemukan navigasi by parent id.
func (s *NavigationStore) FindNavigationByParentID(parentID int) ([]*types.Navigation, error) {
	query := `SELECT id, parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at FROM navigations WHERE parent_id = ? ORDER BY urutan ASC`
	rows, err := s.db.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var navigations []*types.Navigation
	for rows.Next() {
		nav := &types.Navigation{}
		err := rows.Scan(
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
		)
		if err != nil {
			return nil, err
		}
		navigations = append(navigations, nav)
	}
	return navigations, nil
}

// FindRootNavigations adalah fungsi untuk menemukan root navigations.
func (s *NavigationStore) FindRootNavigations() ([]*types.Navigation, error) {
	query := `SELECT id, parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at FROM navigations WHERE parent_id IS NULL ORDER BY urutan ASC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var navigations []*types.Navigation
	for rows.Next() {
		nav := &types.Navigation{}
		err := rows.Scan(
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
		)
		if err != nil {
			return nil, err
		}
		navigations = append(navigations, nav)
	}
	return navigations, nil
}

// FindNavigationByUrutan adalah fungsi untuk menemukan navigasi by urutan.
func (s *NavigationStore) FindNavigationByUrutan(urutan int, parentID *int) (*types.Navigation, error) {
	var query string
	var row *sql.Row
	if parentID == nil {
		query = `SELECT id, parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at FROM navigations WHERE urutan = ? AND parent_id IS NULL`
		row = s.db.QueryRow(query, urutan)
	} else {
		query = `SELECT id, parent_id, icon, title, description, url, route, mode, urutan, created_at, updated_at FROM navigations WHERE urutan = ? AND parent_id = ?`
		row = s.db.QueryRow(query, urutan, *parentID)
	}

	nav := &types.Navigation{}
	err := row.Scan(
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
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("navigation not found")
		}
		return nil, err
	}
	return nav, nil
}
