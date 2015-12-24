package main

import (
	"fmt"
	"os"

	"github.com/temp/prototype/data"
	"github.com/temp/prototype/routes"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"

	"labix.org/v2/mgo"
)

var db *mgo.Database

func main() {

	configPath := "config.json"
	conf, err := ReadConfig(configPath)
	if err != nil {
		fmt.Printf("Failed to read config: %s\n", err)
		os.Exit(-1)
	}

	mongoSession, err := mgo.Dial(conf.ConnectionString)
	if err != nil {
		panic(err)
	}
	db = mongoSession.DB(conf.DatabaseName)
	tokens := data.NewTokensService(db.C(conf.CollectionName))

	e := echo.New()

	e.Hook(stripTrailingSlash)

	e.Use(mw.Logger())
	e.Use(mw.Recover())

	router := routes.NewRouter(tokens)
	e.Get("/", router.HelloWorld)
	e.Post("/authorize", router.Authorize)

	authorizedGroup := e.Group("/api", BearerAuth(tokens))
	authorizedGroup.Post("/secret", func(c *echo.Context) error {
		return c.String(200, "You are authorized!\n")
	})

	adminGroup := e.Group("/admin")
	adminGroup.Static("/assets", "assets")
	adminGroup.Get("", router.AdminIndex)

	e.Run(":1323")
}
