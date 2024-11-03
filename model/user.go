package model

import "encoding/json"

type User struct {
	ID                   string `json:"id"`
	UUID                 string `json:"uuid"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,max=32"`
	PasswordConfirmation string `json:"password_confirmation" validate:"omitempty,min=8,max=32,eqfield=Password"`
	UserEmailVerified    string `json:"user_email_verified" db:"user_email_verified"`
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(u)
}
