package starter

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fossapps/starter/db"
	"github.com/fossapps/starter/transformer"
	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"gopkg.in/matryer/respond.v1"
)

// NewUser for creating a new user
type NewUser struct {
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
}

// UserAvailabilityResponse is used to respond for availability requests
type UserAvailabilityResponse struct {
	Available bool `json:"available"`
}

// UserAvailabilityRequest used for making request asking if email is available
type UserAvailabilityRequest struct {
	Email string `json:"email"`
}

// Ok returns if user is valid
func (user NewUser) Ok() bool {
	if !strings.Contains(user.Email, "@") || (len(user.Password) < 6 && user.Password != "") {
		return false
	}
	return true
}

// CreateUser handler used for creating new users
func (s Server) CreateUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newUser := new(NewUser)
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, "invalid user")
			return
		}
		if !newUser.Ok() {
			s.ErrorResponse(w, r, http.StatusBadRequest, "invalid user")
			return
		}
		database := s.Db.Clone()
		defer database.Close()
		user, err := database.Users().FindByEmail(newUser.Email)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "server error")
			return
		}
		if user != nil {
			s.ErrorResponse(w, r, http.StatusConflict, "duplicate registration")
			return
		}
		validUser := db.User{
			Email:       newUser.Email,
			Password:    newUser.Password,
			Permissions: newUser.Permissions,
		}
		err = database.Users().Create(validUser)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "internal server error")
			return
		}
		s.SuccessResponse(w, r, http.StatusCreated, "created")
	})
}

// ListUsers used for listing all users
func (s Server) ListUsers() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		database := s.Db.Clone()
		defer database.Close()
		users, err := database.Users().List()
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "internal server error")
			return
		}
		respond.With(w, r, http.StatusOK, transformer.TransformUsers(users))
	})
}

// UserAvailability for checking if a email is already taken
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
		user, err := database.Users().FindByEmail(requestedUser.Email)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "server error")
			return
		}
		respond.With(w, r, http.StatusOK, UserAvailabilityResponse{
			Available: user == nil,
		})
	})
}

// UpdateUser to update information about user
func (s Server) UpdateUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["user"]
		user, err := s.Db.Users().FindByID(userID)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "server error")
			return
		}
		if user == nil {
			s.ErrorResponse(w, r, http.StatusPreconditionFailed, "user not found")
			return
		}
		newUser := new(NewUser)
		err = json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil || !newUser.Ok() {
			s.ErrorResponse(w, r, http.StatusBadRequest, "invalid user")
			return
		}
		validUser := db.User{
			Permissions: newUser.Permissions,
			Email:       newUser.Email,
			Password:    newUser.Password,
		}
		err = s.Db.Users().Update(userID, validUser)
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

// GetUser to get information about a user
func (s Server) GetUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["user"]
		user, err := s.Db.Users().FindByID(id)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, "server error")
			return
		}
		if user == nil {
			s.ErrorResponse(w, r, http.StatusNotFound, "not found")
			return
		}
		respond.With(w, r, http.StatusOK, transformer.TransformUser(*user))
	})
}
