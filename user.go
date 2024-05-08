package main

import (
	"database/sql"
)
type User struct {
	Name string `json:"name"`
	Surname string `json:"surname"`
	Email string `json:"email"`
	Password sql.NullString `json:"password"`
	Id string `json:"id"`
}


func NewUser(name, surname, email string) *User {
	return &User{
		Name: name,
		Surname: surname,
		Email: email,
	}
}