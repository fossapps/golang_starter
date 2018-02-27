package crazy_nl_backend

import (
	"net/http"

	"crazy_nl_backend/adapters"

	"github.com/gorilla/mux"
	"gopkg.in/matryer/respond.v1"
)

func NewRouter(s Server) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/device/register", adapters.Adapt(s.RegisterHandler(), adapters.ResponseTime(s.Logger))).
		Methods("POST")
	router.HandleFunc("/session/", s.LoginHandler()).
		Methods("POST")
	router.HandleFunc("/session/refresh", s.RefreshTokenHandler()).Methods("POST")
	router.Handle("/test", adapters.Adapt(func(writer http.ResponseWriter, request *http.Request) {
		respond.With(writer, request, http.StatusOK, "testing")
	}, adapters.MustHavePermission("user.create"))).Methods("GET")
	return router
}
