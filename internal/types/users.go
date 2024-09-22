package types

type User struct {
	ID        string `db:"id"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

type RegisterUser struct {
	ID              string `db:"id"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=12"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
