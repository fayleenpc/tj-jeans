package auth

import (
	"fmt"
	"net/http"
)

func WithLogger(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logNote := fmt.Sprintf("Host : %v, Method : %v, Remote Address : %v, Request URI : %v, Header : \n", r.Host, r.Method, r.RemoteAddr, r.RequestURI)

		for k, v := range r.Header {
			logNote += fmt.Sprintf("%v : %v\n", k, v)
		}

		// log.Println(logNote)

		handlerFunc(w, r)
	}
}
