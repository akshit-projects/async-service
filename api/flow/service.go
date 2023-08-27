package flow_apis

import (
	"time"

	"github.com/akshitbansal-1/async-testing/be/app"
)

type Service interface {
	AddFlow(flow *Flow) (*string, error)
	UpdateFlow(flow *Flow) error
	ValidateSteps(flow *Flow) error
	RunFlow(ch chan<- *StepResponse, flow *Flow) error
	GetFlows() ([]Flow, error)
	GetFlow(id string) (*Flow, error)
}

type service struct {
	app app.App
}

func NewService(app app.App) Service {
	return &service{app}
}

func (s *service) GetFlow(id string) (*Flow, error) {
	return getFlow(s.app, id)
}

func (s *service) GetFlows() ([]Flow, error) {
	return getFlows(s.app)
}

func (s *service) UpdateFlow(flow *Flow) error {
	steps := flow.Steps
	if err := validateSteps(steps); err != nil {
		return err
	}

	flow.ModifiedAt = time.Now().Unix()
	return updateFlow(s.app, flow)
}

func (s *service) AddFlow(flow *Flow) (*string, error) {
	steps := flow.Steps
	if err := validateSteps(steps); err != nil {
		return nil, err
	}

	flow.CreatedAt = time.Now().Unix()
	flow.ModifiedAt = time.Now().Unix()
	return addFlow(s.app, flow)
}

func (s *service) ValidateSteps(flow *Flow) error {
	steps := flow.Steps
	err := validateSteps(steps)
	return err
}

func (s *service) RunFlow(ch chan<- *StepResponse, flow *Flow) error {
	if err := validateSteps(flow.Steps); err != nil {
		ch <- &StepResponse{
			"",
			"",
			ERROR,
			&StepError{
				Error: "Invalid steps data. " + err.Error(),
			},
			"",
		}
		close(ch)
		return nil
	}

	return RunFlow(ch, s.app, flow)
}
