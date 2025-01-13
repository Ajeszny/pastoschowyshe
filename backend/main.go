package main

import (
	"net/http"
	"time"
)

func get_pastas(r http.ResponseWriter, q *http.Request) {
	r.WriteHeader(200)
	r.Write([]byte("[{\"id\": 1, \"name\":\"Вован\", \"tags\": []}]"))
}

func main() {
	http.HandleFunc("/get_pasta_list", get_pastas)
	s := &http.Server{
		Addr:           ":8000",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	print(err)
}
