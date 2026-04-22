package connection

import (
	"context"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// NewSQLServer connects to SQL Server database using sqlx and returns the connection.
// It takes a timeout duration and a database URL as parameters.
// The database URL should be in the format: sqlserver://username:password@host:port/database?query_params
//
// Example: sqlserver://sa:password@localhost:1433/go-template?encrypt=disable
func NewSQLServer(timeout time.Duration, dbURL string) *sqlx.DB {
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// convert dbURL to DSN
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		panic("Invalid database URL: " + err.Error())
	}

	// database in sqlserver is set in query param, not path, so
	query := parsedURL.Query()
	if dbName := query.Get("database"); dbName == "" {
		query.Set("database", parsedURL.Path[1:]) // remove leading slash
		// remove path since sqlserver driver doesn't use it
		parsedURL.Path = ""
	}

	// restore query to raw query
	parsedURL.RawQuery = query.Encode()

	db, err := sqlx.Open("sqlserver", parsedURL.String())
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

// NewSQLServerGORM returns a GORM dialector for SQL Server using the provided timeout and database URL.
// It uses the NewSQLServer function to establish the connection and then creates a GORM dialector with it.
//
// Example usage in main.go:
//
//	gormDB, err := gorm.Open(
//		connection.NewSQLServerGORM(timeoutContext, "sqlserver://sa:password@localhost:1433/go-template?encrypt=disable"),
//		&gorm.Config{Logger: logger.Default.LogMode(logLevel)},
//	)
func NewSQLServerGORM(timeout time.Duration, dbURL string) gorm.Dialector {
	db := NewSQLServer(timeout, dbURL)
	return sqlserver.New(sqlserver.Config{
		Conn: db.DB,
	})
}
