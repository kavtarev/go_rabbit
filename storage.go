package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

type Storage interface {
	GetSome() string
	CreateUser(dto RegisterDto) (*User, error)
	DeleteUser() error
	FindUser() *User
	UpdateUser() error
	Init()
}

type PostgresStorage struct {
	db *sql.DB
}

func (s *PostgresStorage) Init() {
	query := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`;

	_, err := s.db.Exec(query)
	if err != nil {
		panic(err)
	}

	query = `
		create table if not exists users (
			id uuid not null default uuid_generate_v4()
			, name text
			, surname text
			, email text
			, password text
		)
	`;

	_, err = s.db.Exec(query)
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

func (s *PostgresStorage) FindByEmail(email string) (*User, error) {
	query, err := s.db.Query(`select name, surname, email, id, password from users where email = $1`, email)
	if err != nil {
		fmt.Println("find by email error")
		return nil, err
	}
	defer query.Close()

	var user User
	for query.Next() {
		err := query.Scan(
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.Id,
			&user.Password,
		)
		if err != nil {
			fmt.Println("in scan", err)
		}
	}
	return &user, nil
}

func (s *PostgresStorage) CreateUser(dto RegisterDto) (*User, error) {
	existingUser, err := s.FindByEmail(dto.Email)
	if err != nil {
		return nil, err
	}

	// TODO validate correctly
	if existingUser.Email != "" {
		return nil, errors.New("email already taken")
	}

	data, err := s.db.Query(`
		insert into users (
			name
			, surname
			, email
			, password
		) values (
			$1,$2,$3,$4
		) returning name, surname, email, id, password
	`,dto.Name,dto.Surname,dto.Email,dto.Password)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	user := new(User)
	for data.Next() {
		err := data.Scan(
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.Id,
			&user.Password,
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
