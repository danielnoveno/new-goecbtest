/*
   file:           repository/ecbstation.go
   description:    Repositori data untuk ecbstation
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package repository

import (
	"database/sql"
	"go-ecb/app/types"
	"go-ecb/pkg/logging"

	"github.com/go-gorp/gorp"
)

type EcbStationRepository struct {
	dbmap *gorp.DbMap
}

func NewEcbStationRepository(dbmap *gorp.DbMap) *EcbStationRepository {
	return &EcbStationRepository{dbmap: dbmap}
}

func (r *EcbStationRepository) GetEcbStation() ([]types.EcbStation, error) {
	var ecbStations []types.EcbStation
	_, err := r.dbmap.Select(&ecbStations, "SELECT * FROM ecbstations")
	if err != nil {
		logging.Logger().Errorf("Error querying ecbstations: %v", err)
		return nil, err
	}
	return ecbStations, nil
}

func (r *EcbStationRepository) FindEcbStationByIP(ip string) (*types.EcbStation, error) {
	var s types.EcbStation
	err := r.dbmap.SelectOne(&s, "SELECT * FROM ecbstations WHERE ipaddress = ? LIMIT 1", ip)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &s, nil
}

func (r *EcbStationRepository) CreateEcbStation(ecbStation types.EcbStation) (int, error) {
	err := r.dbmap.Insert(&ecbStation)
	if err != nil {
		logging.Logger().Errorf("Error creating ecbstation: %v", err)
		return 0, err
	}
	return ecbStation.ID, nil
}

func (r *EcbStationRepository) UpdateEcbStation(ecbStation types.EcbStation) error {
	_, err := r.dbmap.Update(&ecbStation)
	if err != nil {
		logging.Logger().Errorf("Error updating ecbstation: %v", err)
		return err
	}
	return nil
}

func (r *EcbStationRepository) DeleteEcbStation(id int) error {
	_, err := r.dbmap.Exec("DELETE FROM ecbstations WHERE id = ?", id)
	if err != nil {
		logging.Logger().Errorf("Error deleting ecbstation: %v", err)
		return err
	}
	return nil
}
