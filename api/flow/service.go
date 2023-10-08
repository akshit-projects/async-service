package flow_apis

import (
	"time"

	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/common_structs"
	"github.com/akshitbansal-1/async-testing/lib/structs"
)

type Service interface {
	AddFlow(flow *structs.Flow) (*string, error)
	UpdateFlow(flow *structs.Flow) error
	ValidateSteps(flow *structs.Flow) error
	RunFlow(ch chan<- *structs.StepResponse, flow *structs.Flow) error
	GetFlows(*common_structs.APIFilter) ([]structs.Flow, error)
	GetFlow(id string) (*structs.Flow, error)
}

type service struct {
	app app.App
}

func NewService(app app.App) Service {
	return &service{app}
}

func (s *service) GetFlow(id string) (*structs.Flow, error) {
	return getFlow(s.app, id)
}

func (s *service) GetFlows(filter *common_structs.APIFilter) ([]structs.Flow, error) {
	return getFlows(s.app, filter)
}

func (s *service) UpdateFlow(flow *structs.Flow) error {
	steps := flow.Steps
	if err := validateSteps(steps); err != nil {
		return err
	}

	flow.ModifiedAt = time.Now().Unix()
	return updateFlow(s.app, flow)
}

func (s *service) AddFlow(flow *structs.Flow) (*string, error) {
	steps := flow.Steps
	if err := validateSteps(steps); err != nil {
		return nil, err
	}

	flow.CreatedAt = time.Now().Unix()
	flow.ModifiedAt = time.Now().Unix()
	return addFlow(s.app, flow)
}

func (s *service) ValidateSteps(flow *structs.Flow) error {
	steps := flow.Steps
	err := validateSteps(steps)
	return err
}

func (s *service) RunFlow(ch chan<- *structs.StepResponse, flow *structs.Flow) error {
	if err := validateSteps(flow.Steps); err != nil {
		ch <- &structs.StepResponse{
			"",
			"",
			structs.ERROR,
			&structs.StepError{
				Error: "Invalid steps data. " + err.Error(),
			},
			"",
		}
		close(ch)
		return nil
	}

	_, err := StartFlow(ch, s.app, flow)
	// TODO start polling the status
	return err
}
