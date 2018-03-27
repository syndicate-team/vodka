package vodka

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func (e *Application) validate(ctx *Context) (err error) {
	var errs []string
	if e.Debug {
		log.Printf("Validation: %+v", ctx.Validation)
	}
	v := ctx.Validation

	rt := reflect.TypeOf(v)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Name == "Query" {
			ctx.Query, err = validateMap(field.Type, ctx.Raw.Query)
			if err != nil {
				errs = append(errs, "Query: "+err.Error())
			}
		}
		if field.Name == "Params" {
			ctx.Params, err = validateMap(field.Type, ctx.Raw.Params)
			if err != nil {
				errs = append(errs, "Params: "+err.Error())
			}
		}
	}

	ctx.Options.Set("params", getParamsFromQuery(ctx.Raw.Query))

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}

func getParamsFromQuery(q KeyStorage) interface{} {
	p := make(map[string]interface{})
	if q.Get("__limit") != nil {
		if limit, err := strconv.Atoi(q.GetString("__limit")); err == nil {
			p["limit"] = limit
		}
	}
	if q.Get("__skip") != nil {
		if skip, err := strconv.Atoi(q.GetString("__skip")); err == nil {
			p["skip"] = skip
		}
	}
	if q.Get("__orderBy") != nil {
		p["orderBy"] = q.GetString("__orderBy")
	}
	if q.Get("__order") != nil {
		p["order"] = q.GetString("__order")
	}
	return p
}

func validateMap(rt reflect.Type, dv KeyStorage) (ks KeyStorage, err error) {
	var errs []string
	var typeErr error
	nf := rt.NumField()
	for n := 0; n < nf; n++ {
		var v interface{}
		field := rt.Field(n)
		// fmt.Printf("--- Fields:  %+v — %+v\n", field.Name, field.Tag)
		var fieldName string
		if field.Tag.Get("input") != "" {
			fieldName = field.Tag.Get("input")
		} else {
			fieldName = field.Name
		}
		value := dv.Get(fieldName)
		// fmt.Printf("Value: %s - %+v\n", fieldName, value)
		required := field.Tag.Get("required")
		if dv.Get(fieldName) == nil {
			if required == "true" {
				errs = append(errs, field.Name+" is not defined")
			}
			continue
		}

		v, typeErr = validateType(field.Name, value, field.Type.String())
		if typeErr != nil {
			errs = append(errs, typeErr.Error())
			continue
		}
		ks.Set(field.Name, v)
	}
	if len(errs) > 0 {
		return ks, errors.New(strings.Join(errs, ", "))
	}
	return
}

func validateBody(rule interface{}, b []byte, p KeyStorage) (KeyStorage, error) {
	if rule == nil {
		return KeyStorage{}, nil
	}
	var err error
	if string(b) == "" {
		return p, nil
	}
	if rule == nil {
		return p, nil
	}
	var data interface{}
	if err = json.Unmarshal(b, &data); err != nil {
		return p, err
	}
	for key, v := range data.(map[string]interface{}) {
		p.Set(key, v)
	}
	// err = validateMap(rule, p)
	return p, err
}

func validateType(key string, value interface{}, t string) (res interface{}, err error) {
	if value == nil {
		return value, nil
	}
	switch v := value.(type) {
	case int:
		if t == "int" {
			res = v
			return
		}
		if t == "float" {
			res = float64(v)
			return
		}
		if t == "bool" {
			if v > 0 {
				res = true
			} else {
				res = false
			}
			return
		}
		if t == "string" {
			res = strconv.Itoa(v)
			return
		}

	case string:
		if t == "int" {
			res, err = strconv.Atoi(v)
			return
		}
		if t == "int64" {
			res, err = strconv.ParseInt(v, 10, 64)
			return
		}
		if t == "float" {
			res, err = strconv.ParseFloat(v, 64)
			return
		}
		if t == "bool" {
			if v == "true" || v == "1" || v == "TRUE" {
				res = true
				return
			}
			if v == "false" || v == "0" || v == "FALSE" {
				res = false
				return
			}
		}
		if t == "string" {
			res = v
			return
		}
		if t == "[]int64" {
			var si []int64
			sl := strings.Split(v, ",")
			for _, str := range sl {
				iv, ferr := strconv.ParseInt(str, 10, 64)
				if ferr != nil {
					err = fmt.Errorf("%s (%v): slice element %s is not %s", key, value, str, t)
				}
				si = append(si, iv)
			}
			res = si
			return
		}
		if t == "[]float64" {
			var si []float64
			sl := strings.Split(v, ",")
			for _, str := range sl {
				iv, ferr := strconv.ParseFloat(str, 64)
				if ferr != nil {
					err = fmt.Errorf("%s (%v): slice element %s is not %s", key, value, str, t)
				}
				si = append(si, iv)
			}
			res = si
			return
		}
		if t == "[]string" {
			res = strings.Split(v, ",")
			return
		}
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}

	return nil, formatError(key, value, t)
}

func formatError(key string, value interface{}, t string) error {
	return fmt.Errorf("%s (%v) type is not valid (expected %s)", key, value, t)
}
