package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

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

	server.HandleFunc("/", JWTAccess(MapHandlers(api.Some), api.storage))
	server.HandleFunc("/ctx", MapHandlers(api.Some))

	server.HandleFunc("/register", MapHandlers(api.Register))
	server.HandleFunc("/login", MapHandlers(api.Login))
	server.HandleFunc("/logout", MapHandlers(api.Logout))

	http.ListenAndServe(api.address, server)
}

func (api *Api) Some(w http.ResponseWriter, r *http.Request) error {
	user := r.Context().Value(key).(*User)

	responseAsJson(w, http.StatusOK, user)
	return nil
}

func (api *Api) Register(w http.ResponseWriter, req *http.Request) error {
	if req.Method != http.MethodPost {
		return errors.New("only POST")
	}

	dto := new(RegisterDto)

	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()

	err := d.Decode(dto)
	if err != nil {
		return err
	}
	if dto.Email == "" || dto.Name == "" || dto.Surname == "" || dto.Password == "" {
		return errors.New("invalid registration data")
	}

	res, err := api.storage.CreateUser(*dto)
	if err != nil {
		return err
	}

	return responseAsJson(w, http.StatusCreated, RegisterResponse{Id: res.Id, Email: res.Email})
}

func (api *Api) Login(w http.ResponseWriter, req *http.Request) error {
	if req.Method != http.MethodPost {
		return errors.New("only POST")
	}

	dto := new(LoginDto)

	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()

	err := d.Decode(dto)
	if err != nil {
		return err
	}

	user, err := api.storage.FindByEmail(dto.Email)
	if err != nil {
		return err
	}

	if user.Password.String != dto.Password {
		return errors.New("password is incorrect")
	}

	token, err := createJWT(user.Id, 10000)
	if err != nil {
		return err
	}

	w.Header().Add("x-api-header", fmt.Sprintf("token=%v;Max-Age=90;HttpOnly", token))
	return responseAsJson(w, http.StatusOK, LoginResponse{ Token: token })
}

func (api *Api) Logout(w http.ResponseWriter, req *http.Request) error {
	w.Header().Add("x-api-header", "token=;Max-Age=-;HttpOnly")
	return nil
}
