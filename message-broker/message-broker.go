package message_broker

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akshitbansal-1/async-testing/worker/config"
	"github.com/akshitbansal-1/async-testing/worker/scheduler"
	thirdparty "github.com/akshitbansal-1/async-testing/worker/third_party"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var logger = thirdparty.Logger

func InitConsumer(config *config.Configuration, sc scheduler.Scheduler) {
	kafkaConfig := config.BrokerConfiguration
	conf := kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.Brokers,
		"group.id":          kafkaConfig.GroupId,
	}
	topic := kafkaConfig.Topic
	pullTimeout := kafkaConfig.PullTimeoutMs

	c, err := kafka.NewConsumer(&conf)

	if err != nil {
		logger.Fatal("Failed to create consumer: " + err.Error())
		os.Exit(1)
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	sigchan := setupShutdown()
	logger.Info("Reading the messages")

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			logger.Error("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := c.ReadMessage(time.Duration(pullTimeout) * time.Millisecond)
			if err != nil {
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					// logger.Info("No new message found")
					continue
				}
				logger.Error("Error while reading message from kafka: " + err.Error())
			} else {
				sc.ProcessMessage(ev.Value)
			}
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
