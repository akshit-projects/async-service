package common_structs

type HttpError struct {
	Msg string `json:"msg"`
}

type APIFilter struct {
	Limit   int64
	Filters map[string]interface{}
	Skip    int64
}
