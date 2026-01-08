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
	"go-ecb/pkg/logging"
	"github.com/go-gorp/gorp"
)

type EcbConfigStore interface {
	FindEcbConfigBySectionAndVariable(section, variable string) (*types.EcbConfig, error)
	CreateEcbConfig(config *types.EcbConfig) error
	UpdateEcbConfig(config *types.EcbConfig) error
	FindEcbConfigsBySection(section string) ([]*types.EcbConfig, error)
}

type EcbConfigRepository struct {
	dbmap *gorp.DbMap
}

func NewEcbConfigRepository(dbmap *gorp.DbMap) *EcbConfigRepository {
	return &EcbConfigRepository{dbmap: dbmap}
}

func (r *EcbConfigRepository) FindEcbConfigBySectionAndVariable(section, variable string) (*types.EcbConfig, error) {
	var config types.EcbConfig
	err := r.dbmap.SelectOne(&config, "SELECT * FROM ecbconfigs WHERE section = ? AND variable = ?", section, variable)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ecbconfig not found")
		}
		return nil, err
	}
	return &config, nil
}

func (r *EcbConfigRepository) CreateEcbConfig(config *types.EcbConfig) error {
	err := r.dbmap.Insert(config)
	if err != nil {
		return fmt.Errorf("error creating ecbconfig: %w", err)
	}
	return nil
}

func (r *EcbConfigRepository) UpdateEcbConfig(config *types.EcbConfig) error {
	_, err := r.dbmap.Update(config)
	if err != nil {
		return fmt.Errorf("error updating ecbconfig: %w", err)
	}
	return nil
}

func (r *EcbConfigRepository) FindEcbConfigsBySection(section string) ([]*types.EcbConfig, error) {
	var configs []*types.EcbConfig
	_, err := r.dbmap.Select(&configs, "SELECT * FROM ecbconfigs WHERE section = ?", section)
	if err != nil {
		logging.Logger().Warnf("Error querying EcbConfigs: %v", err)
		return nil, err
	}
	return configs, nil
}
