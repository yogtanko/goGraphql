package variable_controller

import (
	"fmt"
	"log"

	"github.com/yogtanko/goGraphql/domain/mysqldb"
	"github.com/yogtanko/goGraphql/graph/model"
)

type VariableController struct {
	mysqldb mysqldb.MySQLDB
}

func (v VariableController) GetAllVariable() ([]*model.Variable, error) {
	panic(fmt.Errorf("not implemented: GetAllVariable"))
}

func (v VariableController) AddVariable(variable *model.AddVariable) (*model.Variable, error) {
	db, err := v.mysqldb.Connect()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()
	row, err := tx.Query(`
		SELECT MAX(StartDate) FROM Variables
		WHERE 
			Variables.Position = ? 
			AND Variables.Name = ? 
			AND StartDate = ?;
			`, variable.Position, variable.Name, variable.StartDate)
	if err != nil {
		log.Fatalf("Error preparing query 2 %v", err)
		return nil, err
	}
	defer row.Close()
	var startDate string
	for row.Next() {
		if err := row.Scan(&startDate); err != nil {
			log.Fatalf("Error scanning row %v", err)
			return nil, err
		}
	}
	_, err = tx.Exec(`
	UPDATE Variables
	SET EndDate = DATE_SUB(?, INTERVAL 1 SECOND)
	WHERE 
		Variables.Position = ? 
		AND Variables.Name = ? 
		AND StartDate = ?;
	`, variable.StartDate, variable.Position, variable.Name, startDate)
	if err != nil {
		log.Fatalf("Error Exec query 3 %v", err)
		return nil, err
	}
	query, err := tx.Prepare(
		`INSERT INTO Variables SET Name=?, Type=?, Description=? , Position=?, StartDate=?`,
	)
	if err != nil {
		log.Fatalf("Error preparing query 1 %v", err)
		return nil, err
	} else {
		log.Printf("Query: %v", query)
	}
	defer query.Close()
	result, err := query.Exec(variable.Name, variable.Type, variable.Description, variable.Position, variable.StartDate)
	if err != nil {
		log.Fatalf("Error executing query %v", err)
		return nil, err
	} else {
		log.Printf("Result: %v", result)
	}
	resultID, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Error getting last insert ID %v", err)
		return nil, err
	} else {
		log.Printf("Result ID: %v", resultID)
	}
	variableModel := model.Variable{
		ID:          fmt.Sprintf("%d", resultID),
		Name:        variable.Name,
		Description: variable.Description,
		Type:        variable.Type,
		Position:    variable.Position,
		StartDate:   variable.StartDate,
		EndDate:     "9999-12-31",
	}
	for i := 0; i < len(variable.Formula); i++ {
		formula := variable.Formula[i]
		query, err := tx.Prepare(`INSERT INTO VariableFormulas SET VariableID=?,Step=?, TokenType=?, Token=?`)
		if err != nil {
			log.Fatalf("Error preparing query 4 %v", err)
			return nil, err
		}
		defer query.Close()
		_, err = query.Exec(resultID, i+1, formula.TokenType, formula.Token)
		if err != nil {
			log.Fatalf("Error executing query %v", err)
			return nil, err
		}
	}
	tx.Commit()
	return &variableModel, nil
}

func (v VariableController) UpdateVariable(variable *model.AddVariable) (*model.Variable, error) {
	panic(fmt.Errorf("not implemented: UpdateVariable"))
}
