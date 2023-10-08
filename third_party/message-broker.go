package thirdparty

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/akshitbansal-1/async-testing/be/config"
	"github.com/akshitbansal-1/async-testing/lib/structs"
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
		fmt.Println("Invalid exec data")
		return err
	}

	topic := ExecutionTopic
	msg := &kafka.Message{
		Value:          bytes,
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	}
	delivery_chan := make(chan kafka.Event, 10000)
	return k.producer.Produce(msg, delivery_chan)
}

func setupShutdown() chan os.Signal {
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	return sigchan
}
