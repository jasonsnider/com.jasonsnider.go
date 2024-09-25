package types

type User struct {
	ID        string `db:"id"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email,uniqueEmail"`
	Role      string `json:"role" validate:"required,oneof=admin user"`
}

type CreateUser struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email,uniqueEmail"`
	Role      string `json:"role" validate:"required,oneof=admin user"`
}

type RegisterUser struct {
	ID              string `db:"id"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Email           string `json:"email" validate:"required,email,uniqueEmail"`
	Password        string `json:"password" validate:"required,min=12"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
