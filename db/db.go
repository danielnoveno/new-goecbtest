/*
   file:           db/db.go
   description:    Setup database untuk db
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"

	"go-ecb/app/types"
	"go-ecb/configs"
)

var DB *gorp.DbMap
var sqlLogger *log.Logger
var sqlLogEnabled bool

// queryLogger enriches gorp trace output with a duration in milliseconds.
type queryLogger struct {
	logger *log.Logger
}

func (l queryLogger) Printf(format string, v ...interface{}) {
	if len(v) >= 4 {
		prefix, _ := v[0].(string)
		query, qok := v[1].(string)
		argsText, aok := v[2].(string)
		duration, dok := v[3].(time.Duration)
		if qok && aok && dok {
			label := strings.TrimSpace(prefix)
			if label != "" {
				label += " "
			}
			l.logger.Printf("%squery=%s args=%s took=%.2fms", label, query, argsText, duration.Seconds()*1000)
			return
		}
	}
	l.logger.Printf(format, v...)
}

func logSQL(prefix, query string, args []interface{}, duration time.Duration) {
	if !sqlLogEnabled || sqlLogger == nil {
		return
	}
	argsText := fmt.Sprint(args)
	label := strings.TrimSpace(prefix)
	if label != "" {
		label += " "
	}
	sqlLogger.Printf("%squery=%s args=%s took=%.2fms", label, query, argsText, duration.Seconds()*1000)
}

func QueryWith(conn *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := conn.Query(query, args...)
	logSQL("SQL Query", query, args, time.Since(start))
	return rows, err
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return QueryWith(DB.Db, query, args...)
}

func QueryRowWith(conn *sql.DB, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := conn.QueryRow(query, args...)
	logSQL("SQL QueryRow", query, args, time.Since(start))
	return row
}

func QueryRow(query string, args ...interface{}) *sql.Row {
	return QueryRowWith(DB.Db, query, args...)
}

func ExecWith(conn *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	res, err := conn.Exec(query, args...)
	logSQL("SQL Exec", query, args, time.Since(start))
	return res, err
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	return ExecWith(DB.Db, query, args...)
}

// InitDb adalah fungsi untuk inisialisasi db.
func InitDb() (*gorp.DbMap, error) {
	cfg := configs.LoadConfig()
	dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBAddress + ")/" + cfg.DBName + "?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Optimasi connection pool berdasarkan platform
	isRPi := configs.IsRaspberryPi()
	if isRPi {
		// Raspberry Pi: lebih konservatif untuk menghemat memory
		db.SetMaxOpenConns(3)
		db.SetMaxIdleConns(1)
		db.SetConnMaxLifetime(5 * time.Minute)
		log.Println("DB: Raspberry Pi mode detected, using optimized connection pool (max:3, idle:1)")
	} else {
		// Desktop/Server: setting default
		db.SetMaxOpenConns(5)
		db.SetMaxIdleConns(2)
		db.SetConnMaxLifetime(2 * time.Minute)
		log.Println("DB: Desktop mode, using standard connection pool (max:5, idle:2)")
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	appDebug := configs.LoadSimoConfig().AppDebug
	if appDebug {
		logger := log.New(os.Stdout, "", log.LstdFlags)
		sqlLogger = logger
		sqlLogEnabled = true
		dbmap.TraceOn("SQL", queryLogger{logger: logger})
	} else {
		sqlLogEnabled = false
	}

	dbmap.AddTableWithName(types.EcbState{}, "ecbstates").SetKeys(true, "ID")
	dbmap.AddTableWithName(types.EcbData{}, "ecbdatas").SetKeys(true, "ID")
	dbmap.AddTableWithName(types.EcbPo{}, "ecbpos").SetKeys(true, "ID")
	dbmap.AddTableWithName(types.Masterfg{}, "masterfgs").SetKeys(true, "ID")
	dbmap.AddTableWithName(types.Mastersfg{}, "mastersfgs").SetKeys(true, "ID")
	dbmap.AddTableWithName(types.Compressor{}, "compressors").SetKeys(true, "ID")
	dbmap.AddTableWithName(types.Comprefg{}, "comprefgs").SetKeys(true, "ID")
	dbmap.AddTableWithName(types.EcbStation{}, "ecbstations").SetKeys(true, "ID")
	// dbmap.AddTableWithName(types.AccessGroup{}, "access_groups").SetKeys(true, "ID")
	// dbmap.AddTableWithName(types.NavigationAccess{}, "navigation_accesses").SetKeys(true, "ID")
	// dbmap.AddTableWithName(types.HelpCategory{}, "help_categories").SetKeys(true, "ID")
	dbmap.AddTableWithName(types.Theme{}, "themes").SetKeys(true, "ID")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}

	log.Println("DB: Successfully connected!")
	DB = dbmap
	return dbmap, nil
}
