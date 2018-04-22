package repositories

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/niklucky/vodka/builders"
)

func populateStructByMap(rv reflect.Value, data map[string]interface{}) interface{} {
	// val := reflect.ValueOf(str)
	st := rv.Elem()
	t := st.Type()
	// fmt.Printf("val: %+v\n", st)
	// fmt.Printf("data: %+v\n", data)
	// fmt.Printf("type: %+v\n", t)
	for i := 0; i < st.NumField(); i++ {
		key := t.Field(i).Name
		if t.Field(i).Tag.Get("db") != "" {
			key = t.Field(i).Tag.Get("db")
		}
		// fmt.Println("key: ", key)
		if v, ok := data[key]; ok {
			if v == nil {
				continue
			}
			// fmt.Printf("v: %T, %s\n", v, st.Field(i).Type().String())
			switch st.Field(i).Type().String() {
			case "int64":
				st.Field(i).SetInt(getInt64(v))
			case "float64":
				st.Field(i).SetFloat(getFloat64(v))
			case "string":
				st.Field(i).SetString(fmt.Sprintf("%v", v))
			case "bool":
				st.Field(i).SetBool(getBool(v))
			case "time.Time":
				st.Field(i).Set(reflect.ValueOf(getTime(v)))
			}
		}
	}
	return st.Interface()
}

func getTime(v interface{}) time.Time {
	switch v.(type) {
	case int64:
		str := strconv.FormatInt(v.(int64), 10)
		t, _ := time.Parse("U", str)
		return t
	case string:
		t, _ := time.Parse(time.RFC3339, v.(string))
		return t
	case time.Time:
		return v.(time.Time)
	}
	return time.Unix(0, 0)
}

func getInt64(v interface{}) int64 {
	switch v.(type) {
	case int64:
		return v.(int64)
	case int:
		return int64(v.(int))
	case int32:
		return int64(v.(int32))
	case float32:
		return int64(v.(float32))
	case float64:
		return int64(v.(float64))
	case string:
		a, _ := strconv.Atoi(v.(string))
		return int64(a)
	}
	return 0
}

func getFloat64(v interface{}) float64 {
	switch v.(type) {
	case int64:
		return float64(v.(int64))
	case int:
		return float64(v.(int))
	case int32:
		return float64(v.(int32))
	case float32:
		return float64(v.(float32))
	case float64:
		return v.(float64)
	case string:
		a, _ := strconv.ParseFloat(v.(string), 64)
		return a
	}
	return 0
}
func getBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

func parseParams(params interface{}) (m QueryModificator) {
	if params == nil {
		return
	}
	// if p, ok := params.(map[string]interface{}); ok {
	if p, ok := params.(ParamsMap); ok {
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
			var orderParams builders.OrderParam
			var orderParamsArr []builders.OrderParam

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
