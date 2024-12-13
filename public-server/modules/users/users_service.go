package users

import (
	"fmt"
	"mime/multipart"

	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/pkg/storer"
	"github.com/DeepAung/gradient/public-server/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

var ErrInvalidCurrentPassword = fiber.NewError(fiber.StatusBadRequest, "invalid current password")

type usersSvc struct {
	repo   types.UsersRepo
	storer storer.Storer
	cfg    *config.Config
}

func NewUsersSvc(repo types.UsersRepo, storer storer.Storer, cfg *config.Config) types.UsersSvc {
	return &usersSvc{
		repo:   repo,
		storer: storer,
		cfg:    cfg,
	}
}

func (s *usersSvc) GetUser(id int) (types.User, error) {
	return s.repo.FindOneUserById(id)
}

// TODO: rollback
// TODO: write test
func (s *usersSvc) ReplacePicture(
	id int,
	email string,
	picture *multipart.FileHeader,
) (string, error) {
	user, err := s.GetUser(id)
	if err != nil {
		return "", err
	}

	if user.PictureUrl != "" {
		// Delete old picture
		res := storer.NewFileResFromUrl(user.PictureUrl)
		if err := s.storer.Delete(res.Dest); err != nil {
			return "", err
		}
	}

	// Upload new picture
	encryped, err := utils.Encrypt(email, s.cfg.App.AesSecretKey)
	if err != nil {
		return "", err
	}

	dest := fmt.Sprintf("users/%s/%s", string(encryped), picture.Filename)

	res, err := s.storer.UploadMultipart(picture, dest, true)
	if err != nil {
		return "", err
	}

	return res.Url, nil
}

func (s *usersSvc) UpdateUser(id int, req types.UpdateUserReq) error {
	if req.CurrentPassword != nil {
		// check password
		currentHashedPassword, err := s.repo.FindOneUserPasswordById(id)
		if err != nil {
			return err
		}
		if !utils.CompareHash(*req.CurrentPassword, currentHashedPassword) {
			return ErrInvalidCurrentPassword
		}

		newHashedPassword, err := utils.Hash(*req.NewPassword)
		if err != nil {
			return err
		}
		req.NewPassword = &newHashedPassword
	}

	return s.repo.UpdateUser(id, req)
}

func (s *usersSvc) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
