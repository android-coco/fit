package fit

import (
	"net/http"
)

type Request struct {
	*http.Request
}

// Make new Request instance
func newRequest(req *http.Request) *Request {
	request := &Request{
		Request: req,
	}

	return request
}
