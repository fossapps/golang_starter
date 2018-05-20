package starter

import (
	"encoding/json"
	"net/http"

)

// NewRegistration used for making a registration request
type NewRegistration struct {
	Token string `json:"token"`
}

// RegistrationResponse response for registration request
type RegistrationResponse struct {
	Status string `json:"status"`
}

func (registration *NewRegistration) ok() bool {
	if len(registration.Token) < 20 {
		return false
	}
	return true
}

// RegisterHandler is handler for registration requests
func (s *Server) RegisterHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registration := new(NewRegistration)
		err := json.NewDecoder(r.Body).Decode(&registration)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusUnprocessableEntity, "invalid request")
			return
		}
		if !registration.ok() {
			s.ErrorResponse(w, r, http.StatusBadRequest, "invalid token")
			return
		}
		_, pushyErr, err := s.Pushy.DeviceInfo(registration.Token)
		// todo maybe we want to use the cache.Remember thing? so later it'll be a breeze if we want to get data
		if err != nil || pushyErr != nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, "invalid token")
			return
		}
		// save to mongodb
		if s.Db.Devices().Exists(registration.Token) {
			s.ErrorResponse(w, r, http.StatusBadRequest, "already registered")
			return
		}
		err = s.Db.Devices().Register(registration.Token)

		if err != nil {
			s.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		s.SuccessResponse(w, r, http.StatusOK, "success")
	})
}
