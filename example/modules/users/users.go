package users

import (
	"time"

	"github.com/syndicatedb/vodka/adapters"
	"github.com/syndicatedb/vodka/base"
	"github.com/syndicatedb/vodka/repositories"
)

// API — external API for module
type API struct {
	base.Service
}

// User — struct that describes User
type User struct {
	ID          string    `db:"id" uuid:"true" key:"true" json:"id"`
	Name        string    `db:"name" json:"name"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	MobilePhone string    `db:"mobile_phone" json:"mobilePhone"`
}

const source = "public.users"

// New - module constructor
func New(adapter adapters.Adapter) *API {
	var u User
	repo := repositories.NewPostgres(adapter, source, &u)
	// repo.Join("public.statuses", "id", "status_id", "", []string{"name as status_name"})
	return &API{
		Service: base.NewService(repo),
	}
}
