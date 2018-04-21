package base

import (
	"github.com/niklucky/vodka"
)

type Controller interface {
	Find(*vodka.Context) (interface{}, error)
	FindByID(*vodka.Context) (interface{}, error)
	Create(*vodka.Context) (interface{}, error)
	Update(*vodka.Context) (interface{}, error)
	UpdateByID(*vodka.Context) (interface{}, error)
	DeleteByID(*vodka.Context) (interface{}, error)
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

func (c *ctrl) Update(ctx *vodka.Context) (interface{}, error) {
	return c.Service.Update(ctx.Query.Map(), ctx.Body.Map())
}

func (c *ctrl) UpdateByID(ctx *vodka.Context) (interface{}, error) {
	items, err := c.Service.Update(ctx.Params.Map(), ctx.Body.Map())
	if err != nil {
		return items, err
	}
	return items.([]interface{})[0], nil
}

func (c *ctrl) DeleteByID(ctx *vodka.Context) (interface{}, error) {
	res, err := c.Service.DeleteByID(ctx.Params.GetString("id"))
	if err != nil {
		return res, err
	}
	return vodka.ResponseNoContent{}, nil
}
