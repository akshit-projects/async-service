package structs

type StepStatus string

const (
	SUCCESS StepStatus = "SUCCESS"
	ERROR   StepStatus = "ERROR"
)

type Flow struct {
	Name       string `json:"name"`
	Id         string `json:"id" bson:"_id,omitempty"`
	Creator    string `json:"email"`
	Steps      []Step `json:"steps"`
	CreatedAt  int64  `json:"createdAt"`
	ModifiedAt int64  `json:"modifiedAt"`
	TeamId     string `json:"teamId"`
}

type Step struct {
	Name     string      `json:"name"`
	Function string      `json:"function"`
	Type     string      `json:"type,omitempty"`
	Meta     interface{} `json:"meta"`
	Value    interface{} `json:"value" bson:"-"`
	Id       string      `json:"id" bson:"-"`
}

type StepError struct {
	Expected string `json:"expected,omitempty"`
	Actual   string `json:"actual,omitempty"`
	Error    string `json:"error,omitempty"`
}

type StepResponse struct {
	Name     string      `json:"name,omitempty"`
	Function string      `json:"function,omitempty"`
	Status   StepStatus  `json:"status"`
	Response interface{} `json:"response"`
	Id       string      `json:"id"`
}


