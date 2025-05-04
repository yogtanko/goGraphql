package user_controller

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yogtanko/goGraphql/domain/mysqldb"
	"github.com/yogtanko/goGraphql/graph/model"
)

type UserController struct {
	mysqldb mysqldb.MySQLDB
}

func (u UserController) GetAllUser() ([]*model.User, error) {
	users := []*model.User{}
	db, err := u.mysqldb.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
		return nil, err
	}
	defer db.Close()
	query, err := db.Prepare("SELECT * FROM Users")
	if err != nil {
		log.Fatalf("Error preparing query %v", err)
		return nil, err
	}
	defer query.Close()
	rows, err := query.Query()
	if err != nil {
		log.Fatalf("Error executing query %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Address, &user.Position)
		if err != nil {
			log.Fatalf("Error scanning row %v", err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (u UserController) CreateUser(input *model.NewUser) (*model.User, error) {
	db, err := u.mysqldb.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
		return nil, err
	}
	defer db.Close()
	query, err := db.Prepare("INSERT INTO Users (Name, Address, Position) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatalf("Error preparing query %v", err)
		return nil, err
	}
	defer query.Close()
	result, err := query.Exec(input.Name, input.Address, input.Position)
	if err != nil {
		log.Fatalf("Error executing query %v", err)
		return nil, err
	}
	resultID, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Error getting last insert ID %v", err)
		return nil, err
	}
	user := model.User{
		ID:       fmt.Sprintf("%d", resultID),
		Name:     input.Name,
		Address:  input.Address,
		Position: &input.Position,
	}
	return &user, nil
}
