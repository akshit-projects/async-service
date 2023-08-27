package flow_apis

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
	experiment_apis "github.com/akshitbansal-1/async-testing/be/api/experiment"
	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/constants"
	"github.com/akshitbansal-1/async-testing/be/utils"
)

// Run flow step by step
func RunFlow(ch chan<- *StepResponse, app app.App, flow *Flow) error {
	steps := flow.Steps
	for idx := range steps {
		step := &steps[idx]
		var stepResponse *StepResponse
		switch step.Function {
		case constants.API_STEP:
			stepResponse = makeAPICall(app, step)
		case constants.PUBLISH_MESSAGE_STEP:
			stepResponse = publishMessages(app, step)
		case constants.SUBSCRIBE_MESSAGES_STEP:
			stepResponse = subscribeMessages(app, step)
		case "purge-subscriptions":
			stepResponse = purgeMessages(app, step)
		case "override-variant":
			stepResponse = overrideVariant(app, step)
		default:
			stepResponse = createDefaultErrorResponse(step, errors.New("Unsupported function"))
		}

		ch <- stepResponse
		if stepResponse.Status != SUCCESS {
			close(ch)
			return nil
		}
	}
	close(ch)

	return nil
}

// Make HTTP call
func makeAPICall(app app.App, step *Step) *StepResponse {
	var request = step.Value.(*HTTPRequest)
	req, _ := http.NewRequest(request.Method, request.Url, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	var isGetCall = request.Method == "GET"
	if !isGetCall {
		bodyBytes, err := json.Marshal(request.Body)
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
	stepValue := &HTTPResponse{
		Status:   resp.StatusCode,
		Response: buf.String(),
	}

	if request.ExpectedStatus != "" && request.ExpectedStatus !=
		strconv.Itoa(resp.StatusCode) {
		return &StepResponse{
			step.Name,
			step.Function,
			ERROR,
			&StepError{
				request.ExpectedStatus,
				strconv.Itoa(resp.StatusCode),
				"Status code not matching",
			},
			step.Id,
		}
	}

	if err = utils.CompareStrings(&stepValue.Response, &request.ExpectedResponse); err != nil {
		return &StepResponse{
			step.Name,
			step.Function,
			ERROR,
			&StepError{
				request.ExpectedResponse,
				stepValue.Response,
				"Response not matching",
			},
			step.Id,
		}
	}

	return &StepResponse{
		step.Name,
		step.Function,
		SUCCESS,
		stepValue,
		step.Id,
	}
}

// Publish messages in pubsub
func publishMessages(app app.App, step *Step) *StepResponse {
	publishReq := step.Meta.(*PublishRequest)
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

	return &StepResponse{
		step.Name,
		step.Function,
		SUCCESS,
		&PublishResponse{
			messageIds,
		},
		step.Id,
	}
}

// Subscribe to messages in pubsub
func subscribeMessages(app app.App, step *Step) *StepResponse {
	subscribeReq := step.Meta.(*SubscribeRequest)
	projectID := subscribeReq.ProjectId
	subscriptionID := subscribeReq.SubscriptionName

	msgs, err := fetchPubsubMessages(projectID, subscriptionID)
	if err != nil {
		return createDefaultErrorResponse(step, err)
	}

	return &StepResponse{
		step.Name,
		step.Function,
		SUCCESS,
		&SubscribeResponse{
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
func purgeMessages(app app.App, step *Step) *StepResponse {
	purgeReq := step.Meta.(*PurgeSubscriptionsRequest)
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

	return &StepResponse{
		step.Name,
		step.Function,
		SUCCESS,
		nil,
		step.Id,
	}
}

// Override variant in layers for the user
func overrideVariant(app app.App, step *Step) *StepResponse {
	cfg := app.GetConfig()
	layerOverrideReq := step.Meta.(*experiment_apis.OverrideVariantRequest)
	err := experiment_apis.OverrideLayerExperiment(cfg, layerOverrideReq)
	if err != nil {
		return createDefaultErrorResponse(step, err)
	}

	return &StepResponse{
		step.Name,
		step.Function,
		SUCCESS,
		nil,
		step.Id,
	}
}

// Create a default step response with error object
func createDefaultErrorResponse(step *Step, err error) *StepResponse {
	return &StepResponse{
		step.Name,
		step.Function,
		ERROR,
		&StepError{
			Error: err.Error(),
		},
		step.Id,
	}
}
