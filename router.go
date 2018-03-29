package vodka

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type HandlerFunc func(*Context) (interface{}, error)

type Handler interface {
	GetMethod(string) HandlerFunc
}

type Router struct {
	ID       string
	Routes   []string
	router   *httprouter.Router
	dispatch func(*Context)
}

func NewRouter() *Router {
	return &Router{
		router: httprouter.New(),
	}
}

func (r *Router) GetRouter() *httprouter.Router {
	return r.router
}

func (r *Router) GET(path string, h HandlerFunc, v interface{}) {
	r.router.GET(path, r.handle(h, v))
}
func (r *Router) POST(path string, h HandlerFunc, v interface{}) {
	r.router.POST(path, r.handle(h, v))
}
func (r *Router) PUT(path string, h HandlerFunc, v interface{}) {
	r.router.PUT(path, r.handle(h, v))
}
func (r *Router) DELETE(path string, h HandlerFunc, v interface{}) {
	r.router.DELETE(path, r.handle(h, v))
}
func (r *Router) PATCH(path string, h HandlerFunc, v interface{}) {
	r.router.PATCH(path, r.handle(h, v))
}
func (r *Router) OPTIONS(path string, h HandlerFunc, v interface{}) {
	r.router.OPTIONS(path, r.handle(h, v))
}
func (r *Router) HEAD(path string, h HandlerFunc, v interface{}) {
	r.router.HEAD(path, r.handle(h, v))
}

func (r *Router) handle(h HandlerFunc, v interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		var body []byte
		var err error
		if body, err = ioutil.ReadAll(req.Body); err != nil {
			// w.Write().sendResponse(w, nil, NewServerError("reading_body", err))
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
