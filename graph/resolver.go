package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	user_controller "github.com/yogtanko/goGraphql/controller/user"
	variable_controller "github.com/yogtanko/goGraphql/controller/variable"
)

type Resolver struct {
	UserController     user_controller.UserController
	VariableController variable_controller.VariableController
}
