package crazy_nl_backend

import (
	"net/http"
	"gopkg.in/matryer/respond.v1"
	"errors"
	"crazy_nl_backend/helpers"
)

type NewRegistration struct {
	Token string `json:"token"`
}

func (registration *NewRegistration) OK () error {
	// todo talk to pushy guys and see what's good value for this one.
	if len(registration.Token) < 20 {
		return errors.New("registration token invalid")
	}
	return nil
}

func (s *Server) RegisterHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registration, err := decodeRegistrationToken(r)
		if err != nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		err = helpers.QueueDeviceRegistration(registration.Token, s.Redis)

		if err != nil {
			s.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}
		// todo use a cron job to check if tokens are actually correct, then move them to db
		respond.With(w, r, 200, struct {
			Status string `json:"status"`
		}{
			Status: "registration pending",
		})
	})
}

func decodeRegistrationToken(r *http.Request) (*NewRegistration, error) {
	registration := new(NewRegistration)
	Decode(r, &registration)
	err := registration.OK()

	if err != nil {
		return nil, err
	}

	return registration, nil
}
