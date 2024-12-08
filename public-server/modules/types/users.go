package types

type UsersRepo interface {
	CreateUser(username, email, hashedPassword string) (User, error)
	FindOneUserWithPasswordByEmail(email string) (UserWithPassword, error)
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
