package structs

type ExecutionStatus string

const (
	TODO    ExecutionStatus = "TODO"
	QUEUED  ExecutionStatus = "QUEUED"
	RUNNING ExecutionStatus = "RUNNING"
	DONE    ExecutionStatus = "DONE"
	FAILED  ExecutionStatus = "FAILED"
)

type ExecutionFlow struct {
	Name  string `json:"name"`
	Id    string `json:"id" bson:"_id,omitempty"`
	Steps []Step `json:"steps"`
}

type Execution struct {
	Id            string          `json:"id" bson:"_id,omitempty"`
	ExecutionFlow *ExecutionFlow  `json:"executionFlow"`
	Executor      string          `json:"executor"`
	TotalTimeout  int             `json:"totalTimeout"`
	Status        ExecutionStatus `json:"status" bson:"status"`
	CreatedAt     int64           `json:"createdAt"`
	ModifiedAt    int64           `json:"modifiedAt"`
}
