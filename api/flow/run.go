package flow_apis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/akshitbansal-1/async-testing/be/app"
	thirdparty "github.com/akshitbansal-1/async-testing/be/third_party"
	"github.com/akshitbansal-1/async-testing/lib/structs"
	"github.com/akshitbansal-1/async-testing/lib/utils"
	lib_utils "github.com/akshitbansal-1/async-testing/lib/utils"
)

var logger = thirdparty.Logger

// Run flow step by step
func StartFlow(ch chan<- *structs.ExecutionStatusUpdate, app app.App, flow *structs.Flow) (*string, error) {
	execution, err := saveFlow(app, flow)
	if err != nil {
		logger.Error("Unable to save execution in mongodb ", err.Error())
		ch <- lib_utils.CreateErrorExecutionStatus("Internal server error. An unknown error occurred", structs.ES_MONGO_ERROR)
		return nil, err
	}

	err = app.GetMessageBroker().PushExecution(app.GetConfig(), execution)
	if err != nil {
		logger.Error("Unable to publish execution to kafka ", err.Error())
		ch <- lib_utils.CreateErrorExecutionStatus("Internal server error. An unknown error occurred", structs.ES_KAFKA_ERROR)
		return &execution.Id, err
	}

	err = listenToExecWorker(ch, app, execution)
	if err != nil {
		logger.Error("Unable to listen to pubsub for error messages", err.Error())
		ch <- lib_utils.CreateErrorExecutionStatus("Unable to listen to realtime updates. The execution is running in background", structs.ES_RT_UPDATES_ERROR)
		return nil, err
	}

	close(ch)

	return &execution.Id, nil
}

func listenToExecWorker(ch chan<- *structs.ExecutionStatusUpdate, app app.App, execution *structs.Execution) error {
	timeout := getTotalStepsTimeout(execution.ExecutionFlow.Steps)
	pubsubClient := app.GetPubSubClient()
	execTopic := "fr:" + execution.Id
	ctx := context.Background()
	respChan := pubsubClient.SubscribeTopic(ctx, execTopic)

	var er error
	isTimedOut := utils.Race(ctx, func() {
		go func() {
			<-time.After(time.Duration(timeout) * time.Second)
			close(respChan)
		}()

		for resp := range respChan {
			if resp.Err != nil {
				er = resp.Err
				return
			}

			if *resp.Msg == "close" {
				er = nil
				return
			}

			sr := structs.ExecutionStatusUpdate{}
			err := json.Unmarshal([]byte(*resp.Msg), &sr)
			if err != nil {
				logger.Error("Invalid response from execution workers: ", err.Error())
				er = err
				return
			}
			ch <- &sr
		}
	}, timeout)

	if isTimedOut {
		return errors.New("Unable to listen to realtime updates from redis pubsub")
	}

	return er
}

func getTotalStepsTimeout(steps []structs.Step) int {
	timeout := 0
	for _, step := range steps {
		timeout += step.Timeout
	}
	return timeout + 100 // add 100ms as buffer for communication
}
