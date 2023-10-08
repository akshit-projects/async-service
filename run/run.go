package executor

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

	"cloud.google.com/go/pubsub"
	"github.com/akshitbansal-1/async-testing/lib/structs"
	"github.com/akshitbansal-1/async-testing/lib/utils"
	"github.com/akshitbansal-1/async-testing/worker/constants"
)

// Run flow step by step
func RunFlow(ch chan int, exec *structs.Execution) error {
	steps := exec.Flow.Steps
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

		// ch <- stepResponse
		if stepResponse.Status != structs.SUCCESS {
			<-ch
			return nil
		}
	}
	fmt.Println(<-ch)

	return nil
}

// Make HTTP call
func makeAPICall(step *structs.Step) *structs.StepResponse {
	h := structs.HTTPRequest{}
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
	stepValue := &structs.HTTPResponse{
		Status:   resp.StatusCode,
		Response: buf.String(),
	}

	if h.ExpectedStatus != "" && h.ExpectedStatus !=
		strconv.Itoa(resp.StatusCode) {
		return &structs.StepResponse{
			step.Name,
			step.Function,
			structs.ERROR,
			&structs.StepError{
				h.ExpectedStatus,
				strconv.Itoa(resp.StatusCode),
				"Status code not matching",
			},
			step.Id,
		}
	}

	if err = utils.CompareStrings(&stepValue.Response, &h.ExpectedResponse); err != nil {
		return &structs.StepResponse{
			step.Name,
			step.Function,
			structs.ERROR,
			&structs.StepError{
				h.ExpectedResponse,
				stepValue.Response,
				"Response not matching",
			},
			step.Id,
		}
	}

	return &structs.StepResponse{
		step.Name,
		step.Function,
		structs.SUCCESS,
		stepValue,
		step.Id,
	}
}

// Publish messages in pubsub
func publishMessages(step *structs.Step) *structs.StepResponse {
	publishReq := structs.PublishRequest{}
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
			fmt.Println("Failed to publish: %v", err)
			return createDefaultErrorResponse(step, err)
		}
		messageIds = append(messageIds, id)
	}

	return &structs.StepResponse{
		step.Name,
		step.Function,
		structs.SUCCESS,
		&structs.PublishResponse{
			messageIds,
		},
		step.Id,
	}
}

// Subscribe to messages in pubsub
func subscribeMessages(step *structs.Step) *structs.StepResponse {
	subscribeReq := structs.SubscribeRequest{}
	_ = utils.ParseInterface(step.Meta, &subscribeReq)
	projectID := subscribeReq.ProjectId
	subscriptionID := subscribeReq.SubscriptionName

	msgs, err := fetchPubsubMessages(projectID, subscriptionID)
	if err != nil {
		return createDefaultErrorResponse(step, err)
	}

	return &structs.StepResponse{
		step.Name,
		step.Function,
		structs.SUCCESS,
		&structs.SubscribeResponse{
			msgs,
		},
		step.Id,
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
	purgeReq := structs.PurgeSubscriptionsRequest{}
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
		step.Name,
		step.Function,
		structs.SUCCESS,
		nil,
		step.Id,
	}
}

// Create a default step response with error object
func createDefaultErrorResponse(step *structs.Step, err error) *structs.StepResponse {
	return &structs.StepResponse{
		step.Name,
		step.Function,
		structs.ERROR,
		&structs.StepError{
			Error: err.Error(),
		},
		step.Id,
	}
}
