/*
   file:           task/sync_ecbpo.go
   description:    Sinkronisasi ecbpo antara SIMO dan database lokal
*/

package task

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"go-ecb/configs"
)

// SyncEcbPo menyinkronkan catatan PO dari SIMO ke database lokal secara berkala.
func SyncEcbPo() {
	log.Println("[SyncEcbPo] Start syncing ecbpos...")

	cfg := configs.LoadConfig()
	localDB, err := sql.Open("mysql", buildDSN(
		cfg.DBUser,
		cfg.DBPassword,
		fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		cfg.DBName,
	))
	if err != nil {
		log.Println("[SyncEcbPo] local DB open:", err)
		return
	}
	defer localDB.Close()

	simoDB, err := sql.Open("mysql", buildDSN(
		cfg.SimoprdUser,
		cfg.SimoprdPassword,
		fmt.Sprintf("%s:%s", cfg.SimoprdHost, cfg.SimoprdPort),
		cfg.SimoprdDatabase,
	))
	if err != nil {
		log.Println("[SyncEcbPo] simoprd open:", err)
		return
	}
	defer simoDB.Close()

	if err := localDB.Ping(); err != nil {
		log.Println("[SyncEcbPo] local ping:", err)
		return
	}
	if err := simoDB.Ping(); err != nil {
		log.Println("[SyncEcbPo] simo ping:", err)
		return
	}

	syncEcbPoRecords(simoDB, localDB)
}

// syncEcbPoRecords melakukan workcenter-aware sync tanpa membuka koneksi lagi.
func syncEcbPoRecords(simo, local *sql.DB) {
	log.Println("[syncEcbPoRecords] Syncing ecbpos...")
	simoCfg := configs.LoadSimoConfig()
	workcenters := resolveWorkcenters(local, simoCfg)
	if len(workcenters) == 0 {
		log.Println("[syncEcbPoRecords] no workcenters resolved, skip")
		return
	}

	placeholder := make([]string, 0, len(workcenters))
	args := make([]interface{}, 0, len(workcenters))
	for _, wc := range workcenters {
		placeholder = append(placeholder, "?")
		args = append(args, strings.TrimSpace(wc))
	}

	query := fmt.Sprintf(
		"SELECT id, po, sn, ctype, updated_by, workcenter FROM ecbpos WHERE status != 'scanned' AND workcenter IN (%s)",
		strings.Join(placeholder, ","),
	)
	rows, err := simo.Query(query, args...)
	if err != nil {
		log.Println("[syncEcbPoRecords] select simo:", err)
		return
	}
	defer rows.Close()

	stmt, err := local.Prepare(`
		INSERT INTO ecbpos (id, workcenter, po, sn, ctype, updated_by, status, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, '', ?)
		ON DUPLICATE KEY UPDATE workcenter=VALUES(workcenter), po=VALUES(po), sn=VALUES(sn), ctype=VALUES(ctype), updated_by=VALUES(updated_by), updated_at=VALUES(updated_at)
	`)
	if err != nil {
		log.Println("[syncEcbPoRecords] prepare local insert:", err)
		return
	}
	defer stmt.Close()

	now := time.Now()
	imported := 0
	for rows.Next() {
		var id int
		var po, sn, ctype, workcenter string
		var updatedBy sql.NullInt64
		if err := rows.Scan(&id, &po, &sn, &ctype, &updatedBy, &workcenter); err != nil {
			log.Println("[syncEcbPoRecords] scan:", err)
			continue
		}
		if _, err := stmt.Exec(id, workcenter, po, sn, ctype, updatedBy.Int64, now); err != nil {
			log.Println("[syncEcbPoRecords] insert local:", err)
			continue
		}
		imported++
	}

	scannedRows, err := local.Query("SELECT id FROM ecbpos WHERE status = 'scanned'")
	if err != nil {
		log.Println("[syncEcbPoRecords] select scanned local:", err)
		return
	}
	defer scannedRows.Close()

	var scannedIDs []int
	for scannedRows.Next() {
		var id int
		if err := scannedRows.Scan(&id); err == nil {
			scannedIDs = append(scannedIDs, id)
		}
	}
	if len(scannedIDs) > 0 {
		stmtUpdate, err := simo.Prepare("UPDATE ecbpos SET status = 'scanned' WHERE id = ?")
		if err != nil {
			log.Println("[syncEcbPoRecords] prepare update simo:", err)
			return
		}
		defer stmtUpdate.Close()
		for _, id := range scannedIDs {
			if _, err := stmtUpdate.Exec(id); err != nil {
				log.Println("[syncEcbPoRecords] update simo scanned:", err)
			}
		}

		if _, err := local.Exec("DELETE FROM ecbpos WHERE status = 'scanned'"); err != nil {
			log.Println("[syncEcbPoRecords] delete scanned local:", err)
		}
	}

	log.Printf("[syncEcbPoRecords] Done. imported=%d pushedScanned=%d\n", imported, len(scannedIDs))
}
