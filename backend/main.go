package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type pasta struct {
	Id   int
	Name string
	Text string
	Tags []string
}

func get_pastas(w http.ResponseWriter, q *http.Request) {
	if q.Method != "GET" {
		w.WriteHeader(400)
		w.Write([]byte("Wrong request"))
	}
	pastas, err := get_records()
	if err != nil {
		w.WriteHeader(501)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(200)
	mrsh, err := json.Marshal(pastas)
	if err != nil {
		w.WriteHeader(501)
		w.Write([]byte(err.Error()))
	}
	w.Write(mrsh)
}

func fetch_pasta(w http.ResponseWriter, q *http.Request) {
	if q.Method != "GET" {
		w.WriteHeader(400)
		w.Write([]byte("Wrong request"))
		return
	}
	urlPathElements := strings.Split(q.URL.Path, "/")
	id, err := strconv.Atoi(urlPathElements[2])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Id must be integer bro"))
		return
	}
	p, err := get_pasta(id)
	if err != nil {
		if err.Error() == "pasta does not exist" {
			w.WriteHeader(404)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(501)
		w.Write([]byte(err.Error()))
		return
	}
	mrsh, err := json.Marshal(p)
	if err != nil {
		w.WriteHeader(501)
		w.Write([]byte(err.Error()))
	}
	w.Write(mrsh)
}

func main() {
	_ = create_db()
	err := populate_db()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	http.HandleFunc("/get_pasta_list", get_pastas)
	http.HandleFunc("/get_pasta/", fetch_pasta)
	s := &http.Server{
		Addr:           ":8000",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err = s.ListenAndServe()
}
