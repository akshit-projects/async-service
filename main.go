package main

import (
	"github.com/akshitbansal-1/async-testing/worker/config"
	message_broker "github.com/akshitbansal-1/async-testing/worker/message-broker"
	"github.com/akshitbansal-1/async-testing/worker/scheduler"
)

func main() {
	config := config.NewConfig()
	scheduler := scheduler.NewScheduler(config)
	message_broker.InitConsumer(config, scheduler)
}
