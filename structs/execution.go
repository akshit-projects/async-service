package structs

type ExecutionStatus string
type ExecutionStatusCode int

const (
	ES_MONGO_ERROR      = 110
	ES_KAFKA_ERROR      = 111
	ES_RT_UPDATES_ERROR = 112
)

const (
	EXECUTION_TODO    ExecutionStatus = "TODO"
	EXECUTION_QUEUED  ExecutionStatus = "QUEUED"
	EXECUTION_RUNNING ExecutionStatus = "RUNNING"
	EXECUTION_DONE    ExecutionStatus = "DONE"
	EXECUTION_FAILED  ExecutionStatus = "FAILED"
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

type ExecutionStatusUpdate struct {
	Type    ExecutionStatusType `json:"type"`
	SR      *StepResponse       `json:"stepResponse"`
	Message string              `json:"message"`
	Code    ExecutionStatusCode `json:"code"`
}
