package config

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/akshitbansal-1/async-testing/be/constants"
	"github.com/joho/godotenv"
)

type RedisConfiguration struct {
	Hosts    string `json:"hosts"`
	Password string `json:"password"`
}

type MessageBrokerConfig struct {
	Brokers       string `json:"brokers"`
	GroupId       string `json:"groupId"`
	Topic         string `json:"topic"`
	PullTimeoutMs int    `json:"pullTimeoutMS"`
}

type Configuration struct {
	MongoConnectionString     string             `json:"mongoConnectionString"`
	RedisConfiguration        RedisConfiguration `json:"redisConfiguration"`
	ExperimentationServiceURL string             `json:"experimentationServiceUrl"`
	JWTSecret                 string             `json:"jwtSecret"`
	GOOGLE_SIGN_IN_CLIENT_ID  string             `json:"googleSignInClientId"`
	KafkaConfiguration        MessageBrokerConfig `json:"kafka"`
}

func NewConfig() *Configuration {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return getConfig()
}

func getConfig() *Configuration {
	env := os.Getenv(constants.ENVIRONMENT_KEY)
	if env == constants.DEV_ENVIRONMENT_KEY {
		fileData := readFile("./config/config-stag.json")
		return parseConfig(fileData)
	}

	fileData := readFile("./config/config-prod.json")
	return parseConfig(fileData)
}

func parseConfig(fileData []byte) *Configuration {
	config := &Configuration{}
	err := json.Unmarshal(fileData, config)
	if err != nil {
		log.Fatal("Unable to parse config file")
	}

	return config
}

func readFile(fileName string) []byte {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Unable to read config file")
	}

	return byteValue
}
