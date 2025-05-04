package mysqldb

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type MySQLDB struct{}

func (d MySQLDB) Connect() (*sql.DB, error) {
	_ = godotenv.Load(".env")
	MYSQL_USERNAME := os.Getenv("MYSQL_USERNAME")
	MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	MYSQL_DATABASE_NAME := os.Getenv("MYSQL_DATABASE_NAME")
	MYSQL_HOSTNAME := os.Getenv("MYSQL_HOSTNAME")
	MYSQL_PORT := os.Getenv("MYSQL_PORT")

	dsn := "" + MYSQL_USERNAME + ":" + MYSQL_PASSWORD + "@tcp(" + MYSQL_HOSTNAME + ":" + MYSQL_PORT + ")/" + MYSQL_DATABASE_NAME
	// log.Printf("DSN: %s", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
		return nil, err
	}
	return db, nil
}
