package main

import (
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func authorize(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] == nil {
			w.WriteHeader(403)
			w.Write([]byte("Unauthorized"))
			return
		}
		token, err := jwt.Parse(strings.Split(r.Header["Authorization"][0], " ")[1], validate_signature)
		if err != nil {
			w.WriteHeader(403)
			w.Write([]byte("Unauthorized"))
			return
		}
		if !token.Valid {
			w.WriteHeader(403)
			w.Write([]byte("Unauthorized"))
			return
		}
		endpointHandler(w, r)
	})
}
