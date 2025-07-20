package types

import "mime/multipart"

type UsersSvc interface {
	GetUser(id int) (User, error)
	ReplacePicture(id int, email string, picture *multipart.FileHeader) (string, error)
	UpdateUser(id int, req UpdateUserReq) error
	DeleteUser(id int) error
}

type UsersRepo interface {
	CreateUser(username, email, hashedPassword string) (User, error)
	FindOneUserById(id int) (User, error)
	FindOneUserWithPasswordByEmail(email string) (UserWithPassword, error)
	FindOneUserPasswordById(id int) (string, error)
	UpdateUser(id int, req UpdateUserReq) error
	DeleteUser(id int) error
}

type User struct {
	Id         int    `db:"id"`
	Username   string `db:"username"`
	Email      string `db:"email"`
	PictureUrl string `db:"picture_url"`
	IsAdmin    bool   `db:"is_admin"`
}

type UserWithPassword struct {
	Id         int    `db:"id"`
	Username   string `db:"username"`
	Email      string `db:"email"`
	Password   string `db:"password"`
	PictureUrl string `db:"picture_url"`
	IsAdmin    bool   `db:"is_admin"`
}

type UpdateUserReq struct {
	Username   string  `db:"username"`
	PictureUrl *string `db:"picture_url"` // optional field

	// composite optional field?
	CurrentPassword *string
	NewPassword     *string `db:"password"`
}
