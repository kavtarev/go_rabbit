package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Api struct {
	storage Storage
	address string
}

func toJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func JWTAccess(fn http.HandlerFunc) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("in jwt")
		token := r.Header.Get("jwt_token")
		t, err := validateJWT(token)
		fmt.Println("Header", t.Header)
		fmt.Println("Claims", t.Claims)
		fmt.Println("Signature", t.Signature)
		fmt.Println("Valid", t.Valid)
		if err != nil {
			toJson(w, http.StatusBadRequest, "goooooooovna")
			return
		}

		fn(w,r)
	}
}
const secret = "some_secret"

func createJWT(id string) (string, error){
	claims := jwt.MapClaims{
		"sub":  id,
		"exp":  time.Now().Add(time.Second * 10000).Unix(),
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


func MapHandlers (f func(res http.ResponseWriter, req *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			toJson(w, http.StatusBadRequest, err.Error())
		}
	}
}

func NewApi(storage Storage,address string) *Api {
	return &Api{
		storage: storage,
		address: address,
	}
}

func (api *Api) Run() {
	api.storage.Init()
	server := http.NewServeMux()

	server.HandleFunc("/", JWTAccess(MapHandlers(api.Some)))
	server.HandleFunc("/create", MapHandlers(api.CreateUser))
	server.HandleFunc("/user", MapHandlers(api.GetUserById))
	server.HandleFunc("/token", MapHandlers(api.CreateToken))

	http.ListenAndServe(api.address, server)
}

func (api *Api) Some(w http.ResponseWriter, r *http.Request) error {
	toJson(w, http.StatusOK, "it's all good, man")
	return nil
}

func (api *Api) CreateToken(w http.ResponseWriter, req *http.Request) error {
	token, err := createJWT("some")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token)

	return nil;
}

func (api *Api) CreateUser(w http.ResponseWriter, req *http.Request) error {
	if req.Method != http.MethodPost {
		return errors.New("only POST")
	}
	user := new(User)
	json.NewDecoder(req.Body).Decode(user)

	user, err := api.storage.CreateUser(user.Name, user.Surname, user.Email)
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(user)
	return nil
}

func (api *Api) GetUserById(res http.ResponseWriter, req *http.Request) error {
	// json.NewEncoder(res).Encode(NewUser())
	return errors.New("poshel v hui")
}