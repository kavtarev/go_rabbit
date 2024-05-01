package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetSome() string
	CreateUser(name, surname, email string) (*User, error)
	DeleteUser() error
	FindUser() *User
	UpdateUser() error
	Init()
}

type PostgresStorage struct {
	db *sql.DB
}

func (s *PostgresStorage) Init() {
	query := `
		create table if not exists users (
			id int generated always as identity
			, name text
			, surname text
			, email text
			, balance float
		)
	`;

	_, err := s.db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func NewPostgresStore() *PostgresStorage {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5433 user=postgres password=postgres dbname=golang sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("no ping", err)
	}

	return &PostgresStorage{
		db: db,
	}
}

func (s *PostgresStorage) GetSome() string {
	return "get-some"
}
func (s *PostgresStorage) CreateUser(name, surname, email string) (*User, error) {
	data, err := s.db.Query(`
		insert into users (
			name
			, surname
			, email
		) values (
			$1,$2,$3
		) returning name, surname, email, id, balance
	`,name,surname,email)
	if err != nil {
		fmt.Println("before scan", err)
	}
	defer data.Close()

	user := new(User)
	for data.Next() {
		err := data.Scan(
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.Id,
			&user.Balance,
		)
		if err != nil {
			fmt.Println("in scan", err)
		}
	}

	return user, nil
}
func (s *PostgresStorage) DeleteUser() error {
	return nil
}
func (s *PostgresStorage) FindUser() *User {
	return NewUser("","","")
}
func (s *PostgresStorage) UpdateUser() error {
	return nil
}

// type MockStorage struct {

// }

// func (s *MockStorage) GetSome() string {
// 	return "get-some"
// }
// func (s *MockStorage) CreateUser() *User {
// 	return NewUser()
// }
// func (s *MockStorage) DeleteUser() error {
// 	return nil
// }
// func (s *MockStorage) FindUser() *User {
// 	return NewUser()
// }
// func (s *MockStorage) UpdateUser() error {
// 	return nil
// }
