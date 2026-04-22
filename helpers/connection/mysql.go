package connection

import (
	"context"
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMysql connects to MySQL database using sqlx and returns the connection.
// It takes a timeout duration and a database URL as parameters.
// The database URL should be in the format: mysql://username:password@host:port/database?query_params
//
// Example: mysql://root:password@localhost:3306/go-template?charset=utf8mb4&parseTime=True&loc=Local
func NewMysql(timeout time.Duration, dbURL string) *sqlx.DB {
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// convert dbURL to DSN
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		panic("Invalid database URL: " + err.Error())
	}

	user := parsedURL.User.Username()
	password, hasPassword := parsedURL.User.Password()
	host := parsedURL.Host
	dbName := parsedURL.Path
	query := parsedURL.Query()

	// always parseTime=true
	query.Set("parseTime", "true")

	var dsn string
	if hasPassword {
		dsn = fmt.Sprintf("%s:%s@tcp(%s)%s?%s", user, password, host, dbName, query.Encode())
	} else {
		dsn = fmt.Sprintf("%s@tcp(%s)%s?%s", user, host, dbName, query.Encode())
	}

	db, err := sqlx.Open("mysql", dsn)
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

// NewMysqlGORM returns a GORM dialector for MySQL using the provided timeout and database URL.
// It uses the NewMysql function to establish the connection and then creates a GORM dialector with it.
//
// Example usage in main.go:
//
//	  gormDB, err := gorm.Open(
//		  connection.NewMysqlGORM(timeoutContext, "mysql://root:password@localhost:3306/go-template?charset=utf8mb4&parseTime=True&loc=Local"),
//		  &gorm.Config{...},
//	  )
func NewMysqlGORM(timeout time.Duration, dbURL string) gorm.Dialector {
	db := NewMysql(timeout, dbURL)
	return mysql.New(mysql.Config{
		Conn: db,
	})
}
