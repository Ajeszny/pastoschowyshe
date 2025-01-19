package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

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

func validate_signature(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return "Unauthorized", nil
	}
	return []byte(os.Getenv("JWT_KEY")), nil
}
