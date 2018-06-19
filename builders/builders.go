package builders

import (
	"fmt"
	"strconv"
	"strings"
)

// NewPostgres - Postgres SQL builder
func NewPostgres() Builder {
	return &postgres{}
}

// NewMySQL - MySQL SQL builder
func NewMySQL() Builder {
	return &mysql{}
}

/*
Builder - interface for query builder for adapter and data service
*/
type Builder interface {
	Select([]string) Builder
	Insert(string) Builder
	Save(string) Builder
	Update(string) Builder
	Delete() Builder
	ReturnID(string) Builder
	Values(interface{}) Builder
	Set(interface{}) Builder
	From(string) Builder
	Where(map[string]interface{}) Builder
	Limit(int, int) Builder
	Join(Join) Builder
	Order(OrderParam) Builder
	OnConflictAction(string) Builder
	OnConflictFields([]string) Builder
	OnConflictConstraint(string) Builder
	Build() string
}

/*
Join - joining Repository (table to query)
*/
type Join struct {
	Source    string
	Key       string
	TargetKey string
	Fields    []string
	Type      string
}

// OrderParam - ordering params
type OrderParam struct {
	OrderBy string
	Asc     bool
	Desc    bool
}

type parts struct {
	table                string
	fields               []string
	where                map[string]interface{}
	join                 []Join
	order                []OrderParam
	limit                int
	offset               int
	insertData           interface{}
	returnID             string
	onConflictAction     string
	onConflictFields     []string
	onConflictConstraint string
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
