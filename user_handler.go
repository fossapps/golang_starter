package crazy_nl_backend

import (
	"net/http"
)

type NewUser struct {
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
}

func (user NewUser) Ok() bool {
	// check if email is in correct format
	// check if password is strong enough
	// don't need to check if permissions are correct, even if admin tries to add some random permission
	// that permission won't do anything anyway
}

func (s Server) Create() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		database := s.Db.Clone()
		defer database.Close()
		// decode from body
		// make sure the validation passes
		// check if it already exists
		// attempt to add to db
		// handle error
	})
}
