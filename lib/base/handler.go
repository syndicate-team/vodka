package base

import (
	"strconv"
	"strings"

	"github.com/syndicatedb/vodka"
)

// Handler - base handler
type Handler struct {
	Module IModule
}

// Get - getting data by params
func (h *Handler) Get(ctx *vodka.Context) (interface{}, error) {

	return h.Module.Get(ctx.Params.Map(), getModificator(ctx.Query.Map()))
}

// GetByID - getting data by ID
func (h *Handler) GetByID(ctx *vodka.Context) (interface{}, error) {
	return h.Module.GetByID(ctx.Params.Get("id"))
}

// Create - create new event
func (h *Handler) Create(ctx *vodka.Context) (interface{}, error) {
	return h.Module.Create(ctx.Body.Map())
}

// Update - Update entity by query params
func (h *Handler) Update(ctx *vodka.Context) (interface{}, error) {
	query := mergeParamsAndQuery(ctx)
	return h.Module.Update(query, ctx.Body.Map())
}

// Delete - delete entity by query params
func (h *Handler) Delete(ctx *vodka.Context) (interface{}, error) {
	query := mergeParamsAndQuery(ctx)
	_, err := h.Module.Delete(query)
	if err != nil {
		return nil, err
	}
	return vodka.ResponseNoContent{}, nil
}

func mergeParamsAndQuery(ctx *vodka.Context) map[string]interface{} {
	query := make(map[string]interface{})
	for key, v := range ctx.Params.Map() {
		query[key] = v
	}
	for key, v := range ctx.Query.Map() {
		query[key] = v
	}
	return query
}

func getModificator(query map[string]interface{}) Modificator {
	var m Modificator
	if query["__skip"] != nil {
		m.skip, _ = strconv.Atoi(query["__skip"].(string))
	}
	if query["__limit"] != nil {
		m.limit, _ = strconv.Atoi(query["__limit"].(string))
	}
	if query["__fields"] != nil {
		if f, ok := query["__fields"].(string); ok {
			m.fields = strings.Split(f, ",")
		}
	}
	// fmt.Println('M', m)
	return m
}
