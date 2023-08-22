package flow_apis

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	experiment_apis "github.com/akshitbansal-1/async-testing/be/api/experiment"
	"github.com/akshitbansal-1/async-testing/be/constants"
	"github.com/akshitbansal-1/async-testing/be/utils"
)

var validMethods = [4]string{"GET", "POST", "PUT", "DELETE"}

func validateSteps(steps []Step) error {
	if len(steps) == 0 {
		return errors.New("Minimum one step is required")
	}
	for idx, step := range steps {
		if step.Name == "" {
			return errors.New("Name is required for step " + strconv.Itoa(idx))
		}
		err := validateStep(&steps[idx])
		if err != nil {
			return errors.New(step.Name + " -> " + err.Error())
		}
	}
	return nil
}

func validateStep(step *Step) error {
	switch step.Function {
	case constants.API_STEP:
		return validateHttpStep(step)
	case constants.PUBLISH_MESSAGE_STEP:
		return validatePubsubPublish(step)
	case constants.SUBSCRIBE_MESSAGES_STEP:
		return validatePubsubSubscribe(step)
	case "purge-subscriptions":
		return validatePurgeSubscriptions(step)
	case "override-variant":
		return validateOverrideVariant(step)
	default:
		return errors.New("Unsupported function")
	}
}

// E13n variant override request block
func validateOverrideVariant(step *Step) error {
	meta := step.Meta
	var overrideReq experiment_apis.OverrideVariantRequest
	if err := utils.ParseInterface[experiment_apis.OverrideVariantRequest](meta, &overrideReq); err != nil {
		return err
	}

	if err := experiment_apis.ValidateExperimentOverride(&overrideReq); err != nil {
		return err
	}

	step.Meta = &overrideReq

	return nil
}

// Pubsub purge validation block
func validatePurgeSubscriptions(step *Step) error {
	meta := step.Meta
	var purgeRequest PurgeSubscriptionsRequest
	if err := utils.ParseInterface[PurgeSubscriptionsRequest](meta, &purgeRequest); err != nil {
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
func validatePubsubSubscribe(step *Step) error {
	meta := step.Meta
	var subscribeRequest SubscribeRequest
	if err := utils.ParseInterface[SubscribeRequest](meta, &subscribeRequest); err != nil {
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
func validatePubsubPublish(step *Step) error {
	meta := step.Meta
	var publishRequest PublishRequest
	if err := utils.ParseInterface[PublishRequest](meta, &publishRequest); err != nil {
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

// HTTP validation block
func validateHttpStep(step *Step) error {
	meta := step.Meta
	var httpReq HTTPRequest
	if err := utils.ParseInterface[HTTPRequest](meta, &httpReq); err != nil {
		return errors.New("Unable to get http request data")
	}

	httpReq.Method = strings.ToUpper(httpReq.Method)
	if err := validateHTTPMethod(httpReq.Method); err != nil {
		return err
	}

	// Check for specific validation conditions
	if strings.ToUpper(httpReq.Method) == "GET" && httpReq.Body != nil {
		return errors.New("Body can't go with GET method")
	} else if strings.ToUpper(httpReq.Method) != "GET" && httpReq.Body == nil {
		return errors.New("Body is required for " + httpReq.Method)
	}

	_, err := url.ParseRequestURI(httpReq.Url)
	if err != nil {
		return errors.New("Invalid request URL passed")
	}

	step.Value = &httpReq
	return nil
}

func validateHTTPMethod(method string) error {
	for _, m := range &validMethods {
		if method == m {
			return nil
		}
	}

	return errors.New("Invalid http method provided")
}
