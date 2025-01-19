package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

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

func login(w http.ResponseWriter, q *http.Request) {
	if q.Method != "POST" {
		w.WriteHeader(400)
		w.Write([]byte("Wrong request"))
		return
	}
	decoder := json.NewDecoder(q.Body)
	var info user
	err := decoder.Decode(&info)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Wrong format"))
		return
	}
	check, err := check_creds(info.Credentials, info.Password)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("It all went south bro... " + err.Error()))
		return
	}
	if check {
		jwt, err := generate_jwt(info.Credentials)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("It all went south bro... " + err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("{\"token\": \"%s\"}", jwt)))
	} else {
		w.WriteHeader(403)
		w.Write([]byte("Forbidden"))
		return
	}
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

func add_pasta(w http.ResponseWriter, q *http.Request) {
	if q.Method != "POST" {
		w.WriteHeader(400)
		w.Write([]byte("Wrong request"))
		return
	}
	decoder := json.NewDecoder(q.Body)
	var data pasta
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Wrong request"))
		return
	}
	err = add_new_record(data)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Database error" + err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Success"))
}
