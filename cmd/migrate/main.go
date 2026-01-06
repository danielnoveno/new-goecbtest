/*
    file:           cmd/migrate/main.go
    description:    Runner migrasi database untuk main
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package main

import (
	"database/sql"
	"go-ecb/configs"
	"log"
	"os"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// main adalah fungsi untuk utama.
func main() {
	envCfg := configs.LoadConfig()
	cfg := mysqlDriver.Config{
		User:                 envCfg.DBUser,
		Passwd:               envCfg.DBPassword,
		Addr:                 envCfg.DBAddress,
		DBName:               envCfg.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		MultiStatements:      true,
		ParseTime:            true,
	}

	sqlDB, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	driver, err := mysqlMigrate.WithInstance(sqlDB, &mysqlMigrate.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Missing command argument (up, down, down-all, force <version>)")
	}

	cmd := os.Args[1]
	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migration applied successfully.")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Rolled back last migration.")
	default:
		log.Fatalf("Unknown command: %s", cmd)
	}

	v, d, _ := m.Version()
	log.Printf("Current version after command: %d, dirty: %v", v, d)
}
