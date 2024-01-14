package api

type HTTPRequest struct {
	Url              string            `json:"url"`
	Method           string            `json:"method"`
	Body             interface{}       `json:"body"`
	Headers          map[string]string `json:"headers"`
	ExpectedStatus   string            `json:"expectedStatus" bson:"expectedStatus"`
	ExpectedResponse string            `json:"expectedResponse" bson:"expectedResponse"`
}

type HTTPResponse struct {
	Status   int    `json:"status"`
	Response string `json:"response"`
}
