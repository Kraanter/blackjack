package routes

import (
	"net/http"
)

type ApiRoute struct {
	Pattern    string
	Handler    func(http.ResponseWriter, *http.Request)
	MiddleWare []func(http.ResponseWriter, *http.Request) bool
}

func (r *ApiRoute) TotalHandler() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		for _, v := range r.MiddleWare {
			if v(res, req) {
				return
			}
		}

		r.Handler(res, req)
	}
}
