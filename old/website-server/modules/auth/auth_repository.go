package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/DeepAung/gradient/website-server/modules/types"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

var ErrTokenNotFound = fiber.NewError(fiber.StatusBadRequest, "token not found")

type authRepoImpl struct {
	db      *sqlx.DB
	timeout time.Duration
}

func NewAuthRepo(db *sqlx.DB, timeout time.Duration) types.AuthRepo {
	return &authRepoImpl{
		db:      db,
		timeout: timeout,
	}
}

func (r *authRepoImpl) CreateToken(accessToken, refreshToken string) (types.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var token types.Token
	err := r.db.GetContext(ctx, &token,
		`INSERT INTO tokens (access_token, refresh_token)
			VALUES ($1, $2)
		RETURNING id, access_token, refresh_token;`,
		accessToken, refreshToken)

	if err == sql.ErrNoRows {
		return types.Token{}, ErrTokenNotFound
	}

	return token, err
}

func (r *authRepoImpl) DeleteToken(tokenId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, `DELETE FROM tokens WHERE tokens.id = $1;`, tokenId)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrTokenNotFound
	}
	return nil
}

func (r *authRepoImpl) HasToken(id int, refreshToken string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var tmp int
	err := r.db.GetContext(ctx, &tmp,
		`SELECT 1 FROM tokens WHERE tokens.id = $1 AND tokens.refresh_token = $2;`,
		id, refreshToken)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return tmp == 1, err
}

func (r *authRepoImpl) UpdateTokens(id int, newAccessToken, newRefreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx,
		`UPDATE tokens
		SET access_token = $1, refresh_token = $2
		WHERE tokens.id = $3;`,
		newAccessToken, newRefreshToken, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrTokenNotFound
	}
	return nil
}
