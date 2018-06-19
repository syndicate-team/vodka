package items

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

// Item — struct that describes Item
type Item struct {
	ID        string      `db:"id" uuid:"true" key:"true" json:"id"`
	Name      interface{} `db:"name" json:"name" unique:"true"`
	CreatedAt time.Time   `db:"created_at" json:"createdAt"`
	Amount    float64     `db:"amount" json:"amount" unique:"true"`
	Count     int64       `db:"count" json:"count"`
	Status    string      `db:"status" json:"status"`
}

const source = "public.items"

// New - module constructor
func New(adapter adapters.Adapter) *API {
	var u Item
	repo := repositories.NewPostgres(adapter, source, &u)
	return &API{
		Service: base.NewService(repo),
	}
}
