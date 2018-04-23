package crazy_nl_backend

import (
	"net/http"
	"encoding/json"
	"gopkg.in/matryer/respond.v1"
)

type NewRegistration struct {
	Token string `json:"token"`
}

type RegistrationResponse struct {
	Status string `json:"status"`
}

func (registration *NewRegistration) ok() bool {
	if len(registration.Token) < 20 {
		return false
	}
	return true
}

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

		// save to mongodb
		if s.Db.Devices().Exists(registration.Token) {
			s.ErrorResponse(w, r, http.StatusBadRequest, "already registered")
			return
		}
		err = s.Db.Devices().Register(registration.Token)

		if err != nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}
		respond.With(w, r, 200, RegistrationResponse{
			Status: "success",
		})
	})
}
