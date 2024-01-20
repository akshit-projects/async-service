package kafka_validator

import (
	"errors"
	"fmt"

	"github.com/akshitbansal-1/async-testing/lib/structs"
	kafka "github.com/akshitbansal-1/async-testing/lib/structs/kafka"
	"github.com/akshitbansal-1/async-testing/lib/utils"
)

// ValidateKafkaConfig validates the KafkaConfig struct.
func ValidateKafkaConfig(config kafka.KafkaConfig) error {
	if len(config.BootstrapServers) == 0 {
		return errors.New("BootstrapServers cannot be empty")
	}

	return nil
}

// ValidatePublishRequest validates the PublishRequest struct.
func ValidateSubscribeRequest(step *structs.Step) error {
	meta := step.Meta
	var request kafka.SubscribeRequest
	if err := utils.ParseInterface[kafka.SubscribeRequest](meta, &request); err != nil {
		return errors.New("Unable to get subscribe request data")
	}

	step.Meta = &request

	if err := ValidateKafkaConfig(request.KafkaConfig); err != nil {
		return fmt.Errorf("Invalid KafkaConfig: %v", err)
	}

	if request.TopicName == "" {
		return errors.New("TopicName cannot be empty")
	}

	if request.GroupId == "" {
		return errors.New("GroupId cannot be empty")
	}

	return nil
}

// ValidateSubscribeRequest validates the SubscribeRequest struct.
func ValidatePublishRequest(step *structs.Step) error {
	meta := step.Meta
	var request kafka.PublishRequest
	if err := utils.ParseInterface[kafka.PublishRequest](meta, &request); err != nil {
		fmt.Println(err.Error())
		return errors.New("Unable to get publish request data")
	}

	step.Meta = &request

	if err := ValidateKafkaConfig(request.KafkaConfig); err != nil {
		return err
	}

	if request.TopicName == "" {
		return errors.New("TopicName cannot be empty")
	}

	if len(request.Messages) == 0 {
		return errors.New("Messages cannot be empty")
	}

	for idx, message := range request.Messages {
		if message.Value == "" {
			return errors.New("Value cannot be empty for message: " + fmt.Sprint(idx+1))
		}
	}

	return nil
}
