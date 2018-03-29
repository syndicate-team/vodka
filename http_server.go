package vodka

import (
	"log"
	"net/http"
	"strconv"
)

const (
	// ContentTypeJSON - content-type
	ContentTypeJSON = "application/json"

	// ErrorBadRequestCode - server HTTP code for BadRequest 400
	ErrorBadRequestCode = 400
	// ErrorUnathorizedCode - server HTTP code for Unathorized 401
	ErrorUnathorizedCode = 401
	// ErrorServerErrorCode - server HTTP code for ServerError 500
	ErrorServerErrorCode = 500
	// ErrorAccessDeniedCode - server HTTP code for ServerError 403
	ErrorAccessDeniedCode = 403
	// StatusOK - response with code 200
	StatusOK = 200
	// StatusNoContent - response with code 204
	StatusNoContent = 204
)

// ResponseNoContent - empty struct for empty response
type ResponseNoContent struct {
}

// HTTPConfig - HTTP server config
type HTTPConfig struct {
	Host        string
	Port        int
	ContentType string
}

/*
HTTPServer - duh!
*/
type HTTPServer struct {
	Config HTTPConfig
	Router *Router
}

/*
Start - stopping server (duuh!)
*/
func (srv *HTTPServer) Start() {
	log.Println("Starting server: ", "http://"+srv.getHost())
	log.Fatal(http.ListenAndServe(srv.getHost(), srv.Router.GetRouter()))
}

func (srv *HTTPServer) getHost() string {
	return srv.Config.Host + ":" + strconv.Itoa(srv.Config.Port)
}
