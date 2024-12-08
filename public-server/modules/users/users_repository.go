package users

import (
	"database/sql"
	"errors"

	"github.com/DeepAung/gradient/public-server/modules/types"
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

func NewUsersRepo(db *sqlx.DB) types.UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) CreateUser(username, email, hashedPassword string) (types.User, error) {
	var user types.User
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
		return types.User{}, ErrUsernameUnique
	case `pq: duplicate key value violates unique constraint "users_email_key"`:
		return types.User{}, ErrEmailUnique
	default:
		return types.User{}, err
	}
}

func (r *UsersRepo) FindOneUserById(id int) (types.User, error) {
	var user types.User
	err := r.db.Get(
		&user,
		`SELECT id, username, email, picture_url, is_admin
		FROM users WHERE id = $1`,
		id,
	)
	if err == sql.ErrNoRows {
		return types.User{}, ErrUserNotFound
	}
	return user, err
}

func (r *UsersRepo) FindOneUserWithPasswordByEmail(email string) (types.UserWithPassword, error) {
	var user types.UserWithPassword
	err := r.db.Get(
		&user,
		`SELECT id, username, email, password, picture_url, is_admin
		FROM users WHERE email = $1`,
		email,
	)
	if err == sql.ErrNoRows {
		return types.UserWithPassword{}, ErrUserNotFound
	}
	return user, err
}

func (r *UsersRepo) FindOneUserPasswordById(id int) (string, error) {
	var hashedPassword string
	err := r.db.Get(&hashedPassword, `SELECT password FROM users WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return "", ErrUserNotFound
	}
	return hashedPassword, err
}

func (r *UsersRepo) UpdateUser(id int, req types.UpdateUserReq) error {
	result, err := r.db.Exec(`
		UPDATE users
		SET username = $1,
		picture_url = COALESCE($2, picture_url),
		password = COALESCE($3, password)
		WHERE id = $4;`,
		req.Username, req.PictureUrl, req.NewPassword, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UsersRepo) UpdateUserPassword(id int, hashedPassword string) error {
	result, err := r.db.Exec(`UPDATE users SET password = $1 WHERE id = $2;`, hashedPassword, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UsersRepo) DeleteUser(id int) error {
	result, err := r.db.Exec(`DELETE FROM users WHERE id = $1;`, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}
