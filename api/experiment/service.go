package experiment_apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/config"
)

type Service interface {
	GetExperimentsData(teamName string) (*GetExperimentsResponse, error)
	GetLayerExperimentsData(teamName string) (*GetExperimentsResponse, error)
	OverrideVariant(req *OverrideVariantRequest) error
}

type service struct {
	app app.App
}

func NewService(app app.App) Service {
	return &service{app}
}

// GetLayerExperimesData implements Service.
func (s *service) GetLayerExperimentsData(teamName string) (*GetExperimentsResponse, error) {
	config := s.app.GetConfig()
	url := fmt.Sprintf("%s/experiments", config.ExperimentationServiceURL)
	layersRes, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var layersExpResponse []ExperimentServiceResponse
	if err = parseResponse(layersRes, layersExpResponse); err != nil {
		return nil, err
	}

	teamExperiments := []ExperimentData{}
	addE13nResponse(layersExpResponse, teamExperiments, teamName)

	return &GetExperimentsResponse{
		teamExperiments,
	}, nil
}

// GetExperimentsData implements Service.
func (s *service) GetExperimentsData(teamName string) (*GetExperimentsResponse, error) {
	config := s.app.GetConfig()
	url := fmt.Sprintf("%s/experiments", config.ExperimentationServiceURL)
	expRes, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if expRes.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to get data from e13n service")
	}

	var experimentsResponse []ExperimentServiceResponse
	if err = parseResponse(expRes, experimentsResponse); err != nil {
		return nil, err
	}

	teamExperiments := []ExperimentData{}
	addE13nResponse(experimentsResponse, teamExperiments, teamName)

	return &GetExperimentsResponse{
		teamExperiments,
	}, nil
}

func parseResponse(res *http.Response, experimentsResponse []ExperimentServiceResponse) error {
	if data, err := io.ReadAll(res.Body); err != nil {
		return errors.New("Invalid data from e13n service")
	} else {
		if err = json.Unmarshal(data, &experimentsResponse); err != nil {
			return errors.New("Unable to read e13n service response")
		}
	}
	return nil
}

func addE13nResponse(experimentsResponse []ExperimentServiceResponse, teamExperiments []ExperimentData, teamName string) {
	for _, experiment := range experimentsResponse {
		if experiment.Tag == teamName {
			expData := ExperimentData{
				experiment.Name,
				experiment.Id,
				experiment.variants,
				"exp",
				"",
			}
			teamExperiments = append(teamExperiments, expData)
		}
	}
}

// OverrideVariant implements Service.
func (s *service) OverrideVariant(req *OverrideVariantRequest) error {
	if err := ValidateExperimentOverride(req); err != nil {
		return err
	}

	config := s.app.GetConfig()
	if req.LayerId == "" {
		return overrideSimpleExperiment(config, req)
	} else {
		return OverrideLayerExperiment(config, req)
	}
}

func overrideSimpleExperiment(c *config.Configuration, req *OverrideVariantRequest) error {
	return nil
}

func OverrideLayerExperiment(c *config.Configuration, req *OverrideVariantRequest) error {
	url := fmt.Sprintf("%s/experiments", c.ExperimentationServiceURL)
	jsonBody := []byte("{}")
	reader := bytes.NewReader(jsonBody)
	res, err := http.Post(url, "application/json", reader)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusOK {
		return nil
	}

	return errors.New("Unable to update the layers experiment variant")
}
