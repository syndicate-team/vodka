package vodka

import (
	"log"
	"net/http"
	"strconv"
)

const (
	ContentTypeJSON = "application/json"

	RequestError      = 400
	UnathorizedError  = 401
	AccessDeniedError = 403
	ServerError       = 500
	StatusOK          = 200
	StatusNoContent   = 204
)

type ResponseNoContent struct {
}

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
