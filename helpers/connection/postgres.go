package connection

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres connects to PostgreSQL database using sqlx and returns the connection.
// It takes a timeout duration and a database URL as parameters.
// The database URL should be in the format: postgres://username:password@host:port/database?query_params
//
// Example: postgres://user:password@localhost:5432/go-template?sslmode=disable
func NewPostgres(timeout time.Duration, dbURL string) *sqlx.DB {
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	db, err := sqlx.Open("pgx", dbURL)
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

// NewPostgresGORM returns a GORM dialector for PostgreSQL using the provided timeout and database URL.
// It uses the NewPostgres function to establish the connection and then creates a GORM dialector with it.
//
// Example usage in main.go:
//
//	  gormDB, err := gorm.Open(
//		  connection.NewPostgresGORM(timeoutContext, "postgres://user:password@localhost:5432/go-template?sslmode=disable"),
//		  &gorm.Config{...},
//	  )
func NewPostgresGORM(timeout time.Duration, dbURL string) gorm.Dialector {
	db := NewPostgres(timeout, dbURL)
	return postgres.New(postgres.Config{
		Conn: db,
	})
}
