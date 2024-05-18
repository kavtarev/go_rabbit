package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserMetaKey string

var key UserMetaKey = "user_meta"

func createJWT(id string, ttl int) (string, error) {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("no secret key found")
	}
	claims := jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Second * time.Duration(ttl)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("no secret key found")
	}
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

func JWTAccess(fn http.HandlerFunc, storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		ctx := context.WithValue(r.Context(), key, user)

		fn(w, r.WithContext(ctx))
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
