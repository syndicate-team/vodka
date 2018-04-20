package users

import (
	"github.com/niklucky/vodka/adapters"
	"github.com/niklucky/vodka/base"
	"github.com/niklucky/vodka/repositories"
)

type API interface {
	FindByID(interface{}) (interface{}, error)
}

type api struct {
	base.Service
	repository repositories.Recorder
}

type User struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func New(adapter adapters.Adapter) API {
	repo := repositories.NewPostgres(adapter, "users", User{})
	return &api{
		Service: base.NewService(repo),
	}
}
