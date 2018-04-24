package vodka

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	lib "github.com/niklucky/go-lib"
)

type Validator struct {
	Rules map[string]routeRules
}

func (v *Validator) loadRules(fileName string) error {
	fileData, err := lib.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileData, &v.Rules)
	if err != nil {
		return err
	}
	return nil
}

type routeRules map[string]methodRules

type methodRules struct {
	Params  map[string]validation  `json:"params"`
	Query   map[string]validation  `json:"query"`
	Body    map[string]validation  `json:"body"`
	Options map[string]interface{} `json:"options"`
}

type validation struct {
	InputType string `json:"type"`
	Required  bool   `json:"required"`
	Name      string `json:"name"`
}

func (e *Application) validate(ctx *Context) (err error) {
	var errs []string
	if isDebug {
		log.Printf("Validation rules: %+v", ctx.Validation)
	}
	v := ctx.Validation
	if v.Params != nil {
		ctx.Params, err = validateMap(v.Params, ctx.Raw.Params)
		if err != nil {
			errs = append(errs, "Params: "+err.Error())
		}
	}
	if v.Query != nil {
		ctx.Query, err = validateMap(v.Query, ctx.Raw.Query)
		if err != nil {
			errs = append(errs, "Query: "+err.Error())
		}
	}
	if v.Body != nil {
		var body map[string]interface{}
		json.Unmarshal(ctx.Raw.Body, &body)
		var b KeyStorage
		for key, v := range body {
			b.Set(key, v)
		}
		ctx.Body, err = validateMap(v.Body, b)
		if err != nil {
			errs = append(errs, "Params: "+err.Error())
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

func validateMap(vm map[string]validation, dv KeyStorage) (ks KeyStorage, err error) {
	var errs []string
	var typeErr error
	for key, val := range vm {
		var v interface{}
		var name = key
		if val.Name != "" {
			name = val.Name
		}
		value := dv.Get(name)
		fmt.Printf("Value: %s - %+v\n", name, value)
		if value == nil {
			if val.Required {
				errs = append(errs, name+" is not defined")
			}
			continue
		}

		v, typeErr = validateType(name, value, val.InputType)
		if typeErr != nil {
			errs = append(errs, typeErr.Error())
			continue
		}
		ks.Set(key, v)
	}
	if len(errs) > 0 {
		return ks, errors.New(strings.Join(errs, ", "))
	}
	return
}

func validateType(key string, value interface{}, t string) (res interface{}, err error) {
	if value == nil {
		return value, nil
	}
	fmt.Printf("%T\n", value)
	switch v := value.(type) {
	case int64:
		if t == "int64" {
			res = int64(v)
			return
		}
		if t == "float64" {
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
			res = strconv.FormatInt(v, 10)
			return
		}
	case float64:
		if t == "int64" {
			res = int64(v)
			return
		}
		if t == "float64" {
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
			res = strconv.FormatFloat(v, 'f', 10, 64)
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
		if t == "float" || t == "float64" {
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
