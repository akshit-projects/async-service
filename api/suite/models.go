package suite_apis

type Suite struct {
	Name       string   `json:"name"`
	FlowIds    []string `json:"flowIds"`
	CreatedAt  int64    `json:"createdAt"`
	ModifiedAt int64    `json:"modifiedAt"`
}
