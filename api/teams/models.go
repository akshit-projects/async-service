package teams_api

import "go.mongodb.org/mongo-driver/bson/primitive"

type Team struct {
	Org      string             `json:"org"`
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `json:"name"`
	Platform string             `json:"platform,omitempty"`
	Members  []string           `json:"members"`
}

type AddMemberRequest struct {
	TeamId    string `json:"teamId"`
	UserEmail string `json:"userEmail"`
}
