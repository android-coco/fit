package fit

import (
	"net/http"
)

var (
	gHandler *handler
)

type handler struct {
}

func (hdl *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}

func Handler() *handler {
	if gHandler == nil {
		gHandler = &handler{}
	}

	return gHandler
}
