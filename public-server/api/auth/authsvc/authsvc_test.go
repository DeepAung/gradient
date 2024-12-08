package authsvc

import (
	"testing"

	"github.com/DeepAung/gradient/public-server/api/auth"
	"github.com/DeepAung/gradient/public-server/api/auth/authrepo"
	"github.com/DeepAung/gradient/public-server/api/users/usersrepo"
	"github.com/DeepAung/gradient/public-server/pkg/asserts"
	"github.com/DeepAung/gradient/public-server/pkg/config"
	"github.com/DeepAung/gradient/public-server/pkg/database"
	"github.com/jmoiron/sqlx"
)

var (
	migrateSourceName = "../../../migrations/migrate.sql"
	seedSourceName    = "../../../migrations/seed.sql"
	cfg               *config.Config
	db                *sqlx.DB
	usersRepo         *usersrepo.UsersRepo
	authRepo          *authrepo.AuthRepo
	authSvc           *AuthSvc
)

func init() {
	cfg = config.NewConfig("../../../.env.dev")
	db = database.InitDB(cfg.App.DbUrl)
	database.RunSQL(db, migrateSourceName)
	database.RunSQL(db, seedSourceName)
	usersRepo = usersrepo.NewUsersRepo(db)
	authRepo = authrepo.NewAuthRepo(db)
	authSvc = NewAuthSvc(authRepo, usersRepo, cfg)
}

func TestSignUp(t *testing.T) {
	t.Run("unique username", func(t *testing.T) {
		passport, err := authSvc.SignUp("DeepAung", "newuser@gmail.com", "password")
		asserts.EqualError(t, err, usersrepo.ErrUsernameUnique)
		asserts.Equal(t, "passport", passport, auth.Passport{})
	})

	t.Run("unique email", func(t *testing.T) {
		passport, err := authSvc.SignUp("newuser", "i.deepaung@gmail.com", "password")
		asserts.EqualError(t, err, usersrepo.ErrEmailUnique)
		asserts.Equal(t, "passport", passport, auth.Passport{})
	})

	t.Run("normal sign up", func(t *testing.T) {
		passport, err := authSvc.SignUp("newuser", "newuser@gmail.com", "password")
		asserts.EqualError(t, err, nil)
		assertNewUserPassport(t, passport)
	})
}

func TestSignIn(t *testing.T) {
	t.Run("non-exist email", func(t *testing.T) {
		passport, err := authSvc.SignIn("nonexist@gmail.com", "password")
		asserts.EqualError(t, err, usersrepo.ErrUserNotFound)
		asserts.Equal(t, "passport", passport, auth.Passport{})
	})

	t.Run("empty password", func(t *testing.T) {
		passport, err := authSvc.SignIn("newuser@gmail.com", "")
		asserts.EqualError(t, err, ErrInvalidEmailOrPassword)
		asserts.Equal(t, "passport", passport, auth.Passport{})
	})

	t.Run("invalid password", func(t *testing.T) {
		passport, err := authSvc.SignIn("newuser@gmail.com", "invalid")
		asserts.EqualError(t, err, ErrInvalidEmailOrPassword)
		asserts.Equal(t, "passport", passport, auth.Passport{})
	})

	t.Run("normal sign in", func(t *testing.T) {
		passport, err := authSvc.SignIn("newuser@gmail.com", "password")
		asserts.EqualError(t, err, nil)
		assertNewUserPassport(t, passport)
	})
}

func TestSignOut(t *testing.T) {
	t.Run("invalid token id", func(t *testing.T) {
		err := authSvc.SignOut(1000)
		asserts.EqualError(t, err, authrepo.ErrTokenNotFound)
	})

	t.Run("normal sign out", func(t *testing.T) {
		passport, err := authSvc.SignIn("newuser@gmail.com", "password")
		asserts.EqualError(t, err, nil)
		err = authSvc.SignOut(passport.Token.Id)
		asserts.EqualError(t, err, nil)
	})
}

func TestUpdateTokens(t *testing.T) {
	passport, err := authSvc.SignIn("newuser@gmail.com", "password")
	asserts.EqualError(t, err, nil)

	t.Run("invalid token id", func(t *testing.T) {
		token, err := authSvc.UpdateTokens(1000, passport.Token.RefreshToken)
		asserts.EqualError(t, err, ErrInvalidRefreshTokenOrTokenNotFound)
		asserts.Equal(t, "token", token, auth.Token{})
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		token, err := authSvc.UpdateTokens(passport.Token.Id, "invalid refresh token")
		asserts.EqualError(t, err, ErrInvalidRefreshTokenOrTokenNotFound)
		asserts.Equal(t, "token", token, auth.Token{})
	})

	t.Run("empty refresh token", func(t *testing.T) {
		token, err := authSvc.UpdateTokens(passport.Token.Id, "")
		asserts.EqualError(t, err, ErrInvalidRefreshTokenOrTokenNotFound)
		asserts.Equal(t, "token", token, auth.Token{})
	})

	t.Run("normal update tokens", func(t *testing.T) {
		token, err := authSvc.UpdateTokens(passport.Token.Id, passport.Token.RefreshToken)
		asserts.EqualError(t, err, nil)

		asserts.Equal(t, "token id", token.Id, passport.Token.Id)

		// Equal because jwt use time precision of one second
		asserts.Equal(t, "access token", token.AccessToken, passport.Token.AccessToken)
		asserts.Equal(t, "refresh token", token.RefreshToken, passport.Token.RefreshToken)
	})

	err = authSvc.SignOut(passport.Token.Id)
	asserts.EqualError(t, err, nil)
}

func assertNewUserPassport(t *testing.T, passport auth.Passport) {
	asserts.Equal(t, "", passport.User.Username, "newuser")
	asserts.Equal(t, "", passport.User.Email, "newuser@gmail.com")
	asserts.Equal(t, "", passport.User.PictureUrl, "")
	asserts.Equal(t, "", passport.User.IsAdmin, false)
	asserts.NotEqual(t, "", passport.Token.AccessToken, "")
	asserts.NotEqual(t, "", passport.Token.RefreshToken, "")
}
