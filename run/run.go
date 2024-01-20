package run

import (
	"context"
	"encoding/json"
	"errors"

	lib_contants "github.com/akshitbansal-1/async-testing/lib/constants"
	"github.com/akshitbansal-1/async-testing/lib/structs"
	ps "github.com/akshitbansal-1/async-testing/worker/pubsub"
	run_http "github.com/akshitbansal-1/async-testing/worker/run/api"
	run_kafka "github.com/akshitbansal-1/async-testing/worker/run/kafka"
	"github.com/akshitbansal-1/async-testing/worker/utils"
)

// Run flow step by step
func RunFlow(ch chan int, psc ps.PubSub, exec *structs.Execution) error {
	steps := exec.ExecutionFlow.Steps
	pubsubTopic := "fr:" + exec.Id
	for idx := range steps {
		step := &steps[idx]
		var stepResponse *structs.StepResponse
		switch step.Function {
		case lib_contants.HTTP_API_STEP:
			stepResponse = run_http.MakeAPICall(step)
		case lib_contants.PUBLISH_KAFKA_MESSAGE_STEP:
			stepResponse = run_kafka.SendMessage(step)
		case lib_contants.SUBSCRIBE_KAFKA_MESSAGES_STEP:
			stepResponse = run_kafka.ConsumeMessages(step)
		// case "purge-subscriptions":
		// 	stepResponse = purgeMessages(step)
		default:
			stepResponse = utils.CreateDefaultErrorResponse(step, errors.New("Unsupported function"))
		}

		msg, _ := json.Marshal(structs.ExecutionStatusUpdate{
			Type:    structs.EXECUTION_STATUS_SR,
			SR:      stepResponse,
			Message: "",
		})
		psc.PublishMessage(context.Background(), pubsubTopic, string(msg))
		if stepResponse.Status != structs.STEP_SUCCESS {
			publishCloseMessage(psc, pubsubTopic, ch)
			return nil
		}
	}
	publishCloseMessage(psc, pubsubTopic, ch)
	return nil
}

func publishCloseMessage(psc ps.PubSub, pubsubTopic string, ch chan int) {
	psc.PublishMessage(context.Background(), pubsubTopic, "close")
	<-ch
}
