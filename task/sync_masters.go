/*
   file:           task/sync_masters.go
   description:    Penjadwal sinkronisasi master data SIMO
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package task

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"go-ecb/configs"

	_ "github.com/go-sql-driver/mysql"
)

// GetEcbstation hanya menarik master ecbstations (setara simo:ecbstation).
func GetEcbstation() { runSingleSync("[GetEcbstation]", syncEcbstation) }

// GetMasterfg hanya menarik masterfgs (setara simo:masterfg).
func GetMasterfg() { runSingleSync("[GetMasterfg]", syncMasterfg) }

// GetMastersfg hanya menarik mastersfgs (setara simo:mastersfg).
func GetMastersfg() { runSingleSync("[GetMastersfg]", syncMastersfg) }

// GetCompressor hanya menarik compressors (setara simo:compressor).
func GetCompressor() { runSingleSync("[GetCompressor]", syncCompressor) }

// GetComprefg hanya menarik comprefgs (setara simo:comprefg).
func GetComprefg() { runSingleSync("[GetComprefg]", syncComprefg) }

// GetAllMasters menjalankan seluruh sync master secara berurutan (setara simo:getallmasters Laravel; ecbpo tetap tidak disertakan).
func GetAllMasters() {
	start := time.Now()
	log.Println("[GetAllMasters] Starting master data sync...")

	cfg := configs.LoadConfig()
	simoDSN := buildDSN(
		cfg.SimoprdUser,
		cfg.SimoprdPassword,
		fmt.Sprintf("%s:%s", cfg.SimoprdHost, cfg.SimoprdPort),
		cfg.SimoprdDatabase,
	)
	localDSN := buildDSN(
		cfg.DBUser,
		cfg.DBPassword,
		fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		cfg.DBName,
	)

	simoDB, err := sql.Open("mysql", simoDSN)
	if err != nil {
		log.Printf("[GetAllMasters] Cannot connect SIMO: %v", err)
		return
	}
	defer simoDB.Close()

	localDB, err := sql.Open("mysql", localDSN)
	if err != nil {
		log.Printf("[GetAllMasters] Cannot connect local DB: %v", err)
		return
	}
	defer localDB.Close()

	// Ikuti urutan Laravel: ecbstations -> masterfgs -> mastersfgs -> compressors -> comprefgs (ecbpo dikosongkan).
	syncEcbstation(simoDB, localDB)
	syncMasterfg(simoDB, localDB)
	syncMastersfg(simoDB, localDB)
	syncCompressor(simoDB, localDB)
	syncComprefg(simoDB, localDB)

	log.Printf("[GetAllMasters] All master data synced in %s.\n", time.Since(start).Round(time.Second))
}

// runSingleSync membuka koneksi DB dan menjalankan satu fungsi sinkron master.
func runSingleSync(tag string, fn func(*sql.DB, *sql.DB)) {
	start := time.Now()
	log.Println(tag, "Starting sync...")

	cfg := configs.LoadConfig()
	simoDSN := buildDSN(
		cfg.SimoprdUser,
		cfg.SimoprdPassword,
		fmt.Sprintf("%s:%s", cfg.SimoprdHost, cfg.SimoprdPort),
		cfg.SimoprdDatabase,
	)
	localDSN := buildDSN(
		cfg.DBUser,
		cfg.DBPassword,
		fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		cfg.DBName,
	)

	simoDB, err := sql.Open("mysql", simoDSN)
	if err != nil {
		log.Printf("%s Cannot connect SIMO: %v", tag, err)
		return
	}
	defer simoDB.Close()

	localDB, err := sql.Open("mysql", localDSN)
	if err != nil {
		log.Printf("%s Cannot connect local DB: %v", tag, err)
		return
	}
	defer localDB.Close()

	fn(simoDB, localDB)

	log.Printf("%s Done in %s.\n", tag, time.Since(start).Round(time.Second))
}

// syncEcbstation adalah fungsi untuk sync ecbstation.
func syncEcbstation(simo, local *sql.DB) {
	log.Println("[syncEcbstation] Syncing ecbstations...")
	rows, err := simo.Query("SELECT id, ipaddress, location, mode, linetype, lineids, lineactive, ecbstate, theme, tacktime, workcenters, status FROM ecbstations")
	if err != nil {
		log.Println("[syncEcbstation]", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, ipaddress, location, mode, linetype, lineids, lineactive, ecbstate, theme, tacktime, workcenters, status sql.NullString
		rows.Scan(&id, &ipaddress, &location, &mode, &linetype, &lineids, &lineactive, &ecbstate, &theme, &tacktime, &workcenters, &status)
		_, err := local.Exec(`
			INSERT INTO ecbstations (id, ipaddress, location, mode, linetype, lineids, lineactive, ecbstate, theme, tacktime, workcenters, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE ipaddress=VALUES(ipaddress), location=VALUES(location), mode=VALUES(mode), linetype=VALUES(linetype), lineids=VALUES(lineids), lineactive=VALUES(lineactive), ecbstate=VALUES(ecbstate), theme=VALUES(theme), tacktime=VALUES(tacktime), workcenters=VALUES(workcenters), status=VALUES(status)`,
			id, ipaddress, location, mode, linetype, lineids, lineactive, ecbstate, theme, tacktime, workcenters, status)
		if err != nil {
			log.Println("[syncEcbstation] Insert error:", err)
		}
	}
	log.Println("[syncEcbstation] Done.")
}

// syncMasterfg adalah fungsi untuk sync masterfg.
func syncMasterfg(simo, local *sql.DB) {
	log.Println("[syncMasterfg] Syncing masterfgs...")
	rows, err := simo.Query("SELECT id, fgtype, lotinv, mattype, matdesc, aging_tipes_id, kdbar, warna, attrib, category FROM masterfgs")
	if err != nil {
		log.Println("[syncMasterfg]", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, fgtype, lotinv, mattype, matdesc, aging, kdbar, warna, attrib, category sql.NullString
		rows.Scan(&id, &fgtype, &lotinv, &mattype, &matdesc, &aging, &kdbar, &warna, &attrib, &category)
		_, err := local.Exec(`
			INSERT INTO masterfgs (id, fgtype, lotinv, mattype, matdesc, aging_tipes_id, kdbar, warna, attrib, category)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE mattype=?, matdesc=?, aging_tipes_id=?, kdbar=?, warna=?, attrib=?, category=?`,
			id, fgtype, lotinv, mattype, matdesc, aging, kdbar, warna, attrib, category,
			mattype, matdesc, aging, kdbar, warna, attrib, category)
		if err != nil {
			log.Println("[syncMasterfg] Insert error:", err)
		}
	}
	log.Println("[syncMasterfg] Done.")
}

// syncMastersfg adalah fungsi untuk sync mastersfg.
func syncMastersfg(simo, local *sql.DB) {
	log.Println("[syncMastersfg] Syncing mastersfgs...")
	rows, err := simo.Query("SELECT id, plant, mattype, sfgtype, matdesc, sfgdesc FROM mastersfgs")
	if err != nil {
		log.Println("[syncMastersfg]", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, plant, mattype, sfgtype, matdesc, sfgdesc sql.NullString
		rows.Scan(&id, &plant, &mattype, &sfgtype, &matdesc, &sfgdesc)
		_, err := local.Exec(`
			INSERT INTO mastersfgs (id, plant, mattype, sfgtype, matdesc, sfgdesc)
			VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE matdesc=?, sfgdesc=?`,
			id, plant, mattype, sfgtype, matdesc, sfgdesc,
			matdesc, sfgdesc)
		if err != nil {
			log.Println("[syncMastersfg] Insert error:", err)
		}
	}
	log.Println("[syncMastersfg] Done.")
}

// syncCompressor adalah fungsi untuk sync kompresor.
func syncCompressor(simo, local *sql.DB) {
	log.Println("[syncCompressor] Syncing compressors...")
	rows, err := simo.Query("SELECT id, ctype, merk, type, itemcode, force_scan, familycode, status FROM compressors")
	if err != nil {
		log.Println("[syncCompressor]", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, ctype, merk, t, itemcode, forceScan, family, status sql.NullString
		rows.Scan(&id, &ctype, &merk, &t, &itemcode, &forceScan, &family, &status)
		_, err := local.Exec(`
			INSERT INTO compressors (id, ctype, merk, type, itemcode, force_scan, familycode, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE ctype=?, merk=?, type=?, itemcode=?, force_scan=?, familycode=?, status=?`,
			id, ctype, merk, t, itemcode, forceScan, family, status,
			ctype, merk, t, itemcode, forceScan, family, status)
		if err != nil {
			log.Println("[syncCompressor] Insert error:", err)
		}
	}
	log.Println("[syncCompressor] Done.")
}

// syncComprefg adalah fungsi untuk sync comprefg.
func syncComprefg(simo, local *sql.DB) {
	log.Println("[syncComprefg] Syncing comprefgs...")
	rows, err := simo.Query("SELECT id, ctype, barcode, status FROM comprefgs")
	if err != nil {
		log.Println("[syncComprefg]", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, ctype, barcode, status sql.NullString
		rows.Scan(&id, &ctype, &barcode, &status)
		_, err := local.Exec(`
			INSERT INTO comprefgs (id, ctype, barcode, status)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE ctype=?, barcode=?, status=?`,
			id, ctype, barcode, status,
			ctype, barcode, status)
		if err != nil {
			log.Println("[syncComprefg] Insert error:", err)
		}
	}
	log.Println("[syncComprefg] Done.")
}
