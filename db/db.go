/*
   file:           db/db.go
   description:    Setup database untuk db
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package db

import (
	"database/sql"
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


func InitDb() (*gorp.DbMap, error) {
	cfg := configs.LoadConfig()
	dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBAddress + ")/" + cfg.DBName + "?parseTime=true&allowNativePasswords=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	isRPi := configs.IsRaspberryPi()
	if isRPi {
		db.SetMaxOpenConns(3)
		db.SetMaxIdleConns(1)
		db.SetConnMaxLifetime(5 * time.Minute)
		log.Println("DB: Raspberry Pi mode detected, using optimized connection pool (max:3, idle:1)")
	} else {
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
	dbmap.AddTableWithName(types.Theme{}, "themes").SetKeys(true, "ID")
	dbmap.AddTableWithName(types.EcbConfig{}, "ecbconfigs").SetKeys(true, "ID")
	dbmap.AddTableWithName(types.Navigation{}, "navigations").SetKeys(true, "ID")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}

	log.Println("DB: Successfully connected!")
	DB = dbmap
	return dbmap, nil
}
