package controllers

import (
	"github.com/syndicatedb/vodka/base"
	"github.com/syndicatedb/vodka/example/modules/items"
)

// Items - Items controller struct
type Items struct {
	base.Controller
}

// NewItems - Items constructors
func NewItems(m *items.API) *Items {
	return &Items{
		Controller: base.NewController(m),
	}
}
