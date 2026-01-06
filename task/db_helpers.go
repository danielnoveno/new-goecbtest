/*
   file:           task/db_helpers.go
   description:    Shared helpers for database access inside task package
*/

package task

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// buildDSN constructs a MySQL DSN with common timeout helpers.
func buildDSN(user, pass, addr, db string) string {
	base := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, addr, db)
	return base + "?parseTime=true&timeout=5s&readTimeout=5s&writeTimeout=5s"
}

// pingDSN opens and pings a MySQL DSN to ensure reachability.
func pingDSN(dsn string) error {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return conn.PingContext(ctx)
}
