/*
   file:           task/push_ecb_data.go
   description:    Penjadwal push ecbdatas ke SIMO dan bserv
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package task

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go-ecb/app/types"
	"go-ecb/configs"

	_ "github.com/go-sql-driver/mysql"
)

// PostEcbData adalah fungsi untuk mengirim ecb data.
func PostEcbData() {
	fmt.Println("[PostEcbData] Start pushing pending ecbdatas...")
	start := time.Now()

	cfg := configs.LoadConfig()
	localDB, err := sql.Open("mysql", buildDSN(
		cfg.DBUser,
		cfg.DBPassword,
		fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		cfg.DBName,
	))
	if err != nil {
		log.Println("[PostEcbData] local DB open:", err)
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
		log.Println("[PostEcbData] simoprd open:", err)
		return
	}
	defer simoDB.Close()

	bservDB, err := sql.Open("mysql", buildDSN(
		cfg.BservUser,
		cfg.BservPassword,
		fmt.Sprintf("%s:%s", cfg.BservHost, cfg.BservPort),
		cfg.BservDatabase,
	))
	if err != nil {
		log.Println("[PostEcbData] bserv open:", err)
		return
	}
	defer bservDB.Close()

	if err := localDB.Ping(); err != nil {
		log.Println("[PostEcbData] local DB ping:", err)
		return
	}
	if err := simoDB.Ping(); err != nil {
		log.Println("[PostEcbData] simoprd ping:", err)
		return
	}
	if err := bservDB.Ping(); err != nil {
		log.Println("[PostEcbData] bserv ping:", err)
		return
	}

	rows, err := localDB.Query(`
		SELECT id, tgl, jam, wc, prdline, ctgr, sn, fgtype, spc, comptype, compcode, po, sendsts
		FROM ecbdatas
		WHERE sendsts = ''
		ORDER BY tgl, jam
		LIMIT 500`)
	if err != nil {
		log.Println("[PostEcbData] select pending:", err)
		return
	}
	defer rows.Close()

	lineconvs := map[string]int{
		"REF A":         1,
		"REF B":         2,
		"REF C":         3,
		"REF D":         4,
		"REF E":         5,
		"REF F":         6,
		"REF G":         7,
		"REF H":         8,
		"REF I":         9,
		"REF J":         10,
		"PWM":           11,
		"PWM 1":         11,
		"PAW":           12,
		"DISPENSER":     14,
		"DISPENSER 1":   14,
		"PAC":           15,
		"AC_INDOOR":     15,
		"AC_INDOOR 1":   15,
		"SAP":           16,
		"PWM 2":         17,
		"AC_INDOOR 2":   18,
		"AC_OUTDOOR 1":  19,
		"AC_OUTDOOR 2":  20,
		"CHEST_FREEZER": 21,
	}

	type bservRow struct {
		mfgpoststs string
		spc        string
		tgl        time.Time
		jam        time.Time
		line       int
		gbj        string
		ctype      string
		cdesc      string
		ccode      string
		tipe       string
		tipemfg    string
		lotinv     string
	}

	var simoInserts []types.EcbData
	var bservInserts []bservRow
	var changedIDs []int

	for rows.Next() {
		var data types.EcbData
		var jamStr string
		if err := rows.Scan(
			&data.ID, &data.Tgl, &jamStr, &data.Wc, &data.Prdline, &data.Ctgr,
			&data.Sn, &data.Fgtype, &data.Spc, &data.Comptype, &data.Compcode,
			&data.Po, &data.Sendsts,
		); err != nil {
			log.Println("[PostEcbData] scan:", err)
			continue
		}
		parsedJam, err := parseTimeOnly(jamStr)
		if err != nil {
			log.Printf("[PostEcbData] skip id=%d: jam invalid %q: %v", data.ID, jamStr, err)
			continue
		}
		data.Jam = time.Date(
			data.Tgl.Year(), data.Tgl.Month(), data.Tgl.Day(),
			parsedJam.Hour(), parsedJam.Minute(), parsedJam.Second(), parsedJam.Nanosecond(),
			data.Tgl.Location(),
		)

		lineNum := lineconvs[data.Prdline]
		if lineNum == 0 {
			if parsed, ok := parseLineNumber(data.Prdline); ok {
				lineNum = parsed
			}
		}
		if lineNum == 0 {
			log.Printf("[PostEcbData] stop at id=%d: undefined prdline %s", data.ID, data.Prdline)
			break
		}

		if len(data.Sn) < 4 {
			log.Printf("[PostEcbData] skip id=%d: sn too short (%s)", data.ID, data.Sn)
			continue
		}

		var lotinv string
		if err := localDB.QueryRow("SELECT lotinv FROM masterfgs WHERE kdbar = ? LIMIT 1", data.Sn[0:4]).Scan(&lotinv); err != nil {
			log.Printf("[PostEcbData] skip id=%d: masterfg not found for %s", data.ID, data.Sn)
			continue
		}

		compDesc := ""
		if strings.HasPrefix(data.Prdline, "REF") {
			var merk, compType string
			if err := localDB.QueryRow("SELECT merk, type FROM compressors WHERE ctype = ? LIMIT 1", data.Comptype).Scan(&merk, &compType); err != nil {
				log.Printf("[PostEcbData] skip id=%d: compressor %s not found", data.ID, data.Comptype)
				continue
			}
			compDesc = fmt.Sprintf("%s %s", merk, compType)
		}

		var exists int
		if err := simoDB.QueryRow(
			"SELECT COUNT(*) FROM ecbdatas WHERE tgl = ? AND jam = ? AND sn = ?",
			data.Tgl, data.Jam, data.Sn,
		).Scan(&exists); err != nil {
			log.Println("[PostEcbData] simo check:", err)
			continue
		}
		if exists == 0 {
			simoInserts = append(simoInserts, data)
		}

		exists = 0
		if err := bservDB.QueryRow(
			"SELECT COUNT(*) FROM ecb WHERE tgl = ? AND jam = ? AND gbj = ?",
			data.Tgl, data.Jam, data.Sn,
		).Scan(&exists); err != nil {
			log.Println("[PostEcbData] bserv check:", err)
			continue
		}
		if exists == 0 {
			bservInserts = append(bservInserts, bservRow{
				mfgpoststs: "ECB_BL",
				spc:        data.Spc,
				tgl:        data.Tgl,
				jam:        data.Jam,
				line:       lineNum,
				gbj:        data.Sn,
				ctype:      data.Comptype,
				cdesc:      compDesc,
				ccode:      data.Compcode,
				tipe:       data.Fgtype,
				tipemfg:    data.Fgtype,
				lotinv:     lotinv,
			})
		}

		changedIDs = append(changedIDs, data.ID)
	}

	isOK1 := len(simoInserts) == 0
	if len(simoInserts) > 0 {
		stmt, err := simoDB.Prepare(`
			INSERT INTO ecbdatas (tgl, jam, wc, prdline, ctgr, sn, fgtype, spc, comptype, compcode, po, sendsts)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			log.Println("[PostEcbData] prepare simo insert:", err)
		} else {
			success := true
			for _, row := range simoInserts {
				if _, err := stmt.Exec(row.Tgl, row.Jam, row.Wc, row.Prdline, row.Ctgr, row.Sn, row.Fgtype, row.Spc, row.Comptype, row.Compcode, row.Po, row.Sendsts); err != nil {
					log.Println("[PostEcbData] insert simo:", err)
					success = false
					break
				}
			}
			stmt.Close()
			if success {
				isOK1 = true
			}
		}
	}

	isOK2 := len(bservInserts) == 0
	if len(bservInserts) > 0 {
		stmt, err := bservDB.Prepare(`
			INSERT INTO ecb (mfgpoststs, spc, tgl, jam, line, gbj, ctype, cdesc, ccode, tipe, tipemfg, lotinv)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			log.Println("[PostEcbData] prepare bserv insert:", err)
		} else {
			success := true
			for _, row := range bservInserts {
				if _, err := stmt.Exec(row.mfgpoststs, row.spc, row.tgl, row.jam, row.line, row.gbj, row.ctype, row.cdesc, row.ccode, row.tipe, row.tipemfg, row.lotinv); err != nil {
					log.Println("[PostEcbData] insert bserv:", err)
					success = false
					break
				}
			}
			stmt.Close()
			if success {
				isOK2 = true
			}
		}
	}

	if isOK1 && isOK2 && len(changedIDs) > 0 {
		placeholders := make([]string, 0, len(changedIDs))
		args := make([]interface{}, 0, len(changedIDs))
		for _, id := range changedIDs {
			placeholders = append(placeholders, "?")
			args = append(args, id)
		}
		query := fmt.Sprintf("UPDATE ecbdatas SET sendsts = 'SENT' WHERE sendsts = '' AND id IN (%s)", strings.Join(placeholders, ","))
		if _, err := localDB.Exec(query, args...); err != nil {
			log.Println("[PostEcbData] update sendsts:", err)
		}
	}

	if _, err := localDB.Exec("DELETE FROM ecbdatas WHERE sendsts = 'SENT'"); err != nil {
		log.Println("[PostEcbData] cleanup SENT rows:", err)
	}

	fmt.Printf("[PostEcbData] Done. inserted simo=%d bserv=%d in %s\n", len(simoInserts), len(bservInserts), time.Since(start).Round(time.Millisecond))
}

func parseTimeOnly(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, fmt.Errorf("empty time")
	}
	layouts := []string{
		"15:04:05",
		"15:04:05.999999",
		"15:04",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse time %q", raw)
}

// parseLineNumber converts prdline like "L1" or "LINE 2" into numeric line index for bserv.
func parseLineNumber(prdline string) (int, bool) {
	p := strings.TrimSpace(strings.ToUpper(prdline))
	if p == "" {
		return 0, false
	}
	// Try direct integer first.
	if n, err := strconv.Atoi(p); err == nil && n > 0 {
		return n, true
	}
	// Strip non-digit prefix (e.g., L1, LINE 2).
	buf := make([]rune, 0, len(p))
	for _, r := range p {
		if r >= '0' && r <= '9' {
			buf = append(buf, r)
		}
	}
	if len(buf) == 0 {
		return 0, false
	}
	if n, err := strconv.Atoi(string(buf)); err == nil && n > 0 {
		return n, true
	}
	return 0, false
}
