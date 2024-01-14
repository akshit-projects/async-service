package api_validator

import (
	"errors"
	"net/url"
	"strings"

	"github.com/akshitbansal-1/async-testing/lib/structs"
	"github.com/akshitbansal-1/async-testing/lib/structs/api"
	"github.com/akshitbansal-1/async-testing/lib/utils"
)

var validMethods = [4]string{"GET", "POST", "PUT", "DELETE"}

// HTTP validation block
func ValidateHttpStep(step *structs.Step) error {
	meta := step.Meta
	var httpReq api.HTTPRequest
	if err := utils.ParseInterface[api.HTTPRequest](meta, &httpReq); err != nil {
		return errors.New("Unable to get http request data")
	}

	httpReq.Method = strings.ToUpper(httpReq.Method)
	if err := validateHTTPMethod(httpReq.Method); err != nil {
		return err
	}

	// Check for specific validation conditions
	if strings.ToUpper(httpReq.Method) == "GET" && httpReq.Body != nil {
		return errors.New("Body can't go with GET method")
	} else if strings.ToUpper(httpReq.Method) != "GET" && httpReq.Body == nil {
		return errors.New("Body is required for " + httpReq.Method)
	}

	_, err := url.ParseRequestURI(httpReq.Url)
	if err != nil {
		return errors.New("Invalid request URL passed")
	}

	step.Value = &httpReq
	return nil
}

func validateHTTPMethod(method string) error {
	for _, m := range &validMethods {
		if method == m {
			return nil
		}
	}

	return errors.New("Invalid http method provided")
}
