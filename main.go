package main

import (
	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/config"
	"github.com/akshitbansal-1/async-testing/be/server"
)

func main() {
	config := config.NewConfig()
	app := app.NewApp(config)

	server.NewServer(app)
}
