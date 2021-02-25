package user

// User ....
type User struct {
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Password []byte `json:"-" db:"password"`
}

// CreateUser ....
type CreateUser struct {
	Name            string `json:"name" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm" binding:"eqfield=Password"`
}

// UpdateUser ....
type UpdateUser struct {
	CreateUser
}
