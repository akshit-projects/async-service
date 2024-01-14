package structs

type StepStatus string
type ExecutionStatusType string

const (
	STEP_SUCCESS StepStatus = "SUCCESS"
	STEP_ERROR   StepStatus = "ERROR"
)

const (
	EXECUTION_STATUS_ERROR   ExecutionStatusType = "ERROR"
	EXECUTION_STATUS_MESSAGE ExecutionStatusType = "MESSAGE"
	EXECUTION_STATUS_SR      ExecutionStatusType = "SR"
)

type Flow struct {
	Name       string `json:"name"`
	Id         string `json:"id" bson:"_id,omitempty"`
	Creator    string `json:"email"`
	Steps      []Step `json:"steps"`
	CreatedAt  int64  `json:"createdAt"`
	ModifiedAt int64  `json:"modifiedAt"`
}

type Step struct {
	Name     string      `json:"name"`
	Function string      `json:"function"`
	Type     string      `json:"type,omitempty"`
	Meta     interface{} `json:"meta"`
	Value    interface{} `json:"value" bson:"-"`
	Id       string      `json:"id" bson:"-"`
}

type ExecutionStatusUpdate struct {
	Type    ExecutionStatusType `json:"type"`
	SR      *StepResponse       `json:"stepResponse"`
	Message string              `json:"message"`
}

type StepResponse struct {
	Name     string      `json:"name,omitempty"`
	Status   StepStatus  `json:"status"`
	Response interface{} `json:"response"`
	Id       string      `json:"id"`
}

type StepError struct {
	Expected string `json:"expected,omitempty"`
	Actual   string `json:"actual,omitempty"`
	Error    string `json:"error,omitempty"`
}
