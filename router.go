package vodka

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// HandlerFunc - function that handles requests
type HandlerFunc func(*Context) (interface{}, error)

// Handler — interface for handlers. Unused for now
type Handler interface {
	GetMethod(string) HandlerFunc
}

// Router — main struct for routing HTTP-requests
type Router struct {
	ID        string
	Routes    []string
	router    *httprouter.Router
	validator *Validator
	dispatch  func(*Context)
}

// NewRouter - router constructor
func NewRouter() *Router {
	return &Router{
		router: httprouter.New(),
	}
}

// SetValidator - setting validator for routes
func (r *Router) SetValidator(v *Validator) {
	r.validator = v
}

// GetRouter - getting router method
func (r *Router) GetRouter() *httprouter.Router {
	return r.router
}

func (r *Router) getValidationForPath(path, method string) methodRules {
	method = strings.ToLower(method)
	return r.validator.Rules[path][method]
}

// GET - HTTP-method GET setting handler
func (r *Router) GET(path string, h HandlerFunc) {
	r.router.GET(path, r.handle(h, r.getValidationForPath(path, "GET")))
}

// POST - HTTP-method POST setting handler
func (r *Router) POST(path string, h HandlerFunc) {
	r.router.POST(path, r.handle(h, r.getValidationForPath(path, "POST")))
}

// PUT - HTTP-method PUT setting handler
func (r *Router) PUT(path string, h HandlerFunc) {
	r.router.PUT(path, r.handle(h, r.getValidationForPath(path, "PUT")))
}

// DELETE - HTTP-method DELETE setting handler
func (r *Router) DELETE(path string, h HandlerFunc) {
	r.router.DELETE(path, r.handle(h, r.getValidationForPath(path, "DELETE")))
}

// PATCH - HTTP-method PATCH setting handler
func (r *Router) PATCH(path string, h HandlerFunc) {
	r.router.PATCH(path, r.handle(h, r.getValidationForPath(path, "PATCH")))
}

// OPTIONS - HTTP-method OPTIONS setting handler
func (r *Router) OPTIONS(path string, h HandlerFunc) {
	r.router.OPTIONS(path, r.handle(h, r.getValidationForPath(path, "OPTIONS")))
}

// HEAD - HTTP-method HEAD setting handler
func (r *Router) HEAD(path string, h HandlerFunc) {
	r.router.HEAD(path, r.handle(h, r.getValidationForPath(path, "HEAD")))
}

func (r *Router) handle(h HandlerFunc, v methodRules) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		body, err := parseBody(req)
		if err != nil {
			return
		}
		ctx := Context{
			Raw: RawContext{
				Query:  parseQuery(req.URL.Query()),
				Params: parseParams(ps),
				Body:   body,
			},
			Query:       KeyStorage{},
			Params:      KeyStorage{},
			Body:        KeyStorage{},
			Options:     KeyStorage{},
			HandlerFunc: h,
			Request:     req,
			Writer:      w,
			Validation:  v,
		}
		r.dispatch(&ctx)
	}
}

func parseBody(req *http.Request) ([]byte, error) {
	contentType := req.Header.Get("Content-Type")

	if strings.Index(contentType, "multipart/form-data") > -1 {
		req.ParseMultipartForm(1000000)
		d := make(map[string]interface{})
		for key, v := range req.Form {
			d[key] = v[0]
		}
		return json.Marshal(d)
	}

	if strings.Index(contentType, "x-www-form-urlencoded") > -1 {
		req.ParseForm()
		d := make(map[string]interface{})
		for key, v := range req.Form {
			d[key] = v[0]
		}
		return json.Marshal(d)
	}
	return ioutil.ReadAll(req.Body)
}

func parseParams(ps httprouter.Params) (params KeyStorage) {
	for _, param := range ps {
		params.Set(param.Key, param.Value)
	}
	return
}

func parseQuery(q url.Values) (query KeyStorage) {
	for key, v := range q {
		query.Set(key, v[0])
	}
	return
}
