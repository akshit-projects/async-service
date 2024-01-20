package kafka

// KafkaConfig represents the configuration for connecting to Kafka.
type KafkaConfig struct {
	BootstrapServers []string
	Config           map[string]interface{}
}

type KafkaMessage struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PublishRequest struct {
	KafkaConfig KafkaConfig    `json:"kafkaConfig"`
	TopicName   string         `json:"topicName" bson:"topicName"`
	Messages    []KafkaMessage `json:"messages"`
}

type SubscribeRequest struct {
	KafkaConfig   KafkaConfig `json:"kafkaConfig"`
	TopicName     string      `json:"topicName" bson:"topicName"`
	GroupId       string      `json:"groupId"`
	MaxMessages   int         `json:"maxMessages"`
	FromBeginning bool        `json:"fromBeginning"`
}
