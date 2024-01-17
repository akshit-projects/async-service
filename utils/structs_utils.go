package utils

import (
	"github.com/akshitbansal-1/async-testing/lib/structs"
)

func CreateErrorExecutionStatus(msg string, execCode structs.ExecutionStatusCode) *structs.ExecutionStatusUpdate {
	return &structs.ExecutionStatusUpdate{
		Type:    structs.EXECUTION_STATUS_ERROR,
		Message: msg,
		Code:    execCode,
	}
}
