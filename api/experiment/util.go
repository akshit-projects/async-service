package experiment_apis

import (
	"errors"
	"strconv"
)

func ValidateExperimentOverride(req *OverrideVariantRequest) error {
	if req.ExperimentId == "" {
		return errors.New("Experiment id not found in request")
	}

	if req.LayerId == "" {
		return errors.New("Layer id not found in request")
	}

	if req.Variant == "" {
		return errors.New("Variant not found in request")
	}

	if req.UserId == "" {
		return errors.New("User id not found in request")
	}

	if _, err := strconv.Atoi(req.UserId); err != nil {
		return errors.New("Invalid user id in request")
	}

	return nil
}
