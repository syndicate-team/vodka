package vodka

import (
	"encoding/json"
	"log"
	"os"
)

// Application - main app struct
type Application struct {
	Router     *Router
	HTTPServer *HTTPServer

	middlewares []Middleware
	hooks       []Hook
	decorator   Decorator
	Debug       bool
}

// Middleware - middleware service
type Middleware func(*Context) (*Context, error)

// Hook - hook is almost the same as Middleware
type Hook func(*Context) (*Context, error)

// Decorator - decorates response
type Decorator func(interface{}, error) []byte

// Repository - interface that describes repository
type Repository interface {
	Get(string) interface{}
	Set(string, interface{})
}

/*
New - Application constructor
*/
func New() *Application {
	app := Application{
		Router: NewRouter(),
	}
	app.Router.dispatch = app.dispatch
	if os.Getenv("DEBUG") != "" {
		app.Debug = true
	}
	return &app
}

// Run - starting server to run
func (e *Application) Run() {
	log.Println("Running")
	e.HTTPServer.Start()
}

// Start - same as Run
func (e *Application) Start() {
	e.Run()
}

/*
Use - setting server Middleware
*/
func (e *Application) Use(m Middleware) {
	e.middlewares = append(e.middlewares, m)
}

/*
Hook - setting Hook.
*/
func (e *Application) Hook(hook Hook) {
	e.hooks = append(e.hooks, hook)
}

/*
Server - constructor of HTTP server!
*/
func (e *Application) Server(conf HTTPConfig) {
	if conf.ContentType == "" {
		conf.ContentType = ContentTypeJSON
	}
	e.HTTPServer = &HTTPServer{
		Config: conf,
		Router: e.Router,
	}
}

/*
Decorator - setting custom decorator for HTTP response
*/
func (e *Application) Decorator(d Decorator) {
	e.decorator = d
}

func (e *Application) dispatch(ctx *Context) {

	var err error
	// Validating request
	err = e.validate(ctx)
	if err != nil {
		e.sendResponse(ctx, nil, NewBadRequestError("validation", err.Error()))
		return
	}
	// Processing hooks
	if ctx, err = e.applyHooks(ctx); err != nil {
		// Response is sent
		return
	}

	// Processing middlewares
	if len(e.middlewares) > 0 {
		e.applyMiddlewares(ctx)
	} else {
		e.applyHandler(ctx)
	}
}

func (e *Application) applyMiddlewares(ctx *Context) {
	e.applyMiddleware(ctx)
}

func (e *Application) applyMiddleware(ctx *Context) {
	if ctx.iterator < len(e.middlewares) {
		handler := e.middlewares[ctx.iterator]
		ctx.iterator++
		if ctx.iterator < len(e.middlewares) {
			ctx.Next = e.applyMiddleware
		} else {
			ctx.Next = e.applyHandler
		}
		if _, err := handler(ctx); err != nil {
			e.sendResponse(ctx, nil, err)
		}
	} else {
		e.applyHandler(ctx)
	}
}

func (e *Application) applyHandler(ctx *Context) {
	// fmt.Println("Aplying handler")
	result, err := ctx.HandlerFunc(ctx)
	e.sendResponse(ctx, result, err)
}

func (e *Application) applyHooks(ctx *Context) (*Context, error) {
	var err error
	if len(e.hooks) > 0 {
		for _, hook := range e.hooks {
			if ctx, err = hook(ctx); err != nil {
				e.sendResponse(ctx, nil, err)
				return ctx, err
			}
		}
	}
	return ctx, err
}

func (e *Application) decorate(data interface{}, err error) []byte {
	if e.decorator != nil {
		return e.decorator(data, err)
	}
	if data == nil {
		data = make(map[string]string)
	}
	response := make(map[string]interface{})
	response["data"] = data
	if e, ok := err.(Error); ok {
		response["error"] = e
	} else if err != nil {
		response["error"] = err.Error()
	} else {
		response["error"] = nil
	}
	b, err := json.Marshal(response)
	if err != nil {
		return []byte(err.Error())
	}
	return b
}

func (e *Application) sendResponse(ctx *Context, data interface{}, err error) {
	ctx.Writer.Header().Set("Content-Type", e.HTTPServer.Config.ContentType)
	if e, ok := err.(Error); ok {
		ctx.Writer.WriteHeader(e.httpCode)
	} else {
		if err != nil {
			ctx.Writer.WriteHeader(ErrorServerErrorCode)
		}
	}
	if _, ok := data.(ResponseNoContent); ok {
		ctx.Writer.WriteHeader(StatusNoContent)
		ctx.Writer.Write([]byte(""))
	} else {
		ctx.Writer.Write(e.decorate(data, err))
	}
}
