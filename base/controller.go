package base

import (
	"github.com/niklucky/vodka"
)

type Controller interface {
	Find(*vodka.Context) (interface{}, error)
	FindByID(*vodka.Context) (interface{}, error)
	Create(*vodka.Context) (interface{}, error)
}

type ctrl struct {
	Service
}

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
