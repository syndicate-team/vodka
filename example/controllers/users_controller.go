package controllers

import (
	"github.com/syndicatedb/vodka/base"
	"github.com/syndicatedb/vodka/example/modules/users"
)

// Users - users controller struct
type Users struct {
	base.Controller
}

// NewUsers - users constructors
func NewUsers(m *users.API) *Users {
	return &Users{
		Controller: base.NewController(m),
	}
}
