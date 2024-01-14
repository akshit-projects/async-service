package validators

import (
	"errors"
	"strconv"

	"github.com/akshitbansal-1/async-testing/lib/constants"
	"github.com/akshitbansal-1/async-testing/lib/structs"
	api_validator "github.com/akshitbansal-1/async-testing/lib/validators/api"
	gcp_pubsub "github.com/akshitbansal-1/async-testing/lib/validators/gcp/pubsub"
)

func ValidateSteps(steps []structs.Step) error {
	if len(steps) == 0 {
		return errors.New("Minimum one step is required")
	}
	for idx, step := range steps {
		if step.Name == "" {
			return errors.New("Name is required for step " + strconv.Itoa(idx))
		}
		err := ValidateStep(&steps[idx])
		if err != nil {
			return errors.New(step.Name + " -> " + err.Error())
		}
	}
	return nil
}

func ValidateStep(step *structs.Step) error {
	switch step.Function {
	case constants.API_STEP:
		return api_validator.ValidateHttpStep(step)
	case constants.PUBLISH_MESSAGE_STEP:
		return gcp_pubsub.ValidatePubsubPublish(step)
	case constants.SUBSCRIBE_MESSAGES_STEP:
		return gcp_pubsub.ValidatePubsubSubscribe(step)
	case "purge-subscriptions":
		return gcp_pubsub.ValidatePurgeSubscriptions(step)
	default:
		return errors.New("Unsupported function")
	}
}
