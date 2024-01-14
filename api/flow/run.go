package flow_apis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/lib/structs"
)

// Run flow step by step
func StartFlow(ch chan<- *structs.ExecutionStatusUpdate, app app.App, flow *structs.Flow) (*string, error) {
	execution, err := saveFlow(app, flow)
	if err == nil {
		err = app.GetMessageBroker().PushExecution(app.GetConfig(), execution)
		fmt.Println("Published execution in kafka with id: " + execution.Id)

		err = listenToExecWorker(ch, app, execution.Id)
		close(ch)
		if err != nil {
			return nil, err
		}
	}
	return &execution.Id, err
}

func listenToExecWorker(ch chan<- *structs.ExecutionStatusUpdate, app app.App, execId string) error {
	pubsubClient := app.GetPubSubClient()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	execTopic := "fr:" + execId
	respChan := pubsubClient.SubscribeTopic(ctx, execTopic)
	for resp := range respChan {
		if resp.Err != nil {
			return resp.Err
		} else {
			sr := structs.ExecutionStatusUpdate{}
			err := json.Unmarshal([]byte(*resp.Msg), &sr)
			if err != nil {
				fmt.Println("Invalid response from workers: " + err.Error())
				close(respChan)
				return err
			}
			ch <- &sr
		}
	}

	return nil
}
