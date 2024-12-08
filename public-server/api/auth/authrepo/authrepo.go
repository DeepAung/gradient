package authrepo

import (
	"database/sql"
	"errors"

	"github.com/DeepAung/gradient/public-server/api/auth"
	"github.com/jmoiron/sqlx"
)

var ErrTokenNotFound = errors.New("token not found")

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) CreateToken(accessToken, refreshToken string) (auth.Token, error) {
	var token auth.Token
	err := r.db.Get(&token,
		`INSERT INTO tokens (access_token, refresh_token)
			VALUES ($1, $2)
		RETURNING id, access_token, refresh_token;`,
		accessToken, refreshToken)

	if err == sql.ErrNoRows {
		return auth.Token{}, ErrTokenNotFound
	}

	return token, err
}

func (r *AuthRepo) DeleteToken(tokenId int) error {
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

func (r *AuthRepo) HasToken(id int, refreshToken string) (bool, error) {
	var tmp int
	err := r.db.Get(&tmp,
		`SELECT 1 FROM tokens WHERE tokens.id = $1 AND tokens.refresh_token = $2;`,
		id, refreshToken)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return tmp == 1, err
}

func (r *AuthRepo) UpdateTokens(id int, newAccessToken, newRefreshToken string) error {
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
