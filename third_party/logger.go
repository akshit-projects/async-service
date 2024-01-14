package thirdparty

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var logPath = "/Users/akshitbansal/Developer/data/logstash/async-service/logs.log" // change this

func createLogger() *logrus.Logger {
	logger := logrus.New()

	// Create a console formatter (you can use other formatters as needed)
	consoleFormatter := &logrus.TextFormatter{
		ForceColors:            false,
		DisableTimestamp:       false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	}

	// Set the formatter for console output
	logger.SetFormatter(consoleFormatter)

	// Create a file output and set the formatter
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal("Error opening log file:", err)
	}

	// Set the file formatter
	fileWriter := logrus.New()
	fileWriter.SetOutput(file)

	// Create a multi-writer that writes to both console and file
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Set the multiple outputs for the logger
	logger.SetOutput(multiWriter)

	return logger
}

var Logger = createLogger()
