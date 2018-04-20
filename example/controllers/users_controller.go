package controllers

import (
	"github.com/niklucky/vodka"
	"github.com/niklucky/vodka/example/modules/users"
)

type Users struct {
	Service users.API
}

type UserValidation struct {
	FindByID struct {
		Params struct {
			id string `required:"true"`
		}
	}
}

func NewUsers(m users.API) *Users {
	return &Users{
		Service: m,
	}
}

func (ctrl *Users) FindByID(ctx *vodka.Context) (interface{}, error) {
	return ctrl.Service.FindByID(ctx.Params.GetString("id"))
}
