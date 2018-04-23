package base

import (
	"github.com/syndicatedb/vodka"
)

/*
Controller - base controller interface that describes base CRUD methods
*/
type Controller interface {
	// Find — finding in DB rows. Returns array
	Find(*vodka.Context) (interface{}, error)
	// FindByID - finding item in DB by key defined in model
	FindByID(*vodka.Context) (interface{}, error)
	// Create - creating new row
	Create(*vodka.Context) (interface{}, error)
	// Update — updating items by Query
	Update(*vodka.Context) (interface{}, error)
	// UpdateByID — updating item by ID in params
	UpdateByID(*vodka.Context) (interface{}, error)
	// DeleteByID - deleting item. Returns no content
	DeleteByID(*vodka.Context) (interface{}, error)
}

// ctrl - struct that contains injected service
type ctrl struct {
	Service
}

// NewController - controller constructor
func NewController(srv Service) Controller {
	return &ctrl{
		Service: srv,
	}
}

func (c *ctrl) FindByID(ctx *vodka.Context) (interface{}, error) {
	return c.Service.FindByID(ctx.Params.GetString("id"))
}

func (c *ctrl) Find(ctx *vodka.Context) (interface{}, error) {
	var params map[string]interface{}
	if p, ok := ctx.Options.Get("params").(map[string]interface{}); ok {
		params = p
	}
	return c.Service.Find(ctx.Query.Map(), params)
}

func (c *ctrl) Create(ctx *vodka.Context) (interface{}, error) {
	return c.Service.Create(ctx.Body.Map())
}

func (c *ctrl) Update(ctx *vodka.Context) (interface{}, error) {
	return c.Service.Update(ctx.Query.Map(), ctx.Body.Map())
}

func (c *ctrl) UpdateByID(ctx *vodka.Context) (interface{}, error) {
	items, err := c.Service.Update(ctx.Params.Map(), ctx.Body.Map())
	if err != nil {
		return items, err
	}
	if item, ok := items.([]interface{}); ok {
		return item[0], nil
	}
	if _, ok := items.([]int); ok {
		return make(map[string]string, 0), nil
	}
	return items, nil
}

func (c *ctrl) DeleteByID(ctx *vodka.Context) (interface{}, error) {
	res, err := c.Service.DeleteByID(ctx.Params.GetString("id"))
	if err != nil {
		return res, err
	}
	return vodka.ResponseNoContent{}, nil
}
