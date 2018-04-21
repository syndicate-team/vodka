package users

import (
	"time"

	"github.com/niklucky/vodka/adapters"
	"github.com/niklucky/vodka/base"
	"github.com/niklucky/vodka/repositories"
)

type API struct {
	base.Service
}

type User struct {
	ID        string    `db:"id" uuid:"true" key:"true" json:"Id"`
	Name      string    `db:"name" json:"Name"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Amount    float64   `db:"amount"`
	Count     int64     `db:"count"`
}

const source = "users"

func New(adapter adapters.Adapter) *API {
	var u User
	repo := repositories.NewPostgres(adapter, source, &u)
	return &API{
		Service: base.NewService(repo),
	}
}
