package scheduler

import (
	"encoding/json"
	"fmt"

	"github.com/akshitbansal-1/async-testing/lib/structs"
	"github.com/akshitbansal-1/async-testing/lib/utils"
	"github.com/akshitbansal-1/async-testing/worker/config"
	"github.com/akshitbansal-1/async-testing/worker/pubsub"
	"github.com/akshitbansal-1/async-testing/worker/run"
)

var schedules chan int

type Scheduler interface {
	ProcessMessage(km []byte) error
}

type scheduler struct {
	maxExecutions int
	pubsubClient  pubsub.PubSub
}

func (s *scheduler) ProcessMessage(km []byte) error {
	var exec structs.Execution
	fmt.Println(string(km))
	err := json.Unmarshal(km, &exec)
	if err != nil {
		fmt.Println("error while debugging")
		return err
	}

	schedules <- 1
	if r := recover(); r != nil {
		fmt.Println("Recovered")
	}

	fmt.Println(utils.StructToString(exec))

	run.RunFlow(schedules, s.pubsubClient, &exec)

	return nil
}

func NewScheduler(config *config.Configuration) Scheduler {
	schedules = make(chan int, config.MaxExecutions)
	psc := pubsub.NewPubSubClient(config)
	return &scheduler{
		config.MaxExecutions,
		psc,
	}
}
