package users

import (
	"context"
	"database/sql"
	"time"

	"github.com/DeepAung/gradient/website-server/modules/types"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound   = fiber.NewError(fiber.StatusBadRequest, "user not found")
	ErrUsernameUnique = fiber.NewError(fiber.StatusBadRequest, "username already exist")
	ErrEmailUnique    = fiber.NewError(fiber.StatusBadRequest, "email already exist")
)

type usersRepoImpl struct {
	db      *sqlx.DB
	timeout time.Duration
}

func NewUsersRepo(db *sqlx.DB, timeout time.Duration) types.UsersRepo {
	return &usersRepoImpl{
		db:      db,
		timeout: timeout,
	}
}

func (r *usersRepoImpl) CreateUser(username, email, hashedPassword string) (types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var user types.User
	err := r.db.GetContext(ctx, &user,
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

func (r *usersRepoImpl) FindOneUserById(id int) (types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var user types.User
	err := r.db.GetContext(ctx,
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

func (r *usersRepoImpl) FindOneUserWithPasswordByEmail(
	email string,
) (types.UserWithPassword, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var user types.UserWithPassword
	err := r.db.GetContext(ctx,
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

func (r *usersRepoImpl) FindOneUserPasswordById(id int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var hashedPassword string
	err := r.db.GetContext(ctx, &hashedPassword, `SELECT password FROM users WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return "", ErrUserNotFound
	}
	return hashedPassword, err
}

func (r *usersRepoImpl) UpdateUser(id int, req types.UpdateUserReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, `
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

func (r *usersRepoImpl) UpdateUserPassword(id int, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.db.ExecContext(
		ctx,
		`UPDATE users SET password = $1 WHERE id = $2;`,
		hashedPassword,
		id,
	)
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

func (r *usersRepoImpl) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1;`, id)
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
