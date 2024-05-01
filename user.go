package main

import "math/rand"

type User struct {
	Name string `json:"name"`
	Surname string `json:"surname"`
	Amount int64 `json:"amount"`
	Id int64 `json:"id"`
}


func NewUser() *User {
	return &User{
		Name: "name",
		Surname: "surname",
		Amount: int64(rand.Int()),
		Id: int64(rand.Int()),
	}
}