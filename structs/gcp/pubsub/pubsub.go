package gcp_pubsub

type PublishRequest struct {
	ProjectId string   `json:"projectId" bson:"projectId"`
	TopicName string   `json:"topicName" bson:"topicName"`
	Messages  []string `json:"messages"`
}

type PublishResponse struct {
	MessageIds []string `json:"messageIds"`
}

type SubscribeRequest struct {
	ProjectId        string `json:"projectId" bson:"projectId"`
	SubscriptionName string `json:"subscriptionName" bson:"subscriptionName"`
}

type SubscribeResponse struct {
	Messages []string `json:"messagess"`
}

type PurgeSubscriptionsRequest struct {
	ProjectId         string   `json:"projectId"`
	SubscriptionNames []string `json:"subscriptions"`
}
