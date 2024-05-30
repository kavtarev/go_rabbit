package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Api struct {
	storage Storage
	address string
	channel *amqp.Channel
}

func NewApi(storage Storage, address string, channel *amqp.Channel) *Api {
	return &Api{
		storage: storage,
		address: address,
		channel: channel,
	}
}

func MapHandlers(f func(res http.ResponseWriter, req *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			responseAsJson(w, http.StatusBadRequest, err.Error())
		}
	}
}

func (api *Api) Run() {
	server := http.NewServeMux()

	server.HandleFunc("/", JWTAccess(MapHandlers(api.Some), api.storage))

	server.HandleFunc("/register", MapHandlers(api.Register))
	server.HandleFunc("/login", MapHandlers(api.Login))
	server.HandleFunc("/logout", MapHandlers(api.Logout))

	server.HandleFunc("/me", JWTAccess(MapHandlers(api.Me), api.storage))
	server.HandleFunc("/user", JWTAccess(MapHandlers(api.FindUserById), api.storage))
	server.HandleFunc("/users", JWTAccess(MapHandlers(api.ListUsers), api.storage))

	server.HandleFunc("/send-rabbit", MapHandlers(api.SendToRabbit))

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
	return responseAsJson(w, http.StatusOK, LoginResponse{Token: token})
}

func (api *Api) Logout(w http.ResponseWriter, req *http.Request) error {
	w.Header().Add("x-api-header", "token=;Max-Age=-;HttpOnly")
	return nil
}

func (api *Api) Me(w http.ResponseWriter, req *http.Request) error {
	if req.Method != http.MethodGet {
		return errors.New("only GET method allowed")
	}

	user := req.Context().Value(key).(*User)

	responseAsJson(w, http.StatusOK, MeResponse{Id: user.Id, Name: user.Name, Email: user.Email, Surname: user.Surname})
	return nil
}

func (api *Api) FindUserById(w http.ResponseWriter, req *http.Request) error {
	if req.Method != http.MethodGet {
		return errors.New("only GET method allowed")
	}

	id := req.URL.Query()["id"][0]
	if id == "" {
		return errors.New("no id provided")
	}

	user, err := api.storage.FindUserById(id)
	if err != nil {
		return err
	}

	// TODO validate correctly
	if user.Id == "" {
		return errors.New("not found")
	}

	responseAsJson(w, http.StatusOK, UserResponse{Id: user.Id, Name: user.Name, Email: user.Email, Surname: user.Surname})

	return nil
}

func (api *Api) ListUsers(w http.ResponseWriter, req *http.Request) error {
	if req.Method != http.MethodGet {
		return errors.New("only GET method allowed")
	}

	q := QueryParamsParser{values: req.URL.Query()}

	isCorrect := q.CheckCorrectness()
	if !isCorrect {
		return q.Errors[0]
	}

	users, err := api.storage.ListUsers(q)
	if err != nil {
		return err
	}

	result := make([]UserResponse, len(users))
	for i := 0; i < len(users); i++ {
		result[i] = UserResponse{Id: users[i].Id, Name: users[i].Name, Email: users[i].Email, Surname: users[i].Surname}
	}

	responseAsJson(w, http.StatusOK, result)

	return nil
}

func (api *Api) SendToRabbit(w http.ResponseWriter, req *http.Request) error {
	fmt.Println("in send-rabbit")
	err := api.channel.Publish(
		"",
		default_queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(`{"go_message": "some"}`)},
	)

	if err != nil {
		fmt.Println("error publishing to channel/queue", err)
	}
	return nil
}
