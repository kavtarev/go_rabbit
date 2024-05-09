package main

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

const secret = "some_secret"

func createJWT(id string, ttl int) (string, error){
	claims := jwt.MapClaims{
		"sub":  id,
		"exp":  time.Now().Add(time.Second * time.Duration(ttl)).Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}


func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	
}