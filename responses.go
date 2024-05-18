package main

type RegisterResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type MeResponse struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type UserResponse struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
