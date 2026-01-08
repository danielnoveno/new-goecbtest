/*
   file:           services/ecbcore/service.go
   description:    Backend helpers untuk validasi & penyimpanan alur ECB (tanpa UI)
*/

package ecbcore

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-gorp/gorp"
	"go-ecb/app/types"
	"go-ecb/configs"
	"go-ecb/pkg/logging"
)

type LineConfig struct {
	LineID     string
	Workcenter string
	Category   string
}

type RemoteChecker struct {
	targets []remoteTarget
	timeout time.Duration
}

type remoteTarget struct {
	name        string
	dsn         string
	serialQuery string
	spcQuery    string
}

func NewRemoteChecker(cfg configs.Config) RemoteChecker {
	targets := []remoteTarget{}

	if dsn := buildRemoteDSN(cfg.SimoprdUser, cfg.SimoprdPassword, fmt.Sprintf("%s:%s", cfg.SimoprdHost, cfg.SimoprdPort), cfg.SimoprdDatabase); strings.TrimSpace(dsn) != "" {
		targets = append(targets, remoteTarget{
			name:        "simoprd",
			dsn:         dsn,
			serialQuery: "SELECT 1 FROM ecbdatas WHERE sn = ? LIMIT 1",
			spcQuery:    "SELECT 1 FROM ecbdatas WHERE spc = ? LIMIT 1",
		})
	}

	if dsn := buildRemoteDSN(cfg.BservUser, cfg.BservPassword, fmt.Sprintf("%s:%s", cfg.BservHost, cfg.BservPort), cfg.BservDatabase); strings.TrimSpace(dsn) != "" {
		targets = append(targets, remoteTarget{
			name:        "bserv",
			dsn:         dsn,
			serialQuery: "SELECT 1 FROM ecb WHERE gbj = ? LIMIT 1",
			spcQuery:    "SELECT 1 FROM ecb WHERE spc = ? LIMIT 1",
		})
	}

	return RemoteChecker{
		targets: targets,
		timeout: 2 * time.Second,
	}
}

func buildRemoteDSN(user, pass, addr, db string) string {
	if strings.TrimSpace(user) == "" || strings.TrimSpace(addr) == "" || strings.TrimSpace(db) == "" {
		return ""
	}
	base := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, addr, db)
	return base + "?parseTime=true&timeout=5s&readTimeout=5s&writeTimeout=5s"
}

func (r RemoteChecker) ExistsSerial(value string) (bool, error) {
	return r.exists(value, func(t remoteTarget) string { return t.serialQuery })
}

func (r RemoteChecker) ExistsSpc(value string) (bool, error) {
	return r.exists(value, func(t remoteTarget) string { return t.spcQuery })
}

func (r RemoteChecker) exists(value string, querySelector func(remoteTarget) string) (bool, error) {
	if strings.TrimSpace(value) == "" {
		return false, nil
	}
	timeout := r.timeout
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	for _, target := range r.targets {
		query := querySelector(target)
		if strings.TrimSpace(target.dsn) == "" || strings.TrimSpace(query) == "" {
			continue
		}

		db, err := sql.Open("mysql", target.dsn)
		if err != nil {
			logging.Logger().Warnf("Remote %s open error: %v", target.name, err)
			return false, err
		}

		var dummy int
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		err = db.QueryRowContext(ctx, query, value).Scan(&dummy)
		cancel()
		_ = db.Close()

		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			logging.Logger().Warnf("Remote %s query error: %v", target.name, err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func ValidateSnOnlySerial(dbmap *gorp.DbMap, remote RemoteChecker, sn string) (string, error) {
	sn = strings.ToUpper(strings.TrimSpace(sn))
	if sn == "" {
		return "", fmt.Errorf("input kosong")
	}
	if len(sn) < 4 {
		return "", fmt.Errorf("type produk tidak dikenal")
	}
	if dbmap == nil || dbmap.Db == nil {
		return "", fmt.Errorf("server terputus,scan S/N lagi")
	}
	var masterfg types.Masterfg
	if err := dbmap.SelectOne(&masterfg, "SELECT * FROM masterfgs WHERE kdbar = ? LIMIT 1", sn[0:4]); err != nil {
		return "", fmt.Errorf("type produk tidak dikenal")
	}
	if strings.ToUpper(masterfg.Lotinv) == "IDN0" && len(sn) != 12 {
		return "", fmt.Errorf("panjang barcode salah")
	}
	if ok, err := remote.ExistsSerial(sn); err != nil {
		return "", fmt.Errorf("server terputus,scan S/N lagi")
	} else if ok {
		return "", fmt.Errorf("S/N sudah pernah ECB")
	}
	var snCheck int
	if err := dbmap.SelectOne(&snCheck, "SELECT COUNT(*) FROM ecbdatas WHERE sn = ?", sn); err != nil {
		return "", fmt.Errorf("server terputus,scan S/N lagi")
	}
	if snCheck > 0 {
		return "", fmt.Errorf("S/N sudah pernah ECB")
	}
	fg := strings.ToUpper(masterfg.Fgtype)
	if fg == "" {
		fg = deriveFgType(sn)
	}
	return fg, nil
}

func SaveSnOnly(dbmap *gorp.DbMap, cfg LineConfig, sn, fg string) error {
	if dbmap == nil || dbmap.Db == nil {
		return fmt.Errorf("server terputus,scan S/N lagi")
	}
	if sn == "" {
		return fmt.Errorf("S/N belum terisi")
	}
	now := time.Now()
	_, err := dbmap.Exec(`
		INSERT INTO ecbdatas (tgl, jam, wc, prdline, ctgr, sn, fgtype, spc, comptype, compcode, po, status, sendsts, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		now, now, cfg.Workcenter, cfg.LineID, cfg.Category, sn, fg, "", "", "", "", "", "", now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to save data: %v", err)
	}
	return nil
}

func ValidateSpc(dbmap *gorp.DbMap, remote RemoteChecker, spc string) error {
	spc = strings.ToUpper(strings.TrimSpace(spc))
	if len(spc) != 11 {
		return fmt.Errorf("Panjang SPC salah")
	}
	if dbmap == nil || dbmap.Db == nil {
		return fmt.Errorf("server terputus,scan SPC lagi")
	}
	var count int
	if err := dbmap.SelectOne(&count, "SELECT COUNT(*) FROM ecbdatas WHERE spc = ?", spc); err != nil {
		logging.Logger().Errorf("Failed to validate SPC: %v", err)
		return fmt.Errorf("server terputus,scan SPC lagi")
	}
	if count > 0 {
		return fmt.Errorf("SPC sudah pernah ECB")
	}
	if ok, err := remote.ExistsSpc(spc); err != nil {
		return fmt.Errorf("server terputus,scan SPC lagi")
	} else if ok {
		return fmt.Errorf("SPC sudah pernah ECB")
	}
	return nil
}

func ValidateRefrigSerial(dbmap *gorp.DbMap, remote RemoteChecker, sn string) (string, error) {
	sn = strings.ToUpper(strings.TrimSpace(sn))
	if sn == "" {
		return "", fmt.Errorf("input kosong")
	}
	if len(sn) < 4 {
		return "", fmt.Errorf("type produk tidak dikenal")
	}
	if dbmap == nil || dbmap.Db == nil {
		return "", fmt.Errorf("server terputus,scan S/N lagi")
	}
	var masterfg types.Masterfg
	if err := dbmap.SelectOne(&masterfg, "SELECT * FROM masterfgs WHERE kdbar = ? LIMIT 1", sn[0:4]); err != nil {
		return "", fmt.Errorf("type produk tidak dikenal")
	}
	if strings.ToUpper(masterfg.Lotinv) == "IDN0" && len(sn) != 12 {
		return "", fmt.Errorf("panjang barcode salah")
	}
	if ok, err := remote.ExistsSerial(sn); err != nil {
		return "", fmt.Errorf("server terputus,scan S/N lagi")
	} else if ok {
		return "", fmt.Errorf("S/N sudah pernah ECB")
	}
	var snCheck int
	if err := dbmap.SelectOne(&snCheck, "SELECT COUNT(*) FROM ecbdatas WHERE sn = ?", sn); err != nil {
		logging.Logger().Errorf("SN duplicate check failed: %v", err)
		return "", fmt.Errorf("server terputus,scan S/N lagi")
	}
	if snCheck > 0 {
		return "", fmt.Errorf("S/N sudah pernah ECB")
	}
	fg := strings.ToUpper(masterfg.Fgtype)
	if fg == "" {
		fg = deriveFgType(sn)
	}
	return fg, nil
}

func ValidateCompressorType(dbmap *gorp.DbMap, sn, comptype string, isPO bool) (*types.Compressor, *types.EcbPo, error) {
	comptype = extractCompressorType(comptype)
	if comptype == "" {
		return nil, nil, fmt.Errorf("type compressor belum terisi")
	}
	if dbmap == nil || dbmap.Db == nil {
		return nil, nil, fmt.Errorf("server terputus,scan tipe kompresor lagi")
	}
	var compressor types.Compressor
	if err := dbmap.SelectOne(&compressor, "SELECT * FROM compressors WHERE ctype = ? LIMIT 1", comptype); err != nil {
		return nil, nil, fmt.Errorf("type compressor tak dikenal")
	}
	if len(sn) < 4 {
		return nil, nil, fmt.Errorf("S/N belum terisi")
	}
	prefix := sn[:4]
	if !isPO {
		var comprefg types.Comprefg
		if err := dbmap.SelectOne(&comprefg, "SELECT * FROM comprefgs WHERE ctype = ? AND barcode = ? LIMIT 1", comptype, prefix); err != nil {
			return nil, nil, fmt.Errorf("compressor tak sesuai -%s-%s", comptype, prefix)
		}
		return &compressor, nil, nil
	}

	var ecbpo types.EcbPo
	if err := dbmap.SelectOne(&ecbpo, "SELECT * FROM ecbpos WHERE sn = ? AND ctype = ? LIMIT 1", sn, comptype); err != nil {
		return nil, nil, fmt.Errorf("compressor tak sesuai po-%s-%s", comptype, sn)
	}
	return &compressor, &ecbpo, nil
}

func ValidateCompressorCode(comp *types.Compressor, code string) error {
	if comp == nil {
		return fmt.Errorf("type compressor belum dicek")
	}
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return fmt.Errorf("error input")
	}
	if comp.ForceScan == 1 {
		if strings.Repeat(comp.Ctype, 4) == code {
			return fmt.Errorf("Wajib scan S/N comp asli")
		}
	}
	return nil
}

func SaveRefrig(dbmap *gorp.DbMap, cfg LineConfig, sn, fg, spc, compType, compCode, poVal string, isPO bool, poID *int) error {
	if dbmap == nil || dbmap.Db == nil {
		return fmt.Errorf("server terputus,scan S/N lagi")
	}
	if sn == "" {
		return fmt.Errorf("S/N belum terisi")
	}
	now := time.Now()
	_, err := dbmap.Exec(`
		INSERT INTO ecbdatas (tgl, jam, wc, prdline, ctgr, sn, fgtype, spc, comptype, compcode, po, status, sendsts, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		now, now, cfg.Workcenter, cfg.LineID, cfg.Category, sn, fg, spc, compType, compCode, poVal, "", "", now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to save data: %v", err)
	}
	if isPO && poID != nil {
		if _, err := dbmap.Exec("UPDATE ecbpos SET status = 'scanned' WHERE id = ?", *poID); err != nil {
			logging.Logger().Errorf("failed to update ecbpo status: %v", err)
		}
	}
	return nil
}

func deriveFgType(sn string) string {
	if sn == "" {
		return "-"
	}
	if len(sn) >= 4 {
		return fmt.Sprintf("FG-%s", sn[:4])
	}
	return fmt.Sprintf("FG-%s", sn)
}

func extractCompressorType(raw string) string {
	trimmed := strings.ToUpper(strings.TrimSpace(raw))
	if trimmed == "" {
		return ""
	}
	cutoff := len(trimmed)
	for idx, r := range trimmed {
		if !((r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			cutoff = idx
			break
		}
	}
	clean := trimmed[:cutoff]
	if clean == "" {
		return ""
	}
	if len(clean)%4 == 0 {
		part := clean[:len(clean)/4]
		if len(part) > 0 && strings.Repeat(part, 4) == clean {
			return part
		}
	}
	return clean
}
