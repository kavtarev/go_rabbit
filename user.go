package main

import (
	"database/sql"
)
type User struct {
	Name string `json:"name"`
	Surname string `json:"surname"`
	Email string `json:"email"`
	Balance sql.NullFloat64 `json:"balance"`
	Id int64 `json:"id"`
}


func NewUser(name, surname, email string) *User {
	return &User{
		Name: name,
		Surname: surname,
		Email: email,
	}
}