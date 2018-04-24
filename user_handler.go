package crazy_nl_backend

import (
	"net/http"
	"encoding/json"
	"strings"
	"crazy_nl_backend/db"
	"gopkg.in/matryer/respond.v1"
)

type NewUser struct {
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
}

func (user NewUser) Ok() bool {
	if !strings.Contains(user.Email, "@") || len(user.Password) < 6 {
		return false
	}
	return true
}

func (s Server) CreateUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := new(NewUser)
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, "invalid user")
			return
		}
		if !user.Ok() {
			s.ErrorResponse(w, r, http.StatusBadRequest, "invalid user")
			return
		}
		database := s.Db.Clone()
		defer database.Close()
		if database.Users().FindByEmail(user.Email) != nil {
			s.ErrorResponse(w, r, http.StatusConflict, "duplicate registration")
			return
		}
		validUser := db.User{
			Email:user.Email,
			Password:user.Password,
			Permissions:user.Permissions,
		}
		err = database.Users().Create(validUser)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "internal server error")
			return
		}
		s.SuccessResponse(w, r, http.StatusCreated, "created")
	})
}

func (s Server) ListUsers() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		database := s.Db.Clone()
		defer database.Close()
		users, err := database.Users().List()
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "internal server error")
			return
		}
		respond.With(w, r, http.StatusOK, users)
	})
}
