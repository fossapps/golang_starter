package crazy_nl_backend

import (
	"crazy_nl_backend/adapters"

	"github.com/gorilla/mux"
)

func NewRouter(s Server) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/device/register", adapters.Adapt(s.RegisterHandler(), adapters.ResponseTime(s.Logger))).
		Methods("POST")
	router.HandleFunc("/session/", s.LoginHandler()).
		Methods("POST")
	router.HandleFunc("/session/refresh", s.RefreshTokenHandler()).Methods("POST")
	router.HandleFunc("/permissions", adapters.Adapt(s.ListPermissions(), adapters.MustHavePermission("Permission.List")))
	return router
}
