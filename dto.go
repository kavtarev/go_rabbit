package main

type RegisterDto struct {
	Name string `json:"name"`
	Surname string `json:"surname"`
	Email string `json:"email"`
}

type LoginDto struct {
	Email string `json:"email"`
	Password string `json:"password"`
}
