package controllers

import (
	"github.com/niklucky/vodka/base"
	"github.com/niklucky/vodka/example/modules/users"
)

// Users - users controller struct
type Users struct {
	base.Controller
}

// UserValidation - users controller struct
type UserValidation struct {
	FindByID struct {
		Params struct {
			id string `required:"true"`
		}
	}
	Find struct {
		Query struct {
			id        string
			name      string
			createdAt string
		}
	}
	Create struct {
		Body struct {
			name   string  `required:"true"`
			count  int64   `required:"true"`
			amount float64 `required:"true"`
		}
	}
}

// NewUsers - users constructors
func NewUsers(m *users.API) *Users {
	return &Users{
		Controller: base.NewController(m),
	}
}
