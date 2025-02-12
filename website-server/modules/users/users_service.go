package users

import (
	"fmt"
	"mime/multipart"
	"path"

	"github.com/DeepAung/gradient/website-server/config"
	"github.com/DeepAung/gradient/website-server/modules/types"
	"github.com/DeepAung/gradient/website-server/pkg/storer"
	"github.com/DeepAung/gradient/website-server/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

var ErrInvalidCurrentPassword = fiber.NewError(fiber.StatusBadRequest, "invalid current password")

type usersSvcImpl struct {
	repo   types.UsersRepo
	storer storer.Storer
	cfg    *config.Config
}

func NewUsersSvc(repo types.UsersRepo, storer storer.Storer, cfg *config.Config) types.UsersSvc {
	return &usersSvcImpl{
		repo:   repo,
		storer: storer,
		cfg:    cfg,
	}
}

func (s *usersSvcImpl) GetUser(id int) (types.User, error) {
	return s.repo.FindOneUserById(id)
}

// TODO: rollback
// TODO: write test
func (s *usersSvcImpl) ReplacePicture(
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

	// Hash picture path(email) and picture filename
	hashedEmail, err := utils.Hash(email)
	if err != nil {
		return "", err
	}
	ext := path.Ext(picture.Filename) // e.g. ".png"
	hashedName, err := utils.Hash(picture.Filename)
	if err != nil {
		return "", err
	}
	hashedFilename := hashedName + ext

	// Upload new picture
	dest := fmt.Sprintf("users/%s/%s", hashedEmail, hashedFilename)
	res, err := s.storer.UploadMultipart(picture, dest, true)
	if err != nil {
		return "", err
	}

	return res.Url, nil
}

func (s *usersSvcImpl) UpdateUser(id int, req types.UpdateUserReq) error {
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

func (s *usersSvcImpl) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
