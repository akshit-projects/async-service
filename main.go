package main

import (
	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/config"
	"github.com/akshitbansal-1/async-testing/be/server"
	thirdparty "github.com/akshitbansal-1/async-testing/be/third_party"
)

func main() {
	config := config.NewConfig()
	app := app.NewApp(config)

	if r := recover(); r != nil {
		thirdparty.Logger.Error("Error in the app", r)
	}
	server.NewServer(app)
}
