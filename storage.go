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

	FindUserById(id string) (*User, error)
	FindByEmail(email string) (*User, error)

	ListUsers(q QueryParamsParser) ([]*User, error)
}

type PostgresStorage struct {
	db *sql.DB
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
			return nil, err
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
			return nil, err
		}
	}

	return user, nil
}

func (s *PostgresStorage) FindUserById(id string) (*User, error) {
	user := User{}

	res, err := s.db.Query(`select id, name, surname, email, password from users where id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		err := res.Scan(
			&user.Id,
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.Password,
		)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *PostgresStorage) ListUsers(q QueryParamsParser) ([]*User, error) {
	args := []any{q.Limit}

	offsetString := ""
	if q.Page != 0 {
		args = append(args, q.Page * q.Limit)
		offsetString = fmt.Sprintf("offset $%d", len(args))
	}

	qString := ""
	if q.Q != "" {
		args = append(args, q.Q)
		qString = fmt.Sprintf("where name iLike '%%' || $%d || '%%' or email iLike '%%' || $%d || '%%'" ,len(args),len(args))
	}

	query := fmt.Sprintf(
		`
			select id, name, surname, email
			from users
			%v
			%v
			limit $1
		`, qString,  offsetString,
	)

	res, err := s.db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	users := make([]*User, 0)

	for res.Next() {
		user := new(User)
		err := res.Scan(
			&user.Id,
			&user.Name,
			&user.Surname,
			&user.Email,
		)
		
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

