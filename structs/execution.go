package structs

type ExecutionStatus string

const (
	FLOW_SUBMITTED ExecutionStatus = "FLOW_SUBMITTED"
)

type Execution struct {
	Id     string          `json:"id" bson:"_id,omitempty"`
	Status ExecutionStatus `json:"status" bson:"status"`
	Flow   *Flow           `json:"flow"`
}
