package run_gcp_pubsub

import (
	"context"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/akshitbansal-1/async-testing/lib/structs"
	gcp_pubsub "github.com/akshitbansal-1/async-testing/lib/structs/gcp/pubsub"
	"github.com/akshitbansal-1/async-testing/lib/utils"
	thirdparty "github.com/akshitbansal-1/async-testing/worker/third_party"
	worker_utils "github.com/akshitbansal-1/async-testing/worker/utils"
)

var logger = thirdparty.Logger

// Publish messages in pubsub
func publishMessages(step *structs.Step) *structs.StepResponse {
	publishReq := gcp_pubsub.PublishRequest{}
	_ = utils.ParseInterface(step.Meta, &publishReq)
	projectID := publishReq.ProjectId
	topicID := publishReq.TopicName

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		logger.Error("Unable to connect to pubsub", err)
		return worker_utils.CreateDefaultErrorResponse(step, err)
	}

	topic := client.Topic(topicID)

	messageIds := []string{}
	for _, message := range publishReq.Messages {
		result := topic.Publish(ctx, &pubsub.Message{
			Data: []byte(message),
		})
		id, err := result.Get(ctx)
		if err != nil {
			logger.Error("Failed to publish message to pubsub", err)
			return worker_utils.CreateDefaultErrorResponse(step, err)
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
		logger.Error("Unable to fetch messages from pubsub", err)
		return worker_utils.CreateDefaultErrorResponse(step, err)
	}

	return &structs.StepResponse{
		Name:   step.Name,
		Status: structs.STEP_SUCCESS,
		Response: &gcp_pubsub.SubscribeResponse{
			Messages: msgs,
		},
		Id: step.Id,
	}
}

// Fetch messages of given subscription name
func fetchPubsubMessages(projectID string, subscriptionID string) ([]string, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		logger.Error("Unable to connect to pubsub", err)
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
		logger.Error("Error while receiving messages", err)
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
		logger.Error("Error while purging pubsub queue", err)
		return worker_utils.CreateDefaultErrorResponse(step, err)
	}

	return &structs.StepResponse{
		Name:     step.Name,
		Status:   structs.STEP_SUCCESS,
		Response: nil,
		Id:       step.Id,
	}
}
