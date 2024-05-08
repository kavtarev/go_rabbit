package main

import (
	"errors"
	"encoding/json"
	"fmt"
	"net/http"
)

const secret = "some_secret"
type Api struct {
	storage Storage
	address string
}

func NewApi(storage Storage,address string) *Api {
	return &Api{
		storage: storage,
		address: address,
	}
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
			responseAsJson(w, http.StatusBadRequest, "goooooooovna")
			return
		}

		fn(w,r)
	}
}


func MapHandlers (f func(res http.ResponseWriter, req *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			responseAsJson(w, http.StatusBadRequest, err.Error())
		}
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
	responseAsJson(w, http.StatusOK, "it's all good, man")
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