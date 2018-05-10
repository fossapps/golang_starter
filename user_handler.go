package crazy_nl_backend

import (
	"crazy_nl_backend/db"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"gopkg.in/matryer/respond.v1"
)

type NewUser struct {
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
}

type UserAvailabilityResponse struct {
	Available bool `json:"available"`
}

type UserAvailabilityRequest struct {
	Email string `json:"email"`
}

func (user NewUser) Ok() bool {
	if !strings.Contains(user.Email, "@") || (len(user.Password) < 6 && user.Password != "") {
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
			Email:       user.Email,
			Password:    user.Password,
			Permissions: user.Permissions,
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

func (s Server) UserAvailability() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedUser := new(UserAvailabilityRequest)
		err := json.NewDecoder(r.Body).Decode(&requestedUser)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, "bad request")
			return
		}
		database := s.Db.Clone()
		defer database.Close()
		user := database.Users().FindByEmail(requestedUser.Email)
		respond.With(w, r, http.StatusOK, UserAvailabilityResponse{
			Available: user == nil,
		})
	})
}

func (s Server) EditUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := mux.Vars(r)["user"]
		user := s.Db.Users().FindById(userId)
		if user == nil {
			s.ErrorResponse(w, r, http.StatusPreconditionFailed, "user not found")
			return
		}
		newUser := new(NewUser)
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil || !newUser.Ok() {
			s.ErrorResponse(w, r, http.StatusBadRequest, "invalid user")
			return
		}
		validUser := db.User{
			Permissions: newUser.Permissions,
			Email:       newUser.Email,
			Password:    newUser.Password,
		}
		err = s.Db.Users().Edit(userId, validUser)
		if err == mgo.ErrNotFound {
			s.ErrorResponse(w, r, http.StatusPreconditionFailed, "user not found")
			return
		}
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "server error")
			return
		}
		s.SuccessResponse(w, r, http.StatusOK, "user updated")
	})
}

func (s Server) GetUser() http.HandlerFunc {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["user"]
		user := s.Db.Users().FindById(id)
		if user == nil {
			s.ErrorResponse(w, r, http.StatusNotFound, "not found");
			return
		}
		respond.With(w, r, http.StatusOK, user)
	})
}
