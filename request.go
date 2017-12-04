package fit

import (
	"net/http"
	"strconv"
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

func (r *Request) FormIntValue(key string) int {
	intval, err := strconv.Atoi(r.FormValue(key))
	if err != nil {
		intval = 0
	}
	return intval
}

func (r *Request) FormInt64Value(key string) int64 {
	intval, err := strconv.ParseInt(r.FormValue(key), 10, 64)
	if err != nil {
		intval = 0
	}
	return intval
}

