package usersrepo

import (
	"database/sql"
	"errors"

	"github.com/DeepAung/gradient/public-server/api/users"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUsernameUnique = errors.New("username already exist")
	ErrEmailUnique    = errors.New("email already exist")
)

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
		RETURNING id, username, email, picture_url, is_admin;`,
		username, email, hashedPassword,
	)
	if err == nil {
		return user, nil
	}

	switch err.Error() {
	case `pq: duplicate key value violates unique constraint "users_username_key"`:
		return users.User{}, ErrUsernameUnique
	case `pq: duplicate key value violates unique constraint "users_email_key"`:
		return users.User{}, ErrEmailUnique
	default:
		return users.User{}, err
	}
}

func (r *UsersRepo) FindOneUserWithPasswordByEmail(email string) (users.UserWithPassword, error) {
	var user users.UserWithPassword
	err := r.db.Get(
		&user,
		`SELECT id, username, email, password, picture_url, is_admin
		FROM users WHERE users.email = $1`,
		email,
	)
	if err == sql.ErrNoRows {
		return users.UserWithPassword{}, ErrUserNotFound
	}
	return user, err
}
