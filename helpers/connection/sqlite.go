package connection

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewSQLite connects to SQLite database using sqlx and returns the connection.
// It takes a timeout duration and a database URL as parameters.
// The database URL should be in the format: file:path/to/database.db?query_params
//
// Example:
//
//	NewSQLite(timeout, "database.db?cache=shared&mode=rwc")
//	NewSQLite(timeout, "file::memory:?cache=shared")
func NewSQLite(timeout time.Duration, dbURL string) *sqlx.DB {
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	db, err := sqlx.Open("sqlite3", dbURL)
	if err != nil {
		panic("Cannot connect to database: " + err.Error())
	}

	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(2 * time.Minute)
	db.SetMaxOpenConns(40)

	if err = db.Ping(); err != nil {
		panic("Database not reachable: " + err.Error())
	}

	return db
}

// NewSQLiteGORM returns a GORM dialector for SQLite using the provided timeout and database URL.
// It uses the NewSQLite function to establish the connection and then creates a GORM dialector with it.
//
// Example usage in main.go:
//
//	  gormDB, err := gorm.Open(
//		  connection.NewSQLiteGORM(timeoutContext, "database.db?cache=shared&mode=rwc"),
//		  // connection.NewSQLiteGORM(timeoutContext, "file::memory:?cache=shared"),
//		  &gorm.Config{...},
//	  )
func NewSQLiteGORM(timeout time.Duration, dbURL string) gorm.Dialector {
	db := NewSQLite(timeout, dbURL)
	return sqlite.New(sqlite.Config{
		Conn: db,
	})
}
