package models

type User struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Permissions []string `json:"permissions"`
}
