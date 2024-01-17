package thirdparty

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/akshitbansal-1/async-testing/be/config"
	"github.com/akshitbansal-1/async-testing/lib/structs"
	"github.com/akshitbansal-1/async-testing/lib/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	ExecutionTopic = "executions"
)

type MessageBroker interface {
	PushExecution(config *config.Configuration, flow *structs.Execution) error
}

type broker struct {
	producer *kafka.Producer
}

func InitBroker(config *config.Configuration) MessageBroker {
	conf := kafka.ConfigMap{}
	kafkaConfig := config.KafkaConfiguration
	conf["bootstrap.servers"] = kafkaConfig.Brokers
	fmt.Print(kafkaConfig.Brokers)
	prd, err := kafka.NewProducer(&conf)
	if err != nil {
		fmt.Println("Unable to run kafka")
		os.Exit(1)
	}

	return &broker{
		prd,
	}
}

func (k *broker) PushExecution(config *config.Configuration, exec *structs.Execution) error {
	bytes, err := json.Marshal(exec)
	if err != nil {
		Logger.Error("Invalid execution data ", utils.StructToString(*exec))
		return err
	}

	topic := ExecutionTopic
	msg := &kafka.Message{
		Value:          bytes,
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	}

	completionSignal := make(chan kafka.Event, 1)
	isTimedout := utils.Race(context.Background(), func() {
		err = k.producer.Produce(msg, completionSignal)
		<-completionSignal
	}, 1000)

	if isTimedout {
		return errors.New("Unable to publish execution flow to kafka. API timed out")
	}

	return err
}

func setupShutdown() chan os.Signal {
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	return sigchan
}
