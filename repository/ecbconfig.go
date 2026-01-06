/*
    file:           repository/ecbconfig.go
    description:    Repositori data untuk ecbconfig
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package repository

import (
	"database/sql"
	"fmt"
	"go-ecb/app/types"
	"log"
)

type EcbConfigStore interface {
	FindEcbConfigBySectionAndVariable(section, variable string) (*types.EcbConfig, error)
	CreateEcbConfig(config *types.EcbConfig) error
	UpdateEcbConfig(config *types.EcbConfig) error
	FindEcbConfigsBySection(section string) ([]*types.EcbConfig, error)
}

type EcbConfigRepository struct {
	db *sql.DB
}

// NewEcbConfigRepository adalah fungsi untuk baru ecb konfigurasi repository.
func NewEcbConfigRepository(db *sql.DB) *EcbConfigRepository {
	return &EcbConfigRepository{db: db}
}

// FindEcbConfigBySectionAndVariable adalah fungsi untuk menemukan ecb konfigurasi by section and variable.
func (r *EcbConfigRepository) FindEcbConfigBySectionAndVariable(section, variable string) (*types.EcbConfig, error) {
	query := `SELECT id, section, variable, value, ordering, created_at, updated_at FROM ecbconfigs WHERE section = ? AND variable = ?`
	row := r.db.QueryRow(query, section, variable)

	config := &types.EcbConfig{}
	err := row.Scan(
		&config.ID,
		&config.Section,
		&config.Variable,
		&config.Value,
		&config.Ordering,
		&config.CreatedAt,
		&config.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ecbconfig not found")
		}
		return nil, err
	}
	return config, nil
}

// CreateEcbConfig adalah fungsi untuk membuat ecb konfigurasi.
func (r *EcbConfigRepository) CreateEcbConfig(config *types.EcbConfig) error {
	query := `INSERT INTO ecbconfigs (section, variable, value, ordering, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, config.Section, config.Variable, config.Value, config.Ordering, config.CreatedAt, config.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating ecbconfig: %w", err)
	}
	return nil
}

// UpdateEcbConfig adalah fungsi untuk memperbarui ecb konfigurasi.
func (r *EcbConfigRepository) UpdateEcbConfig(config *types.EcbConfig) error {
	query := `UPDATE ecbconfigs SET value = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, config.Value, config.UpdatedAt, config.ID)
	if err != nil {
		return fmt.Errorf("error updating ecbconfig: %w", err)
	}
	return nil
}

// FindEcbConfigsBySection adalah fungsi untuk menemukan ecb configs by section.
func (r *EcbConfigRepository) FindEcbConfigsBySection(section string) ([]*types.EcbConfig, error) {
	query := `SELECT id, section, variable, value, ordering, created_at, updated_at FROM ecbconfigs WHERE section = ?`
	rows, err := r.db.Query(query, section)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*types.EcbConfig
	for rows.Next() {
		config := &types.EcbConfig{}
		err := rows.Scan(
			&config.ID,
			&config.Section,
			&config.Variable,
			&config.Value,
			&config.Ordering,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning EcbConfig row: %v", err)
			continue
		}
		configs = append(configs, config)
	}
	return configs, nil
}
