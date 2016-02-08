package main

import (
	"fmt"
	"os"

	"bitbucket.org/AlexShkor/cozytime/data"
	"bitbucket.org/AlexShkor/cozytime/routes"
	"bitbucket.org/AlexShkor/cozytime/settings"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"

	"labix.org/v2/mgo"
)

var db *mgo.Database

func main() {

	conf := settings.Get()

	if conf == nil {
		fmt.Println("No config found, terminating.")
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
	e.Get("/sendsms", router.SendSms)
	e.Post("/authorize", router.Register)
    e.Post("/submitcode", router.SubmitCode)
	
	authorizedGroup := e.Group("/api", BearerAuth(tokens))
	authorizedGroup.Post("/setname", router.SetName)
	adminGroup := e.Group("/admin")

	adminGroup.Static("/assets", "assets")
	adminGroup.Get("", router.AdminIndex)

	e.Run(":" + conf.Port)
}
