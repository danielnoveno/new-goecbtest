/*
   file:           repository/navigation.go
   description:    Repositori data untuk navigation
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package repository

import (
	"database/sql"
	"go-ecb/app/types"

	"github.com/go-gorp/gorp"
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
	dbmap *gorp.DbMap
}

func NewNavigationRepository(dbmap *gorp.DbMap) NavigationStore {
	return &navigationRepo{dbmap: dbmap}
}

func (r *navigationRepo) Insert(navigation *types.Navigation) error {
	return r.dbmap.Insert(navigation)
}

func (r *navigationRepo) Update(navigation *types.Navigation) error {
	_, err := r.dbmap.Update(navigation)
	return err
}

func (r *navigationRepo) Delete(id int) error {
	_, err := r.dbmap.Exec("DELETE FROM navigations WHERE id = ?", id)
	return err
}

func (r *navigationRepo) GetAll() ([]types.Navigation, error) {
	var navigations []types.Navigation
	_, err := r.dbmap.Select(&navigations, `
		SELECT * FROM navigations
		ORDER BY COALESCE(parent_id, 0) ASC, urutan ASC, title ASC`)
	return navigations, err
}

func (r *navigationRepo) FindNavigationByID(id int) (*types.Navigation, error) {
	var nav types.Navigation
	err := r.dbmap.SelectOne(&nav, "SELECT * FROM navigations WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &nav, nil
}

func (r *navigationRepo) FindNavigationByParentID(parentID int) ([]*types.Navigation, error) {
	var navigations []*types.Navigation
	_, err := r.dbmap.Select(&navigations, `
		SELECT * FROM navigations
		WHERE parent_id = ?
		ORDER BY urutan ASC, title ASC`, parentID)
	return navigations, err
}

func (r *navigationRepo) FindRootNavigations() ([]*types.Navigation, error) {
	var navigations []*types.Navigation
	_, err := r.dbmap.Select(&navigations, `
		SELECT * FROM navigations
		WHERE parent_id IS NULL
		ORDER BY urutan ASC, title ASC`)
	return navigations, err
}

func (r *navigationRepo) ListAll() ([]*types.Navigation, error) {
	var navigations []*types.Navigation
	_, err := r.dbmap.Select(&navigations, `
		SELECT * FROM navigations
		ORDER BY COALESCE(parent_id, 0) ASC, urutan ASC, title ASC`)
	return navigations, err
}

func (r *navigationRepo) FindNavigationByUrutan(urutan int, parentID *int) (*types.Navigation, error) {
	var nav types.Navigation
	var err error
	
	query := "SELECT * FROM navigations WHERE urutan = ?"
	if parentID == nil {
		err = r.dbmap.SelectOne(&nav, query + " AND parent_id IS NULL", urutan)
	} else {
		err = r.dbmap.SelectOne(&nav, query + " AND parent_id = ?", urutan, *parentID)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &nav, nil
}
