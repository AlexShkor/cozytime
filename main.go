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
	games := data.NewGamesService(db.C("games"))
	e := echo.New()

	e.Hook(stripTrailingSlash)

	e.Use(mw.Logger())
	e.Use(mw.Recover())

	router := routes.NewRouter(tokens)
	e.Get("/", router.HelloWorld)
	e.Post("/authorize", router.Register)
    e.Post("/submitcode", router.SubmitCode)
	
	authorizedGroup := e.Group("/api", BearerAuth(tokens))
	authorizedGroup.Post("/setname", router.SetName)
	
	game := routes.NewGameRouter(games)
	
	authorizedGroup.Post("/game/start", game.StartGame)
	authorizedGroup.Post("/game/create", game.CreateGame)
	authorizedGroup.Post("/game/join", game.JoinGame)
	authorizedGroup.Post("/game/stop", game.StopGame)
	
	authorizedGroup.Post("/secret", func(c *echo.Context) error {
		return c.String(200, "You are authorized!\n")
	})
	
	adminGroup := e.Group("/admin")

	adminGroup.Static("/assets", "assets")
	adminGroup.Get("", router.AdminIndex)

	e.Run(":" + conf.Port)
}
