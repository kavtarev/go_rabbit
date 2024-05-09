package main

import (
	"fmt"
	"net/http"
	"strings"
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

func JWTAccess(fn http.HandlerFunc, storage Storage) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		cookieDict := splitCookies(r.Header.Get("x-api-header"))

		t, err := validateJWT(cookieDict["token"])
		if err != nil {
			responseAsJson(w, http.StatusBadRequest, "i see what you did there")
			return
		}
		if !t.Valid {
			responseAsJson(w, http.StatusBadRequest, "i see what you did there")
			return
		}

		sub, err := t.Claims.GetSubject()
		if err != nil {
			responseAsJson(w, http.StatusBadRequest, "i see what you did there")
			return
		}

		user, err := storage.FindUserById(sub)
		if err != nil {
			responseAsJson(w, http.StatusBadRequest, err.Error())
			return
		}
		if user.Id == "" {
			responseAsJson(w, http.StatusBadRequest, "i see what you did there punk")
			return
		}

		fn(w,r)
	}
}

func splitCookies(s string) map[string]string {
	arr := strings.Split(s, ";")
	dict := make(map[string]string)

	for i := 0; i < len(arr); i++ {
		pair := strings.Split(arr[i], "=")
		if len(pair) == 2 {
			dict[pair[0]] = pair[1]
		}
	}

	return dict
}