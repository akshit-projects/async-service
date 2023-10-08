package scheduler

import (
	"encoding/json"
	"fmt"

	"github.com/akshitbansal-1/async-testing/lib/structs"
	"github.com/akshitbansal-1/async-testing/worker/config"
)

var schedules chan int

type Scheduler interface {
	ProcessMessage(km []byte) error
}

type scheduler struct {
	maxExecutions int
}

func (s *scheduler) ProcessMessage(km []byte) error {
	var exec structs.Execution
	err := json.Unmarshal(km, &exec)
	if err != nil {
		return err
	}

	schedules <- 1
	if r := recover(); r != nil {
		fmt.Println("Recovered")
	}

	RunFlow(schedules, &exec)

	return nil
}

func NewScheduler(config *config.Configuration) Scheduler {
	schedules = make(chan int, config.MaxExecutions)
	return &scheduler{
		config.MaxExecutions,
	}
}
