package repositories

const (
	defaultLimit = 100
)

/*
Recorder - Repository interface
*/
type Recorder interface {
	Join(string, string, string, string)
	Find(QueryMap, ParamsMap) (interface{}, error)
	FindByID(interface{}) (interface{}, error)
	Create(interface{}) (interface{}, error)
	Delete(QueryMap) (interface{}, error)
	DeleteByID(interface{}) (interface{}, error)
	Update(QueryMap, map[string]interface{}) (interface{}, error)
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
	orderBy map[string]string
}

// Mapper - mapping interface
type Mapper interface {
	Collection([]interface{}) (interface{}, error)
	Item(interface{}) (interface{}, error)
}

type joinRepository struct {
	source         string
	model          interface{}
	condition      map[string]interface{}
	conditionValue map[string]interface{}
	joinType       string
}
