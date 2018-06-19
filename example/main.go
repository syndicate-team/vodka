package main

import (
	"log"

	"github.com/syndicatedb/vodka/example/modules/items"

	"github.com/syndicatedb/vodka"
	"github.com/syndicatedb/vodka/adapters"
	"github.com/syndicatedb/vodka/example/controllers"
	"github.com/syndicatedb/vodka/example/modules/orders"
	"github.com/syndicatedb/vodka/example/modules/users"
)

type infrastructure struct {
	Postgres adapters.Adapter
	MySQL    adapters.Adapter
}

var repos infrastructure
var config Config
var err error

var userModule *users.API
var orderModule *orders.API
var itemsModule *items.API

func init() {
	if config, err = NewConfig("./config.json"); err != nil {
		log.Fatalln("Error reading config: ", err)
	}
	repos.Postgres = adapters.NewPostgres(config.Postgres)
	repos.MySQL = adapters.NewMySQL(config.MySQL)

	userModule = users.New(repos.Postgres)
	orderModule = orders.New(repos.MySQL)
	itemsModule = items.New(repos.Postgres)
}

func main() {
	engine := vodka.New()
	engine.Server(config.HTTPServer)
	engine.Validation("./validation.json")

	userCtrl := controllers.NewUsers(userModule)
	orderCtrl := controllers.NewOrders(orderModule)
	itemsCtrl := controllers.NewItems(itemsModule)

	engine.Router.GET("/users", userCtrl.Find)
	engine.Router.GET("/users/:id", userCtrl.FindByID)
	engine.Router.POST("/users", userCtrl.Create)
	engine.Router.PUT("/users/:id", userCtrl.UpdateByID)
	engine.Router.PUT("/users", userCtrl.Update)
	engine.Router.DELETE("/users/:id", userCtrl.DeleteByID)

	engine.Router.GET("/orders", orderCtrl.Find)

	engine.Router.POST("/items", itemsCtrl.Save)

	for {
		engine.Start()
	}
}
