package handlerwrapper

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func HandleFunc(f func(w http.ResponseWriter, r *http.Request)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		f(w, r)
	}
}
