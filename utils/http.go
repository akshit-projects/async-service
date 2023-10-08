package utils

import "net/http"

var client = http.Client{}

func CallHTTP(req *http.Request) (*http.Response, error) {
	return client.Do(req)
}
