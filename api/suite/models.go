package suite_apis

type Suite struct {
	Id         string   `json:"id" bson:"_id,omitempty"`
	Name       string   `json:"name"`
	FlowIds    []string `json:"flowIds"`
	CreatedAt  int64    `json:"createdAt"`
	ModifiedAt int64    `json:"modifiedAt"`
}
