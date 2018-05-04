package repositories

import (
	"github.com/syndicatedb/vodka/builders"
)

const (
	defaultLimit = 100
)

/*
Recorder - Repository interface
*/
type Recorder interface {
	Join(source, key, targetKey, joinType string, fields []string)
	Find(QueryMap, ParamsMap) (interface{}, error)
	FindByID(interface{}) (interface{}, error)
	Create(interface{}) (interface{}, error)
	Delete(QueryMap) (interface{}, error)
	DeleteByID(interface{}) (interface{}, error)
	Update(QueryMap, map[string]interface{}) (interface{}, error)
	// SetMapper - setting mapper to build collection
	SetMapper(mapper Mapper)
	Exec(string) (interface{}, error)
}

/*
QueryMap - simple map key=value
*/
type QueryMap map[string]interface{}

/*
ParamsMap - simple map key=value
*/
type ParamsMap map[string]interface{}

/*
QueryModificator - modification of query
*/
type QueryModificator struct {
	fields  []string
	skip    int
	limit   int
	orderBy []builders.OrderParam
}

// Mapper - mapping interface
type Mapper interface {
	Collection([]interface{}) (interface{}, error)
	Item(interface{}) (interface{}, error)
}
