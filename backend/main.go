package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

type pasta struct {
	Id   int
	Name string
	Text string
	Tags []string
}

type user struct {
	Credentials string
	Password    string
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
	http.HandleFunc("/login", login)
	http.HandleFunc("/add_pasta", authorize(add_pasta))
	s := &http.Server{
		Addr:           ":8000",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err = s.ListenAndServeTLS("credentials/cert.crt", "credentials/k.key")
}
