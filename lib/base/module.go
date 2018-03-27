package base

type Module struct {
	DataService Repo
}

type IModule interface {
	Get(map[string]interface{}, Modificator) (interface{}, error)
	GetByID(interface{}) (interface{}, error)
	Create(map[string]interface{}) (interface{}, error)
	Delete(map[string]interface{}) (interface{}, error)
	Update(Query, map[string]interface{}) (interface{}, error)
}
