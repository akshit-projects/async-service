package gcp_pubsub_validator

import (
	"errors"

	"github.com/akshitbansal-1/async-testing/lib/structs"
	gcp_pubsub "github.com/akshitbansal-1/async-testing/lib/structs/gcp/pubsub"
	"github.com/akshitbansal-1/async-testing/lib/utils"
)

// Pubsub purge validation block
func ValidatePurgeSubscriptions(step *structs.Step) error {
	meta := step.Meta
	var purgeRequest gcp_pubsub.PurgeSubscriptionsRequest
	if err := utils.ParseInterface[gcp_pubsub.PurgeSubscriptionsRequest](meta, &purgeRequest); err != nil {
		return errors.New("Unable to get purge request data")
	}

	if purgeRequest.ProjectId == "" {
		return errors.New("Project id is required for subscription step")
	}
	if len(purgeRequest.SubscriptionNames) == 0 {
		return errors.New("Atleast one subscription is required for purging")
	}
	step.Meta = &purgeRequest

	return nil
}

// Pubsub subscribe validation block
func ValidatePubsubSubscribe(step *structs.Step) error {
	meta := step.Meta
	var subscribeRequest gcp_pubsub.SubscribeRequest
	if err := utils.ParseInterface[gcp_pubsub.SubscribeRequest](meta, &subscribeRequest); err != nil {
		return errors.New("Unable to get subscribe request data")
	}

	if subscribeRequest.ProjectId == "" {
		return errors.New("Project id is required for subscription step")
	}
	if subscribeRequest.SubscriptionName == "" {
		return errors.New("Subscription name is required for subscription step")
	}
	step.Meta = &subscribeRequest

	return nil
}

// Pubsub publish validation block
func ValidatePubsubPublish(step *structs.Step) error {
	meta := step.Meta
	var publishRequest gcp_pubsub.PublishRequest
	if err := utils.ParseInterface[gcp_pubsub.PublishRequest](meta, &publishRequest); err != nil {
		return errors.New("Unable to get publish request data")
	}

	if publishRequest.ProjectId == "" {
		return errors.New("Project id required for publish request")
	}
	if publishRequest.TopicName == "" {
		return errors.New("Topic name is required for publish request")
	}
	if len(publishRequest.Messages) == 0 {
		return errors.New("At least one message is required for publish request")
	}
	step.Meta = &publishRequest

	return nil
}
