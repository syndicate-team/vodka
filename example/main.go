package main

import (
	"log"

	"github.com/niklucky/vodka"
	"github.com/niklucky/vodka/adapters"
	"github.com/niklucky/vodka/example/controllers"
	"github.com/niklucky/vodka/example/modules/users"
)

type infrastructure struct {
	Postgres adapters.Adapter
}

var repos infrastructure
var config Config
var err error

var userModule *users.API

func init() {
	if config, err = NewConfig("./config.json"); err != nil {
		log.Fatalln("Error reading config: ", err)
	}
	repos.Postgres = adapters.NewPostgres(config.Postgres)

	userModule = users.New(repos.Postgres)
}

func main() {
	engine := vodka.New()
	engine.Server(config.HTTPServer)

	userCtrl := controllers.NewUsers(userModule)
	var userValidation controllers.UserValidation

	engine.Router.GET("/users", userCtrl.Find, userValidation.Find)
	engine.Router.GET("/users/:id", userCtrl.FindByID, userValidation.FindByID)
	engine.Router.POST("/users", userCtrl.Create, userValidation.Create)

	for {
		engine.Start()
	}
}
