package utils

import "github.com/akshitbansal-1/async-testing/lib/structs"

// Create a default step response with error object
func CreateDefaultErrorResponse(step *structs.Step, err error) *structs.StepResponse {
	return &structs.StepResponse{
		Name:   step.Name,
		Status: structs.STEP_ERROR,
		Response: &structs.StepError{
			Error: err.Error(),
		},
		Id: step.Id,
	}
}
