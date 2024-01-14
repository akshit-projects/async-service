package run

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"os"

	"cloud.google.com/go/pubsub"
	"github.com/akshitbansal-1/async-testing/lib/structs"
	"github.com/akshitbansal-1/async-testing/lib/structs/api"
	gcp_pubsub "github.com/akshitbansal-1/async-testing/lib/structs/gcp/pubsub"
	"github.com/akshitbansal-1/async-testing/lib/utils"
	"github.com/akshitbansal-1/async-testing/worker/constants"
	ps "github.com/akshitbansal-1/async-testing/worker/pubsub"
)

// Run flow step by step
func RunFlow(ch chan int, psc ps.PubSub, exec *structs.Execution) error {
	steps := exec.ExecutionFlow.Steps
	pubsubTopic := "fr:" + exec.Id
	for idx := range steps {
		step := &steps[idx]
		var stepResponse *structs.StepResponse
		switch step.Function {
		case constants.API_STEP:
			stepResponse = makeAPICall(step)
		case constants.PUBLISH_MESSAGE_STEP:
			stepResponse = publishMessages(step)
		case constants.SUBSCRIBE_MESSAGES_STEP:
			stepResponse = subscribeMessages(step)
		case "purge-subscriptions":
			stepResponse = purgeMessages(step)
		default:
			stepResponse = createDefaultErrorResponse(step, errors.New("Unsupported function"))
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

// Make HTTP call
func makeAPICall(step *structs.Step) *structs.StepResponse {
	h := api.HTTPRequest{}
	_ = utils.ParseInterface(step.Meta, &h)
	req, _ := http.NewRequest(h.Method, h.Url, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}

	var isGetCall = h.Method == "GET"
	if !isGetCall {
		bodyBytes, err := json.Marshal(h.Body)
		if err != nil {
			fmt.Println("Error marshaling body:", err)
			return createDefaultErrorResponse(step, err)
		}
		// set body
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	// Perform the HTTP request
	resp, err := utils.CallHTTP(req)
	if err != nil {
		fmt.Println("HTTP request error:", err)
		if errors.Is(err, context.DeadlineExceeded) {
			err = errors.New("Request timed out")
		}
		return createDefaultErrorResponse(step, err)
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	stepValue := &api.HTTPResponse{
		Status:   resp.StatusCode,
		Response: buf.String(),
	}
	fmt.Println("Made API call")

	if h.ExpectedStatus != "" && h.ExpectedStatus !=
		strconv.Itoa(resp.StatusCode) {
		return &structs.StepResponse{
			Name:   step.Name,
			Status: structs.STEP_ERROR,
			Response: &structs.StepError{
				Expected: h.ExpectedStatus,
				Actual:   strconv.Itoa(resp.StatusCode),
				Error:    "Status code not matching",
			},
			Id: step.Id,
		}
	}

	if err = utils.CompareStrings(&stepValue.Response, &h.ExpectedResponse); err != nil {
		return &structs.StepResponse{
			Name:   step.Name,
			Status: structs.STEP_ERROR,
			Response: &structs.StepError{
				Expected: h.ExpectedResponse,
				Actual:   stepValue.Response,
				Error:    "Response not matching",
			},
			Id: step.Id,
		}
	}

	return &structs.StepResponse{
		Name:     step.Name,
		Status:   structs.STEP_SUCCESS,
		Response: stepValue,
		Id:       step.Id,
	}
}

// Publish messages in pubsub
func publishMessages(step *structs.Step) *structs.StepResponse {
	publishReq := gcp_pubsub.PublishRequest{}
	_ = utils.ParseInterface(step.Meta, &publishReq)
	projectID := publishReq.ProjectId
	topicID := publishReq.TopicName

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fmt.Println("Unable to connect to pubsub", err)
		return createDefaultErrorResponse(step, err)
	}

	topic := client.Topic(topicID)

	messageIds := []string{}
	for _, message := range publishReq.Messages {
		result := topic.Publish(ctx, &pubsub.Message{
			Data: []byte(message),
		})
		id, err := result.Get(ctx)
		if err != nil {
			fmt.Fprintln(os.Stdout, []any{"Failed to publish: %v", err}...)
			return createDefaultErrorResponse(step, err)
		}
		messageIds = append(messageIds, id)
	}

	return &structs.StepResponse{
		Name:   step.Name,
		Status: structs.STEP_SUCCESS,
		Response: &gcp_pubsub.PublishResponse{
			MessageIds: messageIds,
		},
		Id: step.Id,
	}
}

// Subscribe to messages in pubsub
func subscribeMessages(step *structs.Step) *structs.StepResponse {
	subscribeReq := gcp_pubsub.SubscribeRequest{}
	_ = utils.ParseInterface(step.Meta, &subscribeReq)
	projectID := subscribeReq.ProjectId
	subscriptionID := subscribeReq.SubscriptionName

	msgs, err := fetchPubsubMessages(projectID, subscriptionID)
	if err != nil {
		return createDefaultErrorResponse(step, err)
	}

	return &structs.StepResponse{
		Name:   step.Name,
		Status: structs.STEP_SUCCESS,
		Response: &gcp_pubsub.SubscribeResponse{
			Messagess: msgs,
		},
		Id: step.Id,
	}
}

// Fetch messages of given subscription name
func fetchPubsubMessages(projectID string, subscriptionID string) ([]string, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fmt.Println("Unable to connect to pubsub", err)
		return nil, err
	}

	subscription := client.Subscription(subscriptionID)
	subscriptionTimeout := 10 * time.Second
	subCtx, cancel := context.WithTimeout(ctx, subscriptionTimeout)
	defer cancel()

	msgs := []string{}
	err = subscription.Receive(subCtx, func(ctx context.Context, msg *pubsub.Message) {
		msgs = append(msgs, string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		fmt.Printf("Error while receiving messages: %v\n", err)
		return nil, err
	}

	return msgs, nil
}

// Purge messages from given subscription names in parallel
func purgeMessages(step *structs.Step) *structs.StepResponse {
	purgeReq := gcp_pubsub.PurgeSubscriptionsRequest{}
	_ = utils.ParseInterface(step.Meta, &purgeReq)
	projectID := purgeReq.ProjectId

	wg := sync.WaitGroup{}
	var err error
	for _, subscriptionID := range purgeReq.SubscriptionNames {
		go func() {
			_, er := fetchPubsubMessages(projectID, subscriptionID)
			if er != nil && err == nil {
				err = er
			}
			wg.Add(1)
		}()
	}
	wg.Wait()
	if err != nil {
		return createDefaultErrorResponse(step, err)
	}

	return &structs.StepResponse{
		Name:     step.Name,
		Status:   structs.STEP_SUCCESS,
		Response: nil,
		Id:       step.Id,
	}
}

// Create a default step response with error object
func createDefaultErrorResponse(step *structs.Step, err error) *structs.StepResponse {
	return &structs.StepResponse{
		Name:   step.Name,
		Status: structs.STEP_ERROR,
		Response: &structs.StepError{
			Error: err.Error(),
		},
		Id: step.Id,
	}
}
