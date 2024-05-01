package main

import (
	"encoding/json"
	"errors"
	"net/http"
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

	server.HandleFunc("/", MapHandlers(api.Some))
	server.HandleFunc("/create", MapHandlers(api.CreateUser))
	server.HandleFunc("/user", MapHandlers(api.GetUserById))

	http.ListenAndServe(api.address, server)
}

func (api *Api) Some(w http.ResponseWriter, r *http.Request) error {
	toJson(w, http.StatusOK, "it's all good, man")
	return nil
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