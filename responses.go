package main

type RegisterResponse struct {
	Id string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}