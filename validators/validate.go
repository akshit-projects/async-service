package validators

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/akshitbansal-1/async-testing/lib/constants"
	"github.com/akshitbansal-1/async-testing/lib/structs"
	api_validator "github.com/akshitbansal-1/async-testing/lib/validators/api"
	kafka "github.com/akshitbansal-1/async-testing/lib/validators/kafka"
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
	fmt.Println(step.Function)
	switch step.Function {
	case constants.HTTP_API_STEP:
		return api_validator.ValidateHttpStep(step)
	case constants.PUBLISH_KAFKA_MESSAGE_STEP:
		return kafka.ValidatePublishRequest(step)
	case constants.SUBSCRIBE_KAFKA_MESSAGES_STEP:
		return kafka.ValidateSubscribeRequest(step)
	default:
		return errors.New("Unsupported function")
	}
}
