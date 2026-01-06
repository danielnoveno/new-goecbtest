/*
   file:           repository/ecbstation.go
   description:    Repositori data untuk ecbstation
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package repository

import (
	"go-ecb/app/types"
	"go-ecb/db"
	"log"
)

// GetEcbStation adalah fungsi untuk mengambil ecb stasiun.
func GetEcbStation() ([]types.EcbStation, error) {
	var ecbStations []types.EcbStation
	rows, err := db.Query("SELECT id, ipaddress, location, mode, linetype, lineids, lineactive, theme, tacktime, workcenters, status, created_at, updated_at FROM ecbstations")
	if err != nil {
		log.Printf("Error querying ecbstations: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ecbStation types.EcbStation
		if err := rows.Scan(&ecbStation.ID, &ecbStation.Ipaddress, &ecbStation.Location, &ecbStation.Mode, &ecbStation.Linetype, &ecbStation.Lineids, &ecbStation.Lineactive, &ecbStation.Theme, &ecbStation.Tacktime, &ecbStation.Workcenters, &ecbStation.Status, &ecbStation.CreatedAt, &ecbStation.UpdatedAt); err != nil {
			log.Printf("Error scanning ecbstation row: %v", err)
			return nil, err
		}
		ecbStations = append(ecbStations, ecbStation)
	}

	return ecbStations, nil
}

// CreateEcbStation adalah fungsi untuk membuat ecb stasiun.
func CreateEcbStation(ecbStation types.EcbStation) (int, error) {
	result, err := db.Exec("INSERT INTO ecbstations (ipaddress, location, mode, linetype, lineids, lineactive, theme, tacktime, workcenters, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		ecbStation.Ipaddress, ecbStation.Location, ecbStation.Mode, ecbStation.Linetype, ecbStation.Lineids, ecbStation.Lineactive, ecbStation.Theme, ecbStation.Tacktime, ecbStation.Workcenters, ecbStation.Status)
	if err != nil {
		log.Printf("Error creating ecbstation: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for ecbstation: %v", err)
		return 0, err
	}

	return int(id), nil
}

// UpdateEcbStation adalah fungsi untuk memperbarui ecb stasiun.
func UpdateEcbStation(ecbStation types.EcbStation) error {
	_, err := db.Exec("UPDATE ecbstations SET ipaddress = ?, location = ?, mode = ?, linetype = ?, lineids = ?, lineactive = ?, theme = ?, tacktime = ?, workcenters = ?, status = ? WHERE id = ?",
		ecbStation.Ipaddress, ecbStation.Location, ecbStation.Mode, ecbStation.Linetype, ecbStation.Lineids, ecbStation.Lineactive, ecbStation.Theme, ecbStation.Tacktime, ecbStation.Workcenters, ecbStation.Status, ecbStation.ID)
	if err != nil {
		log.Printf("Error updating ecbstation: %v", err)
		return err
	}
	return nil
}

// DeleteEcbStation adalah fungsi untuk menghapus ecb stasiun.
func DeleteEcbStation(id int) error {
	_, err := db.Exec("DELETE FROM ecbstations WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting ecbstation: %v", err)
		return err
	}
	return nil
}
