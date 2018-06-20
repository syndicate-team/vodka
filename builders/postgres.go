package builders

import (
	"fmt"
	"strconv"
	"strings"
)

/*
postgres - abstract builder for SQL-queries. Now adapted for Postgres
*/
type postgres struct {
	queryType string
	parts     parts
	sources   map[string]string // map that contains tables with aliases
}

/*
Select - will set query type to SELECT and sets fields array.
*/
func (sql *postgres) Select(fields []string) Builder {
	sql.queryType = queryTypeSelect
	sql.parts.fields = append(sql.parts.fields, fields...)
	return sql
}

/*
Insert - will set query type to INSERT and sets table
*/
func (sql *postgres) Insert(table string) Builder {
	sql.queryType = queryTypeInsert
	sql.parts.table = table
	return sql
}

/*
Save - will set query type to SAVE and sets table
*/
func (sql *postgres) Save(table string) Builder {
	sql.queryType = queryTypeSave
	sql.parts.table = table
	return sql
}

/*
Update — will set queryType to UPDATE and sets table
*/
func (sql *postgres) Update(table string) Builder {
	// setting table
	sql.queryType = queryTypeUpdate
	sql.parts.table = table
	sql.addToSources(table, tablePrefix)
	return sql
}

/*
Delete — will set queryType to DELETE and sets table
*/
func (sql *postgres) Delete() Builder {
	sql.queryType = queryTypeDelete
	return sql
}

/*
Set - alias for Values()
*/
func (sql *postgres) Set(data interface{}) Builder {
	return sql.Values(data)
}

/*
Values - map that will be users for Insert.
— key is for column
— value for column value
*/
func (sql *postgres) Values(data interface{}) Builder {
	sql.parts.insertData = data
	return sql
}

/*
From - will set table for query
*/
func (sql *postgres) From(table string) Builder {
	sql.parts.table = table
	sql.addToSources(table, tablePrefix)
	return sql
}

/*
ReturnID - return auto increment `id` after INSERT query
*/
func (sql *postgres) ReturnID(id string) Builder {
	sql.parts.returnID = id
	return sql
}

/*
Where - map that contains keys=values for SELECT/UPDATE/DELETE
*/
func (sql *postgres) Where(where map[string]interface{}) Builder {
	sql.parts.where = where
	return sql
}

/*
Join - join source with params into query.
Every table in SQL query have to have Alias. If you'll not provide - it will be generated
*/
func (sql *postgres) Join(jp Join) Builder {
	sql.parts.join = append(sql.parts.join, jp)
	length := len(sql.sources) + 1
	sql.addToSources(jp.Source, fmt.Sprintf("t%d", length))
	return sql
}

/*
Order - will set order by params for query
*/
func (sql *postgres) Order(o OrderParam) Builder {
	sql.parts.order = append(sql.parts.order, o)
	return sql
}

/*
Limit - limit and offset.
— offset by default is 0
- limit by default is defaultLimit
*/
func (sql *postgres) Limit(limit, offset int) Builder {
	sql.parts.limit = limit
	sql.parts.offset = offset
	return sql
}

/*
OnConflictAction - setting on conflict action.
Action must be 'UPDATE' of 'NOTHING'
*/
func (sql *postgres) OnConflictAction(action string) Builder {
	sql.parts.onConflictAction = action
	return sql
}

/*
OnConflictFields - setting unique fields to run on conflict script on.
*/
func (sql *postgres) OnConflictFields(fields []string) Builder {
	sql.parts.onConflictFields = fields
	return sql
}

/*
OnConflictConstraint - setting constraint name to run on conflict script.
*/
func (sql *postgres) OnConflictConstraint(constraint string) Builder {
	sql.parts.onConflictConstraint = constraint
	return sql
}

/*
Build - method that builds from params into SQL string
*/
func (sql postgres) Build() string {
	if sql.queryType == queryTypeSelect {
		return sql.buildSelect()
	}
	if sql.queryType == queryTypeInsert {
		return sql.buildInsert()
	}
	if sql.queryType == queryTypeSave {
		return sql.buildSave()
	}
	if sql.queryType == queryTypeDelete {
		return sql.buildDelete()
	}
	if sql.queryType == queryTypeUpdate {
		return sql.buildUpdate()
	}
	return ""
}

func (sql *postgres) buildUpdate() (SQL string) {
	SQL = queryTypeUpdate
	SQL += sql.buildTable(true)
	SQL += sql.buildSetter()
	SQL += sql.buildWhere()
	return
}
func (sql *postgres) buildInsert() (SQL string) {
	SQL = queryTypeInsert
	SQL += " INTO " + sql.parts.table
	SQL += sql.buildValues()
	if sql.parts.returnID != "" {
		SQL += " RETURNING " + sql.parts.returnID
	}
	return
}
func (sql *postgres) buildSave() (SQL string) {
	SQL = sql.buildInsert()
	SQL += sql.buildOnConflict()
	return
}
func (sql *postgres) buildDelete() (SQL string) {
	SQL = queryTypeDelete
	SQL += sql.buildFrom(true)
	SQL += sql.buildWhere()
	return
}
func (sql *postgres) buildValues() string {
	var keys []string
	var values []string

	if data, ok := sql.parts.insertData.(map[string]interface{}); ok {
		for key, value := range data {
			keys = append(keys, ""+key+"")
			values = append(values, toString(value))
		}
	}
	return "(" + strings.Join(keys, ",") + ") VALUES (" + strings.Join(values, ",") + ")"
}

func (sql *postgres) buildSelect() (SQL string) {
	SQL = queryTypeSelect
	SQL += sql.buildFields()
	SQL += sql.buildFrom(true)
	SQL += sql.buildJoin()
	SQL += sql.buildWhere()
	SQL += sql.buildOrderBy()
	SQL += sql.buildLimit()
	return
}

func (sql *postgres) buildFrom(alias bool) string {
	return " FROM " + sql.buildTable(alias)
}
func (sql *postgres) buildTable(alias bool) (t string) {
	if alias == false {
		return " " + sql.parts.table
	}
	return " " + sql.parts.table + " as " + sql.getAliasBySource(sql.parts.table)
}
func (sql *postgres) buildFields() string {
	var fields []string
	if len(sql.parts.fields) == 0 {
		sql.parts.fields = []string{"*"}
	}
	for _, f := range sql.parts.fields {
		fields = append(fields, sql.getAliasBySource(sql.parts.table)+"."+f)
	}
	for _, j := range sql.parts.join {
		for _, f := range j.Fields {
			fields = append(fields, sql.getAliasBySource(j.Source)+"."+f)
		}
	}
	return " " + strings.Join(fields, ", ")
}

func (sql *postgres) buildJoin() (join string) {
	if len(sql.parts.join) == 0 {
		return
	}
	for _, j := range sql.parts.join {
		src := sql.getAliasBySource(j.Source)
		join += " " + strings.ToUpper(j.Type) + " JOIN " + j.Source + " AS " + src + " ON "
		join += src + "." + j.Key + " = " + sql.getAliasBySource(sql.parts.table) + "." + j.TargetKey
	}
	return
}

func (sql *postgres) buildWhere() (where string) {
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

func (sql *postgres) buildSetter() (where string) {
	if len(sql.parts.where) == 0 {
		return
	}
	where = " SET "
	var w []string
	if data, ok := sql.parts.insertData.(map[string]interface{}); ok {
		for key, value := range data {
			str := toString(value)
			w = append(w, ""+key+" = "+str)
		}
	}
	return where + strings.Join(w, ", ")
}

func (sql *postgres) buildLimit() (limit string) {
	if sql.parts.limit != 0 {
		limit = " LIMIT "
		limit += strconv.Itoa(sql.parts.limit)
		limit += " OFFSET "
		limit += strconv.Itoa(sql.parts.offset)
	}
	return
}

func (sql *postgres) buildOrderBy() (order string) {
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

func (sql *postgres) buildOnConflict() (conflict string) {
	if sql.checkOnConflictParts() {
		if sql.parts.onConflictConstraint != "" {
			conflict = " ON CONFLICT ON CONSTRAINT " + sql.parts.onConflictConstraint
			conflict += sql.buildOnConflictAction()
			return
		}
		if len(sql.parts.onConflictFields) != 0 {
			conflict = " ON CONFLICT (" + strings.Join(sql.parts.onConflictFields, ",") + ")"
			conflict += sql.buildOnConflictAction()
			return
		}
	}

	return
}

func (sql *postgres) buildOnConflictAction() (conflict string) {
	action := sql.parts.onConflictAction

	if action != "" {
		if strings.ToUpper(action) == "NOTHING" {
			conflict = " DO NOTHING"
		}
		if strings.ToUpper(action) == "UPDATE" {
			conflict = " DO UPDATE SET "

			var w []string
			if data, ok := sql.parts.insertData.(map[string]interface{}); ok {
				for key, value := range data {
					str := toString(value)
					w = append(w, ""+key+" = "+str)
				}
			}
			conflict += strings.Join(w, ", ") + " RETURNING *"
		}
		return
	}
	return
}

func (sql *postgres) checkOnConflictParts() bool {
	action := strings.ToUpper(sql.parts.onConflictAction)
	if action == "UPDATE" || action == "NOTHING" {
		if sql.parts.onConflictConstraint != "" || len(sql.parts.onConflictFields) != 0 {
			return true
		}
		return false
	}
	return false
}

func (sql *postgres) addToSources(table, id string) {
	if sql.sources == nil {
		sql.sources = make(map[string]string)
	}
	sql.sources[table] = id
}

func (sql *postgres) getAliasBySource(source string) string {
	if sql.sources[source] != "" {
		return sql.sources[source]
	}
	return source
}
