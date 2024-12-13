package auth

import (
	"database/sql"

	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

var ErrTokenNotFound = fiber.NewError(fiber.StatusBadRequest, "token not found")

type authRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) types.AuthRepo {
	return &authRepo{
		db: db,
	}
}

func (r *authRepo) CreateToken(accessToken, refreshToken string) (types.Token, error) {
	var token types.Token
	err := r.db.Get(&token,
		`INSERT INTO tokens (access_token, refresh_token)
			VALUES ($1, $2)
		RETURNING id, access_token, refresh_token;`,
		accessToken, refreshToken)

	if err == sql.ErrNoRows {
		return types.Token{}, ErrTokenNotFound
	}

	return token, err
}

func (r *authRepo) DeleteToken(tokenId int) error {
	result, err := r.db.Exec(`DELETE FROM tokens WHERE tokens.id = $1;`, tokenId)
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

func (r *authRepo) HasToken(id int, refreshToken string) (bool, error) {
	var tmp int
	err := r.db.Get(&tmp,
		`SELECT 1 FROM tokens WHERE tokens.id = $1 AND tokens.refresh_token = $2;`,
		id, refreshToken)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return tmp == 1, err
}

func (r *authRepo) UpdateTokens(id int, newAccessToken, newRefreshToken string) error {
	result, err := r.db.Exec(
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
