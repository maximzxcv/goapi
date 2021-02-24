package user

// User ....
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateUser ....
type CreateUser struct {
	Name            string `json:"name" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm" binding:"eqfield=Password"`
}
