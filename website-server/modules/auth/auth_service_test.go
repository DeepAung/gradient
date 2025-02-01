package auth

import (
	"testing"

	"github.com/DeepAung/gradient/website-server/config"
	"github.com/DeepAung/gradient/website-server/database"
	"github.com/DeepAung/gradient/website-server/modules/types"
	"github.com/DeepAung/gradient/website-server/modules/users"
	"github.com/DeepAung/gradient/website-server/pkg/asserts"
	"github.com/jmoiron/sqlx"
)

var (
	migrateSourceName = "../../migrations/migrate.sql"
	seedSourceName    = "../../migrations/seed.sql"
	cfg               *config.Config
	db                *sqlx.DB
	myUsersRepo       types.UsersRepo
	myAuthRepo        types.AuthRepo
	myAuthSvc         types.AuthSvc
)

func init() {
	cfg = config.NewConfig("../../.env.dev")
	db = database.InitDB(cfg.App.DbUrl)
	database.RunSQL(db, migrateSourceName)
	database.RunSQL(db, seedSourceName)
	myUsersRepo = users.NewUsersRepo(db, cfg.App.Timeout)
	myAuthRepo = NewAuthRepo(db, cfg.App.Timeout)
	myAuthSvc = NewAuthSvc(myAuthRepo, myUsersRepo, cfg)
}

func TestSignUp(t *testing.T) {
	t.Run("unique username", func(t *testing.T) {
		passport, err := myAuthSvc.SignUp("DeepAung", "newuser@gmail.com", "password")
		asserts.EqualError(t, err, users.ErrUsernameUnique)
		asserts.Equal(t, "passport", passport, types.Passport{})
	})

	t.Run("unique email", func(t *testing.T) {
		passport, err := myAuthSvc.SignUp("newuser", "i.deepaung@gmail.com", "password")
		asserts.EqualError(t, err, users.ErrEmailUnique)
		asserts.Equal(t, "passport", passport, types.Passport{})
	})

	t.Run("normal sign up", func(t *testing.T) {
		passport, err := myAuthSvc.SignUp("newuser", "newuser@gmail.com", "password")
		asserts.EqualError(t, err, nil)
		assertNewUserPassport(t, passport)
	})
}

func TestSignIn(t *testing.T) {
	t.Run("non-exist email", func(t *testing.T) {
		passport, err := myAuthSvc.SignIn("nonexist@gmail.com", "password")
		asserts.EqualError(t, err, users.ErrUserNotFound)
		asserts.Equal(t, "passport", passport, types.Passport{})
	})

	t.Run("empty password", func(t *testing.T) {
		passport, err := myAuthSvc.SignIn("newuser@gmail.com", "")
		asserts.EqualError(t, err, ErrInvalidEmailOrPassword)
		asserts.Equal(t, "passport", passport, types.Passport{})
	})

	t.Run("invalid password", func(t *testing.T) {
		passport, err := myAuthSvc.SignIn("newuser@gmail.com", "invalid")
		asserts.EqualError(t, err, ErrInvalidEmailOrPassword)
		asserts.Equal(t, "passport", passport, types.Passport{})
	})

	t.Run("normal sign in", func(t *testing.T) {
		passport, err := myAuthSvc.SignIn("newuser@gmail.com", "password")
		asserts.EqualError(t, err, nil)
		assertNewUserPassport(t, passport)
	})
}

func TestSignOut(t *testing.T) {
	t.Run("invalid token id", func(t *testing.T) {
		err := myAuthSvc.SignOut(1000)
		asserts.EqualError(t, err, ErrTokenNotFound)
	})

	t.Run("normal sign out", func(t *testing.T) {
		passport, err := myAuthSvc.SignIn("newuser@gmail.com", "password")
		asserts.EqualError(t, err, nil)
		err = myAuthSvc.SignOut(passport.Token.Id)
		asserts.EqualError(t, err, nil)
	})
}

func TestUpdateTokens(t *testing.T) {
	passport, err := myAuthSvc.SignIn("newuser@gmail.com", "password")
	asserts.EqualError(t, err, nil)

	t.Run("invalid token id", func(t *testing.T) {
		token, _, err := myAuthSvc.UpdateTokens(1000, passport.Token.RefreshToken)
		asserts.EqualError(t, err, ErrInvalidRefreshTokenOrTokenNotFound)
		asserts.Equal(t, "token", token, types.Token{})
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		token, _, err := myAuthSvc.UpdateTokens(passport.Token.Id, "invalid refresh token")
		asserts.EqualError(t, err, ErrInvalidRefreshTokenOrTokenNotFound)
		asserts.Equal(t, "token", token, types.Token{})
	})

	t.Run("empty refresh token", func(t *testing.T) {
		token, _, err := myAuthSvc.UpdateTokens(passport.Token.Id, "")
		asserts.EqualError(t, err, ErrInvalidRefreshTokenOrTokenNotFound)
		asserts.Equal(t, "token", token, types.Token{})
	})

	t.Run("normal update tokens", func(t *testing.T) {
		token, _, err := myAuthSvc.UpdateTokens(passport.Token.Id, passport.Token.RefreshToken)
		asserts.EqualError(t, err, nil)

		asserts.Equal(t, "token id", token.Id, passport.Token.Id)

		// Equal because jwt use time precision of one second
		asserts.Equal(t, "access token", token.AccessToken, passport.Token.AccessToken)
		asserts.Equal(t, "refresh token", token.RefreshToken, passport.Token.RefreshToken)
	})

	err = myAuthSvc.SignOut(passport.Token.Id)
	asserts.EqualError(t, err, nil)
}

func assertNewUserPassport(t *testing.T, passport types.Passport) {
	asserts.Equal(t, "", passport.User.Username, "newuser")
	asserts.Equal(t, "", passport.User.Email, "newuser@gmail.com")
	asserts.Equal(t, "", passport.User.PictureUrl, "")
	asserts.Equal(t, "", passport.User.IsAdmin, false)
	asserts.NotEqual(t, "", passport.Token.AccessToken, "")
	asserts.NotEqual(t, "", passport.Token.RefreshToken, "")
}
