package fit

import (
	"net/http"
	"strconv"
	"time"
	"errors"
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

func (r *Request) FormTimeStruct(key string) (datetime time.Time, err error) {
	datetime, err = time.ParseInLocation("2006-01-02 15:04:05", r.FormValue("datetime"), time.Local)
	if err != nil  {
		return
	}
	if datetime.IsZero() {
		err = errors.New("date is zero!")
		return
	}
	return datetime, nil
}
