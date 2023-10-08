package flow_apis

import (
	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/lib/structs"
)

// Run flow step by step
func StartFlow(ch chan<- *structs.StepResponse, app app.App, flow *structs.Flow) (*string, error) {
	execution, err := submitFlow(app, flow)
	if err == nil {
		err = app.GetMessageBroker().PushExecution(app.GetConfig(), execution)
	}
	return &execution.Id, err
}
