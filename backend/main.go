package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
	"net/http"
	"os"
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

type user struct {
	Credentials string
	Password    string
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

func generate_jwt(id string) (string, error) {
	expire := time.Now().Add(time.Hour * 8) //JWT żyje 8 godzin
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = expire.Unix()
	res, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	return res, err
}

func hash_pwd(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)
	var hashed_pwd string
	for _, b := range hash {
		hashed_pwd += fmt.Sprintf("%x", b)
	}
	return hashed_pwd
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

func validate_signature(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return "Unauthorized", nil
	}
	return []byte(os.Getenv("JWT_KEY")), nil
}

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
	err = s.ListenAndServe()
}
