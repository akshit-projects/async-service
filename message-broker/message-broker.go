package message_broker

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akshitbansal-1/async-testing/worker/config"
	"github.com/akshitbansal-1/async-testing/worker/scheduler"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func InitConsumer(config *config.Configuration, sc scheduler.Scheduler) {
	conf := kafka.ConfigMap{}
	kafkaConfig := config.BrokerConfiguration
	conf["bootstrap.servers"] = kafkaConfig.Brokers
	conf["group.id"] = kafkaConfig.GroupId
	topic := kafkaConfig.Topic
	pullTimeout := kafkaConfig.PullTimeoutMs

	c, err := kafka.NewConsumer(&conf)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	sigchan := setupShutdown()

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := c.ReadMessage(time.Duration(pullTimeout) * time.Millisecond)
			if err != nil {
				continue
			}
			sc.ProcessMessage(ev.Value)
		}
	}

	c.Close()
}

func setupShutdown() chan os.Signal {
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	return sigchan
}
