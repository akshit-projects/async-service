package suite_apis

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/common_structs"
	"github.com/gofiber/fiber/v2"
)

const MAX_SUITE_FLOWS = 10

type Service interface {
	AddSuite(suite *Suite) (*string, *common_structs.APIError)
	GetSuites(*common_structs.APIFilter) ([]Suite, error)
}

type service struct {
	app app.App
}

func NewService(app app.App) Service {
	return &service{app}
}

func (s *service) GetSuites(filter *common_structs.APIFilter) ([]Suite, error) {
	return getSuites(s.app, filter)
}

func (s *service) AddSuite(suite *Suite) (*string, *common_structs.APIError) {
	if err := validateSuite(suite); err != nil {
		return nil, &common_structs.APIError{
			Status: http.StatusBadRequest,
			Msg:    err.Error(),
		}
	}

	id, err := addSuite(s.app, suite)
	if err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			return nil, &common_structs.APIError{
				Status: http.StatusBadRequest,
				Msg:    "Can't add non-existing flow ids",
			}
		} else {
			return nil, &common_structs.APIError{
				Status: http.StatusInternalServerError,
				Msg:    err.Error(),
			}
		}
	} else {
		return id, nil
	}
}

func validateSuite(suite *Suite) error {
	if suite.Name == "" {
		return errors.New("Suite Name is required")
	}

	if len(suite.FlowIds) == 0 {
		return errors.New("Atleast one flow id is required")
	}

	if len(suite.FlowIds) > MAX_SUITE_FLOWS {
		return errors.New(fmt.Sprintf("Atmost %d flow ids can be saved in suite", MAX_SUITE_FLOWS))
	}

	return nil
}
