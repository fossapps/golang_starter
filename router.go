package crazy_nl_backend

import (
	"github.com/gorilla/mux"
	"crazy_nl_backend/adapters"
)

func NewRouter(s Server) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/device/register", adapters.Adapt(s.RegisterHandler(), adapters.ResponseTime(s.Logger))).
		Methods("POST")
	router.HandleFunc("/session/", s.LoginHandler()).
		Methods("POST")
	router.HandleFunc("/session/refresh", s.RefreshTokenHandler()).Methods("POST")
	return router
}
