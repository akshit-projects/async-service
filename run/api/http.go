package run_http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/akshitbansal-1/async-testing/lib/structs"
	"github.com/akshitbansal-1/async-testing/lib/structs/api"
	"github.com/akshitbansal-1/async-testing/lib/utils"
	lib_utils "github.com/akshitbansal-1/async-testing/lib/utils"
	thirdparty "github.com/akshitbansal-1/async-testing/worker/third_party"
	worker_utils "github.com/akshitbansal-1/async-testing/worker/utils"
)

var logger = thirdparty.Logger

// Make HTTP call
func MakeAPICall(step *structs.Step) *structs.StepResponse {
	h := api.HTTPRequest{}
	_ = lib_utils.ParseInterface(step.Meta, &h)
	req, _ := http.NewRequest(h.Method, h.Url, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}

	var isGetCall = h.Method == "GET"
	if !isGetCall {
		bodyBytes, err := json.Marshal(h.Body)
		if err != nil {
			logger.Error("Error while marshaling HTTP request body: ", err.Error())
			return worker_utils.CreateDefaultErrorResponse(step, err)
		}
		// set body
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	// Perform the HTTP request
	resp, err := lib_utils.CallHTTP(req)
	if err != nil {
		logger.Println("HTTP request error: ", err)
		if errors.Is(err, context.DeadlineExceeded) {
			err = errors.New("Request timed out")
		}
		return worker_utils.CreateDefaultErrorResponse(step, err)
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	stepValue := &api.HTTPResponse{
		Status:   resp.StatusCode,
		Response: buf.String(),
	}
	logger.Info("Made HTTP call")

	if h.ExpectedStatus != "" && h.ExpectedStatus != strconv.Itoa(resp.StatusCode) {
		logger.Infof("Expected status not matching for HTTP request: %s, result: %s", utils.StructToString(h), resp.Status)
		return &structs.StepResponse{
			Name:   step.Name,
			Status: structs.STEP_ERROR,
			Response: &structs.StepError{
				Expected: h.ExpectedStatus,
				Actual:   strconv.Itoa(resp.StatusCode),
				Error:    "Status code not matching",
			},
			Id: step.Id,
		}
	}

	if err = lib_utils.CompareStrings(&stepValue.Response, &h.ExpectedResponse); err != nil {
		logger.Infof("Expected responses not matching for HTTP request: %s, result: %s", utils.StructToString(h), stepValue.Response)
		return &structs.StepResponse{
			Name:   step.Name,
			Status: structs.STEP_ERROR,
			Response: &structs.StepError{
				Expected: h.ExpectedResponse,
				Actual:   stepValue.Response,
				Error:    "Response not matching",
			},
			Id: step.Id,
		}
	}

	return &structs.StepResponse{
		Name:     step.Name,
		Status:   structs.STEP_SUCCESS,
		Response: stepValue,
		Id:       step.Id,
	}
}
