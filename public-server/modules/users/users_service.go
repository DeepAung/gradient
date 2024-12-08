package users

import (
	"errors"

	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/pkg/utils"
)

var ErrInvalidCurrentPassword = errors.New("invalid current password")

type UsersSvc struct {
	repo types.UsersRepo
}

func NewUsersSvc(repo types.UsersRepo) types.UsersSvc {
	return &UsersSvc{
		repo: repo,
	}
}

func (s *UsersSvc) GetUser(id int) (types.User, error) {
	return s.repo.FindOneUserById(id)
}

func (s *UsersSvc) UpdateUser(id int, req types.UpdateUserReq) error {
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

func (s *UsersSvc) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
