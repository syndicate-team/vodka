package adapters

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// JoinLeft - constant for SQL query builder
	JoinLeft = "LEFT"
	// JoinRight - constant for SQL query builder
	JoinRight = "RIGHT"
	// JoinInner - constant for SQL query builder
	JoinInner = "INNER"
)

const (
	queryTypeSelect = "SELECT"
	queryTypeInsert = "INSERT"
	queryTypeUpdate = "UPDATE"
	queryTypeDelete = "DELETE"
	tablePrefix     = "t"
	defaultLimit    = 100
)

type parts struct {
	table      string
	fields     []string
	where      map[string]interface{}
	join       []JoinParam
	order      []OrderParam
	limit      int
	offset     int
	insertData interface{}
	returnID   string
}

/*
SQLBuilder - abstract builder for SQL-queries. Now adapted for Postgres
*/
type SQLBuilder struct {
	queryType string
	parts     parts
	sources   map[string]string // map that contains tables with aliases
}

/*
Select - will set query type to SELECT and sets fields array.
*/
func (sql *SQLBuilder) Select(fields []string) Builder {
	sql.queryType = queryTypeSelect
	sql.parts.fields = append(sql.parts.fields, fields...)
	return sql
}

/*
Insert - will set query type to INSERT and sets table
*/
func (sql *SQLBuilder) Insert(table string) Builder {
	sql.queryType = queryTypeInsert
	sql.parts.table = table
	return sql
}

/*
Update — will set queryType to UPDATE and sets table
*/
func (sql *SQLBuilder) Update(table string) Builder {
	// setting table
	sql.queryType = queryTypeUpdate
	sql.parts.table = table
	sql.addToSources(table, tablePrefix)
	return sql
}

/*
Delete — will set queryType to DELETE and sets table
*/
func (sql *SQLBuilder) Delete() Builder {
	sql.queryType = queryTypeDelete
	return sql
}

/*
Set - alias for Values()
*/
func (sql *SQLBuilder) Set(data interface{}) Builder {
	return sql.Values(data)
}

/*
Values - map that will be users for Insert.
— key is for column
— value for column value
*/
func (sql *SQLBuilder) Values(data interface{}) Builder {
	sql.parts.insertData = data
	return sql
}

/*
From - will set table for query
*/
func (sql *SQLBuilder) From(table string) Builder {
	sql.parts.table = table
	sql.addToSources(table, tablePrefix)
	return sql
}

/*
ReturnID - return auto increment `id` after INSERT query
*/
func (sql *SQLBuilder) ReturnID(id string) Builder {
	sql.parts.returnID = id
	return sql
}

/*
Where - map that contains keys=values for SELECT/UPDATE/DELETE
*/
func (sql *SQLBuilder) Where(where map[string]interface{}) Builder {
	sql.parts.where = where
	return sql
}

/*
Join - join source with params into query.
Every table in SQL query have to have Alias. If you'll not provide - it will be generated
*/
func (sql *SQLBuilder) Join(jp JoinParam) Builder {
	if jp.SourceID == "" {
		jp.SourceID = tablePrefix + strconv.Itoa(len(sql.parts.join)+1)
	}
	sql.parts.join = append(sql.parts.join, jp)
	sql.addToSources(jp.Source, jp.SourceID)
	return sql
}

/*
Order - will set order by params for query
*/
func (sql *SQLBuilder) Order(o OrderParam) Builder {
	sql.parts.order = append(sql.parts.order, o)
	return sql
}

/*
Limit - limit and offset.
— offset by default is 0
- limit by default is defaultLimit
*/
func (sql *SQLBuilder) Limit(limit, offset int) Builder {
	sql.parts.limit = limit
	sql.parts.offset = offset
	return sql
}

/*
Build - method that builds from params into SQL string
*/
func (sql SQLBuilder) Build() string {
	if sql.queryType == queryTypeSelect {
		return sql.buildSelect()
	}
	if sql.queryType == queryTypeInsert {
		return sql.buildInsert()
	}
	if sql.queryType == queryTypeDelete {
		return sql.buildDelete()
	}
	if sql.queryType == queryTypeUpdate {
		return sql.buildUpdate()
	}
	return ""
}

func (sql *SQLBuilder) buildUpdate() (SQL string) {
	SQL = queryTypeUpdate
	SQL += sql.buildTable(true)
	SQL += sql.buildSetter()
	SQL += sql.buildWhere()
	return
}
func (sql *SQLBuilder) buildInsert() (SQL string) {
	SQL = queryTypeInsert
	SQL += " INTO " + sql.parts.table
	SQL += sql.buildValues()
	if sql.parts.returnID != "" {
		SQL += " RETURNING " + sql.parts.returnID
	}
	return
}
func (sql *SQLBuilder) buildDelete() (SQL string) {
	SQL = queryTypeDelete + " " + sql.getAliasBySource(sql.parts.table)
	SQL += sql.buildFrom(true)
	SQL += sql.buildWhere()
	return
}

func (sql *SQLBuilder) buildValues() string {
	var keys []string
	var values []string

	if data, ok := sql.parts.insertData.(map[string]interface{}); ok {
		for key, value := range data {
			keys = append(keys, "'"+key+"'")
			values = append(values, toString(value))
		}
	}
	return "(" + strings.Join(keys, ",") + ") VALUES (" + strings.Join(values, ",") + ")"
}

func (sql *SQLBuilder) buildSelect() (SQL string) {
	SQL = queryTypeSelect
	SQL += sql.buildFields()
	SQL += sql.buildFrom(true)
	SQL += sql.buildJoin()
	SQL += sql.buildWhere()
	SQL += sql.buildOrderBy()
	SQL += sql.buildLimit()
	return
}

func (sql *SQLBuilder) buildFrom(alias bool) string {
	return " FROM " + sql.buildTable(alias)
}
func (sql *SQLBuilder) buildTable(alias bool) (t string) {
	if alias == false {
		return " " + sql.parts.table
	}
	return " " + sql.parts.table + " as " + sql.getAliasBySource(sql.parts.table)
}
func (sql *SQLBuilder) buildFields() string {
	var fields []string
	if len(sql.parts.fields) == 0 {
		sql.parts.fields = []string{"*"}
	}
	for _, f := range sql.parts.fields {
		fields = append(fields, sql.getAliasBySource(sql.parts.table)+"."+f)
	}
	for _, j := range sql.parts.join {
		for _, f := range j.Fields {
			fields = append(fields, j.SourceID+"."+f)
		}
	}
	return " " + strings.Join(fields, ", ")
}

func (sql *SQLBuilder) buildJoin() (join string) {
	if len(sql.parts.join) == 0 {
		return
	}
	for _, j := range sql.parts.join {
		join += " " + strings.ToUpper(j.Type) + " JOIN " + j.Source + " AS " + j.SourceID + " ON "
		var keys []string
		for _, on := range j.On {
			var key string
			if on.Source != "" {
				key += sql.getAliasBySource(on.Source) + "." + on.SourceKey
			} else {
				key += sql.getAliasBySource(sql.parts.table) + "." + on.SourceKey
			}
			if on.JoinValue != nil {
				key += formatValue(on.JoinValue)
			} else {
				key += "=" + j.SourceID + "." + on.JoinKey
			}
			keys = append(keys, key)
		}
		join += strings.Join(keys, " AND ")
	}
	return
}

func formatValue(value interface{}) (fv string) {
	if v, ok := value.(string); ok {
		fv = "= '" + v + "'"
		return
	}
	if v, ok := value.([]int64); ok {
		var vs []string
		for _, n := range v {
			vs = append(vs, fmt.Sprintf("%v", n))
		}
		fv = " IN (" + strings.Join(vs, ",") + ")"
		return
	}
	if v, ok := value.([]float64); ok {
		var vs []string
		for _, n := range v {
			vs = append(vs, fmt.Sprintf("%v", n))
		}
		fv = " IN (" + strings.Join(vs, ",") + ")"
		return
	}
	if v, ok := value.([]string); ok {
		var vs []string
		for _, n := range v {
			vs = append(vs, "'"+n+"'")
		}
		fv = " IN (" + strings.Join(vs, ",") + ")"
		return
	}
	fv = "=" + fmt.Sprintf("%v", value)
	return
}

func (sql *SQLBuilder) buildWhere() (where string) {
	if len(sql.parts.where) == 0 {
		return
	}
	where = " WHERE "
	var w []string
	for key, value := range sql.parts.where {
		if sl, ok := value.([]int64); ok {
			var str []string
			for _, st := range sl {
				str = append(str, strconv.FormatInt(st, 10))
			}
			w = append(w, sql.getAliasBySource(sql.parts.table)+"."+key+" IN ("+strings.Join(str, ",")+")")
			continue
		}
		if sl, ok := value.([]string); ok {
			var str []string
			for _, st := range sl {
				str = append(str, `'`+st+`'`)
			}
			w = append(w, sql.getAliasBySource(sql.parts.table)+"."+key+" IN ("+strings.Join(str, ",")+")")
			continue
		}
		str := toString(value)
		sign := ""
		if strings.Index(key, "=") == -1 && strings.Index(key, ">") == -1 && strings.Index(key, "<") == -1 {
			sign = "="
		}
		w = append(w, sql.getAliasBySource(sql.parts.table)+"."+key+sign+str)
	}
	return where + strings.Join(w, " AND ")
}

func (sql *SQLBuilder) buildSetter() (where string) {
	if len(sql.parts.where) == 0 {
		return
	}
	where = " SET "
	var w []string
	if data, ok := sql.parts.insertData.(map[string]interface{}); ok {
		for key, value := range data {
			str := toString(value)
			w = append(w, key+" = "+str)
		}
	}
	return where + strings.Join(w, ", ")
}

func (sql *SQLBuilder) buildLimit() (limit string) {
	if sql.parts.limit != 0 {
		limit = " LIMIT "
		limit += strconv.Itoa(sql.parts.limit)
		limit += " OFFSET "
		limit += strconv.Itoa(sql.parts.offset)
	}
	return
}

func (sql *SQLBuilder) buildOrderBy() (order string) {
	if len(sql.parts.order) > 0 {
		var arr []string
		for _, o := range sql.parts.order {
			var item string
			if strings.Contains(o.OrderBy, ".") == false {
				item = sql.getAliasBySource(sql.parts.table) + "." + o.OrderBy
			} else {
				item = o.OrderBy
			}
			if o.Asc {
				item += " ASC"
			}
			if o.Desc {
				item += " DESC"
			}
			arr = append(arr, item)
		}
		order = " ORDER BY " + strings.Join(arr, ",")
	}
	return
}

func (sql *SQLBuilder) addToSources(table, id string) {
	if sql.sources == nil {
		sql.sources = make(map[string]string)
	}
	sql.sources[table] = id
}

func (sql *SQLBuilder) getAliasBySource(source string) string {
	if sql.sources[source] != "" {
		return sql.sources[source]
	}
	return source
}

func toString(value interface{}) (str string) {
	if v, ok := value.(float64); ok {
		str = strconv.FormatFloat(v, 'f', 8, 64)
	} else if v, ok := value.(int64); ok {
		str = strconv.FormatInt(v, 10)
	} else if v, ok := value.(int); ok {
		str = strconv.Itoa(v)
	} else {
		str = fmt.Sprint("'", value, "'")
	}
	return
}
