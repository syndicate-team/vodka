package base

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/niklucky/vodka"

	"github.com/niklucky/vodka/lib/adapters"

	lib "github.com/niklucky/go-lib"
)

const (
	defaultLimit = 100
)

/*
Repository - is responsible for storing/fetching data
*/
type Repository struct {
	Adapter            adapters.Adapter
	Model              interface{}
	Source             string
	AutoID             string
	BuildCollection    func([]interface{}) (interface{}, error)
	BuildItem          func(interface{}) (interface{}, error)
	Debug              bool
	joinedRepositories map[string]joinRepository
}

type joinRepository struct {
	source         string
	model          interface{}
	condition      map[string]interface{}
	conditionValue map[string]interface{}
	joinType       string
}

/*
Repo - Repository interface
*/
type Repo interface {
	Find(Query, Modificator) (interface{}, error)
	Join(string, string, interface{}, map[string]string, string)
	FindByID(interface{}) (interface{}, error)
	Create(interface{}) (interface{}, error)
	Delete(Query) (interface{}, error)
	Update(Query, map[string]interface{}) (interface{}, error)
}

/*
Query - simple map key=value
*/
type Query map[string]interface{}

/*
Modificator - modification of query
*/
type Modificator struct {
	fields  []string
	skip    int
	limit   int
	orderBy []adapters.OrderParam
}

/*
NewRepository - Repository constructor
*/
func NewRepository(adapter adapters.Adapter, source string, model interface{}) Repository {
	var debug bool
	if os.Getenv("DEBUG") != "" {
		debug = true
	}
	return Repository{
		Adapter:            adapter,
		Source:             source,
		Model:              model,
		Debug:              debug,
		joinedRepositories: make(map[string]joinRepository),
	}
}

// Join - joining additional repository for queries
func (ds *Repository) Join(sourceID string, source string, model interface{}, joinType string) {
	ds.joinedRepositories[sourceID] = joinRepository{
		model:    model,
		source:   source,
		joinType: joinType,
	}
}

func (ds *Repository) SetJoinCondition(sourceID, key string, value interface{}) {
	item := ds.joinedRepositories[sourceID]
	if item.condition == nil {
		item.condition = make(map[string]interface{})
	}
	item.condition[key] = value
	ds.joinedRepositories[sourceID] = item
}
func (ds *Repository) SetJoinConditionValue(sourceID, key string, value interface{}) {
	item := ds.joinedRepositories[sourceID]
	if item.conditionValue == nil {
		item.conditionValue = make(map[string]interface{})
	}
	item.conditionValue[key] = value
	ds.joinedRepositories[sourceID] = item
}

/*
Create - save data to Storage with Adapter
*/
func (ds Repository) Create(data interface{}) (interface{}, error) {
	builder := ds.Adapter.Builder()
	builder.Insert(ds.Source).Values(data).ReturnID(ds.AutoID)
	SQL := builder.Build()
	if ds.Debug {
		fmt.Println("Create SQL: ", SQL)
	}
	result, err := ds.Adapter.Exec(SQL)
	if err != nil {
		return nil, err
	}
	if id, err := result.LastInsertId(); err == nil {
		return ds.FindByID(id)
	}
	// 	rows, _ := result.RowsAffected()
	// var id interface{}
	// if err := row.Scan(&id); err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return nil, err
	// }
	// if id != nil {
	// 	return ds.FindByID(id)
	// }
	return data, nil
}

/*
Delete - deleteing from storage by query
*/
func (ds Repository) Delete(q Query) (interface{}, error) {
	builder := ds.Adapter.Builder()
	SQL := builder.Delete().From(ds.Source).Where(q).Build()
	if ds.Debug {
		fmt.Println("Delete SQL: ", SQL)
	}

	rows, err := ds.Adapter.Exec(SQL)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

/*
DeleteByID - deleteing from storage by query
*/
func (ds Repository) DeleteByID(id interface{}) (interface{}, error) {
	builder := ds.Adapter.Builder()
	q := make(map[string]interface{})
	q["id"] = id
	SQL := builder.Delete().From(ds.Source).Where(q).Build()
	if ds.Debug {
		fmt.Println("DeleteByID SQL: ", SQL)
	}
	rows, err := ds.Adapter.Exec(SQL)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

/*
Update - updating item in storage by query and payload
*/
func (ds Repository) Update(q Query, payload map[string]interface{}) (interface{}, error) {
	builder := ds.Adapter.Builder()
	SQL := builder.Update(ds.Source).Set(payload).Where(q).Limit(1, 0).Build()
	if ds.Debug {
		fmt.Println("Update SQL: ", SQL)
	}
	result, err := ds.Adapter.Exec(SQL)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	fmt.Printf("Update Result: %+v\n", id)
	if q[ds.AutoID] != nil {
		return ds.FindByID(q[ds.AutoID])
	}
	return nil, nil
}

/*
Find - Finding data by query (map key=value) and Modificator
Will return Collection
*/
func (ds Repository) Find(query Query, params interface{}) (interface{}, error) {
	rows, err := ds.fetch(query, params)
	if err != nil {
		return nil, err
	}
	result, err := ds.buildCollection(rows)
	if d, ok := result.([]interface{}); ok {
		if len(d) == 0 {
			return make([]int, 0), err
		}
	}
	return result, err
}

/*
FindByID - fetching Object by id. interface{} because id could be string or int
*/
func (ds Repository) FindByID(id interface{}) (interface{}, error) {
	q := make(map[string]interface{})
	q["id"] = id
	data, err := ds.fetch(q, nil)
	if err != nil {
		return nil, err
	}
	if len(data) > 0 {
		return ds.buildItem(data[0])
	}
	return nil, vodka.NewError(404, "not_found", "Item not found")
}

func (ds Repository) fetch(query Query, params interface{}) ([]interface{}, error) {
	qb := ds.Adapter.Builder()
	var fields []string
	mod := parseParams(params)
	if len(mod.fields) == 0 {
		fields = lib.GetStructTags(ds.Model, "db", true)
	} else {
		fields = mod.fields
	}
	if mod.limit == 0 {
		mod.limit = defaultLimit
	}
	qb.Select(fields).
		From(ds.Source).
		Where(query).
		Limit(mod.limit, mod.skip)

	if len(ds.joinedRepositories) > 0 {
		fmt.Printf("Join: %+v\n", ds.joinedRepositories)
		for sourceID, j := range ds.joinedRepositories {
			var on []adapters.JoinParamOn
			if j.condition != nil {
				for key, v := range j.condition {
					on = append(on, adapters.JoinParamOn{
						SourceKey: fmt.Sprintf("%v", v),
						JoinKey:   key,
					})
				}
			}
			if j.conditionValue != nil {
				for key, v := range j.conditionValue {
					on = append(on, adapters.JoinParamOn{
						Source:    j.source,
						SourceKey: key,
						JoinValue: v,
					})
				}
			}
			qb.Join(adapters.JoinParam{
				SourceID: sourceID,
				Source:   j.source,
				Fields:   lib.GetStructTags(j.model, "db", true),
				Type:     j.joinType,
				On:       on,
			})
		}
	}

	if len(mod.orderBy) > 0 {
		for _, o := range mod.orderBy {
			qb.Order(o)
		}
	}

	SQL := qb.Build()
	if ds.Debug {
		fmt.Println("Fetch SQL: ", SQL)
	}
	rows, err := ds.Adapter.Query(SQL)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer rows.Close()
	return ds.buildResult(rows)
}

func (ds *Repository) buildResult(rows *sql.Rows) ([]interface{}, error) {
	var result []interface{}
	i := 0
	cols, _ := rows.Columns()
	dest := make([]interface{}, len(cols))
	rawResult := make([]interface{}, len(cols))

	for c := range cols {
		dest[c] = &rawResult[c]
	}

	for rows.Next() {
		data := make(map[string]interface{})
		i++
		if err := rows.Scan(dest...); err != nil {
			fmt.Println("Error: ", err)
			return nil, err
		}
		for key, v := range cols {
			if a, ok := rawResult[key].([]byte); ok == true {
				// data[v] = string(a)
				f, e := strconv.ParseFloat(string(a), 64)
				if e != nil {
					data[v] = string(a)
				} else {
					data[v] = f
				}
			} else {
				data[v] = rawResult[key]
			}
		}
		result = append(result, data)
	}
	return result, nil
}

func (ds *Repository) buildCollection(data []interface{}) (interface{}, error) {
	if ds.BuildCollection != nil {
		return ds.BuildCollection(data)
	}

	return data, nil
}

func (ds *Repository) buildItem(data interface{}) (interface{}, error) {
	if ds.BuildItem != nil {
		return ds.BuildItem(data)
	}
	return data, nil
}

func parseParams(params interface{}) (m Modificator) {
	if params == nil {
		return
	}
	if p, ok := params.(map[string]interface{}); ok {
		if p["fields"] != nil {
			m.fields = p["fields"].([]string)
		}
		if p["skip"] != nil {
			m.skip = p["skip"].(int)
		}
		if p["limit"] != nil {
			m.limit = p["limit"].(int)
		}
		if p["orderBy"] != nil {
			var orderParams adapters.OrderParam
			var orderParamsArr []adapters.OrderParam

			orderParams.OrderBy = p["orderBy"].(string)

			if p["order"] == "asc" {
				orderParams.Asc = true
			} else {
				orderParams.Desc = true
			}

			orderParamsArr = append(orderParamsArr, orderParams)
			m.orderBy = orderParamsArr
		}
	}
	return
}
