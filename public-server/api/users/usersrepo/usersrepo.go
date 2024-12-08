package usersrepo

import (
	"database/sql"
	"errors"

	"github.com/DeepAung/gradient/public-server/api/users"
	"github.com/jmoiron/sqlx"
)

var ErrUserNotFound = errors.New("user not found")

type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) CreateUser(username, email, hashedPassword string) (users.User, error) {
	var user users.User
	err := r.db.Get(&user,
		`INSERT INTO users (username, email, password)
			VALUES ($1, $2, $3)
		RETURNING *;`,
		username, email, hashedPassword,
	)
	if err == sql.ErrNoRows {
		return users.User{}, ErrUserNotFound
	}
	return user, err
}

func (r *UsersRepo) FindOneUserWithPasswordByEmail(email string) (users.UserWithPassword, error) {
	var user users.UserWithPassword
	err := r.db.Get(&user,
		`SELECT * from users WHERE users.email = $1`, email)
	if err == sql.ErrNoRows {
		return users.UserWithPassword{}, ErrUserNotFound
	}
	return user, err
}
