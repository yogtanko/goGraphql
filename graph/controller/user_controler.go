package user_controller

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/yogtanko/goGraphql/graph/model"
)

func GetAllUser() ([]*model.User, error) {
	users := []*model.User{
		{ID: "1", Name: "John Doe", Address: "123 Main St"},
		{ID: "2", Name: "Jane Smith", Address: "456 Elm St"},
	}
	return users, nil
}

func CreateUser(input model.NewUser) (*model.User, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
		return nil, err
	}

	MYSQL_USERNAME := os.Getenv("MYSQL_USERNAME")
	MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	MYSQL_DATABASE_NAME := os.Getenv("MYSQL_DATABASE_NAME")
	MYSQL_HOSTNAME := os.Getenv("MYSQL_HOSTNAME")
	MYSQL_PORT := os.Getenv("MYSQL_PORT")

	dsn := "" + MYSQL_USERNAME + ":" + MYSQL_PASSWORD + "@tcp(" + MYSQL_HOSTNAME + ":" + MYSQL_PORT + ")/" + MYSQL_DATABASE_NAME
	log.Printf("DSN: %s", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
		return nil, err
	}
	defer db.Close()
	query, err := db.Prepare("INSERT INTO users (name, address) VALUES (?, ?)")
	if err != nil {
		log.Fatalf("Error preparing query %v", err)
		return nil, err
	}
	defer query.Close()
	row, err := query.Query(input.Name, input.Address)
	if err != nil {
		log.Fatalf("Error executing query %v", err)
		return nil, err
	}
	defer row.Close()
	user := model.User{
		ID:      "",
		Name:    "",
		Address: "",
	}
	for row.Next() {
		var id int
		var name string
		var address string
		err = row.Scan(&id, &name, &address)
		if err != nil {
			log.Fatalf("Error scanning row %v", err)
			return nil, err
		}
		user.ID = fmt.Sprintf("%d", id)
		user.Name = name
		user.Address = address
	}
	return &user, nil
}
