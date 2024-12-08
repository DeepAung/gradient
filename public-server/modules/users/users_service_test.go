package users

import (
	"testing"

	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/database"
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/pkg/asserts"
	"github.com/DeepAung/gradient/public-server/pkg/utils"
	"github.com/jmoiron/sqlx"
)

var (
	migrateSourceName = "../../migrations/migrate.sql"
	seedSourceName    = "../../migrations/seed.sql"
	cfg               *config.Config
	db                *sqlx.DB
	repo              types.UsersRepo
	svc               types.UsersSvc
)

func init() {
	cfg = config.NewConfig("../../.env.dev")
	db = database.InitDB(cfg.App.DbUrl)
	database.RunSQL(db, migrateSourceName)
	database.RunSQL(db, seedSourceName)
	repo = NewUsersRepo(db)
	svc = NewUsersSvc(repo)
}

func TestGetUser(t *testing.T) {
	t.Run("id not found", func(t *testing.T) {
		user, err := svc.GetUser(1000)
		asserts.EqualError(t, err, ErrUserNotFound)
		asserts.Equal(t, "user", user, types.User{})
	})

	t.Run("normal get user", func(t *testing.T) {
		user, err := svc.GetUser(1)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "user", user, types.User{
			Id:         1,
			Username:   "DeepAung",
			Email:      "i.deepaung@gmail.com",
			PictureUrl: "",
			IsAdmin:    false,
		})
	})
}

func TestUpdateUser(t *testing.T) {
	pictureUrl := "newPictureUrl"
	currentPassword := "password"
	invalidCurrentPassword := "invalid password here"
	newPassword := "newPassword"

	t.Run("id not found", func(t *testing.T) {
		err := svc.UpdateUser(1000, types.UpdateUserReq{
			CurrentPassword: nil,
		})
		asserts.EqualError(t, err, ErrUserNotFound)
	})

	t.Run("id not found 02", func(t *testing.T) {
		err := svc.UpdateUser(1000, types.UpdateUserReq{
			CurrentPassword: &currentPassword,
		})
		asserts.EqualError(t, err, ErrUserNotFound)
	})

	t.Run("invalid current password", func(t *testing.T) {
		err := svc.UpdateUser(1, types.UpdateUserReq{
			Username:        "NewDeepAung",
			PictureUrl:      &pictureUrl,
			CurrentPassword: &invalidCurrentPassword,
			NewPassword:     &newPassword,
		})
		asserts.EqualError(t, err, ErrInvalidCurrentPassword)
	})

	t.Run("partial update", func(t *testing.T) {
		err := svc.UpdateUser(1, types.UpdateUserReq{
			Username:        "NewDeepAung",
			PictureUrl:      nil,
			CurrentPassword: nil,
			NewPassword:     nil,
		})
		asserts.EqualError(t, err, nil)

		user, err := repo.FindOneUserWithPasswordByEmail("i.deepaung@gmail.com")
		asserts.EqualError(t, err, nil)

		asserts.Equal(t, "user id", user.Id, 1)
		asserts.Equal(t, "username", user.Username, "NewDeepAung")
		asserts.Equal(t, "picture url", user.PictureUrl, "")

		if !utils.CompareHash(currentPassword, user.Password) {
			t.Fatalf(
				"utils.CompareHash password, expect=%v, got=%v",
				currentPassword, user.Password,
			)
		}
	})

	t.Run("full update", func(t *testing.T) {
		err := svc.UpdateUser(1, types.UpdateUserReq{
			Username:        "NewDeepAung02",
			PictureUrl:      &pictureUrl,
			CurrentPassword: &currentPassword,
			NewPassword:     &newPassword,
		})
		asserts.EqualError(t, err, nil)

		user, err := repo.FindOneUserWithPasswordByEmail("i.deepaung@gmail.com")
		asserts.EqualError(t, err, nil)

		asserts.Equal(t, "user id", user.Id, 1)
		asserts.Equal(t, "username", user.Username, "NewDeepAung02")
		asserts.Equal(t, "picture url", user.PictureUrl, pictureUrl)

		if !utils.CompareHash(newPassword, user.Password) {
			t.Fatalf("utils.CompareHash password, expect=%v, got=%v", newPassword, user.Password)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("id not found", func(t *testing.T) {
		err := svc.DeleteUser(1000)
		asserts.EqualError(t, err, ErrUserNotFound)
	})

	t.Run("normal delete", func(t *testing.T) {
		err := svc.DeleteUser(1)
		asserts.EqualError(t, err, nil)
	})
}
