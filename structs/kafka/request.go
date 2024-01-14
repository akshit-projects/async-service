package kafka

type PublishRequest struct {
	Name      string   `json:"name" bson:"name"`
	TopicName string   `json:"topicName" bson:"topicName"`
	Messages  []string `json:"messages"`
}

type SubscribeRequest struct {
	Name      string `json:"name" bson:"name"`
	TopicName string `json:"topicName" bson:"topicName"`
}
