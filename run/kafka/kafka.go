package run_kafka

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/akshitbansal-1/async-testing/lib/structs"
	kafka_struct "github.com/akshitbansal-1/async-testing/lib/structs/kafka"
	"github.com/akshitbansal-1/async-testing/lib/utils"
	thirdparty "github.com/akshitbansal-1/async-testing/worker/third_party"
	worker_utils "github.com/akshitbansal-1/async-testing/worker/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var logger = thirdparty.Logger

// SendMessage sends a message to a Kafka topic.
func SendMessage(step *structs.Step) *structs.StepResponse {
	var request kafka_struct.PublishRequest
	err := utils.ParseInterface[kafka_struct.PublishRequest](step.Meta, &request)
	if err != nil {
		logger.Error("Error while decoding kafka publish request", err)
		return &structs.StepResponse{
			Name:     "",
			Status:   structs.STEP_ERROR,
			Response: "Error while decoding kafka publish request",
			Id:       "",
		}
	}

	// Create Kafka producer configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(request.KafkaConfig.BootstrapServers, ","),
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		logger.Error("Error while creating kafka producer", err)
		return worker_utils.CreateDefaultErrorResponse(step, err)
	}
	defer producer.Close()

	// Produce messages
	delChannels := make(chan kafka.Event, len(request.Messages))
	for _, msg := range request.Messages {
		message := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &request.TopicName, Partition: kafka.PartitionAny},
			Key:            []byte(msg.Key),
			Value:          []byte(msg.Value),
		}

		err := producer.Produce(message, delChannels)
		if err != nil {
			logger.Error("Unable to publish a message", err)
			return worker_utils.CreateDefaultErrorResponse(step, err)
		}
	}
	<-delChannels
	close(delChannels)

	return &structs.StepResponse{
		Name:     step.Name,
		Status:   structs.STEP_SUCCESS,
		Response: nil,
		Id:       step.Id,
	}
}

// ConsumeMessages consumes messages from a Kafka topic.
func ConsumeMessages(step *structs.Step) *structs.StepResponse {
	var request kafka_struct.SubscribeRequest
	err := utils.ParseInterface[kafka_struct.SubscribeRequest](step.Meta, &request)
	if err != nil {
		logger.Error("Error while decoding kafka subscribe request", err)
		return &structs.StepResponse{
			Name:     "",
			Status:   structs.STEP_ERROR,
			Response: "Error while decoding kafka subscribe request",
			Id:       "",
		}
	}

	config := kafka.ConfigMap{
		"bootstrap.servers": strings.Join(request.KafkaConfig.BootstrapServers, ","),
	}

	if request.GroupId != "" {
		config.SetKey("group.id", request.GroupId)
	}

	if request.FromBeginning {
		config.SetKey("auto.offset.reset", kafka.OffsetBeginning.String())
	}

	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		logger.Error("Error while creating kafka consumer", err)
		return worker_utils.CreateDefaultErrorResponse(step, err)
	}

	// Subscribe to topic
	err = consumer.SubscribeTopics([]string{request.TopicName}, nil)
	if err != nil {
		logger.Error("Error while subscribing to kafka topics", err)
		return worker_utils.CreateDefaultErrorResponse(step, err)
	}
	defer consumer.Close()

	var er error
	var messages []kafka_struct.KafkaMessage
	utils.Race(context.Background(), func() {
		i := 0
		startTime := time.Now().UnixMilli()
		for {
			fmt.Println("reading a message")
			if i >= request.MaxMessages {
				break
			}
			if startTime+int64(step.Timeout) < time.Now().UnixMilli() {
				break
			}
			i += 1
			msg, err := consumer.ReadMessage(time.Second * 2)
			if err != nil {
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				er = errors.New("Error while reading kafka messages")
				logger.Error("Error while reading kafka messages: ", err)
				break
			}
			messages = append(messages, kafka_struct.KafkaMessage{
				Value: string(msg.Value),
				Key:   string(msg.Key),
			})
			time.Sleep(time.Second)
		}
	}, step.Timeout)

	if er != nil {
		return worker_utils.CreateDefaultErrorResponse(step, er)
	}

	return &structs.StepResponse{
		Name:     step.Name,
		Status:   structs.STEP_SUCCESS,
		Response: messages,
		Id:       step.Id,
	}
}
