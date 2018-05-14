package crazy_nl_backend

import (
	"crazy_nl_backend/adapters"

	"github.com/gorilla/mux"
	"time"
	"net/http"
	"gopkg.in/matryer/respond.v1"
	"crazy_nl_backend/helpers"
)

type router struct {
	router *mux.Router
	perm   Permissions
	server Server
}

func (r router) build() {
	r.perm = Const().Permissions
	r.userResource()
	r.deviceResource()
	r.permissionsResource()
	r.authResource()
}

func (r router) deviceResource() {
	r.router.HandleFunc("/device/register", adapters.Adapt(r.server.RegisterHandler(), adapters.ResponseTime(r.server.Logger))).
		Methods("POST")
}

func (r router) userResource() {
	r.router.HandleFunc(
		"/users",
		adapters.Adapt(r.server.CreateUser(), adapters.MustHavePermission(r.perm.User.Create))).
		Methods("POST")

	r.router.HandleFunc(
		"/users", adapters.Adapt(r.server.ListUsers(),
			adapters.MustHavePermission(r.perm.User.List))).
		Methods("GET")

	r.router.Handle("/users/available", r.server.UserAvailability()).
		Methods("POST")
	r.router.Handle(
		"/users/{user}",
		adapters.Adapt(r.server.EditUser(), adapters.MustHavePermission(r.perm.User.Edit))).
		Methods("PUT")
	r.router.Handle("/users/{user}", r.server.GetUser()).Methods("GET")
}

func (r router) permissionsResource() {
	r.router.HandleFunc(
		"/permissions", adapters.Adapt(r.server.ListPermissions(),
			adapters.AuthMw(r.server.ReqHelper))).
		Methods("GET")
}

func (r router) authResource() {
	r.router.HandleFunc("/session/", r.server.LoginHandler()).
		Methods("POST")

	r.router.HandleFunc("/session/refresh", r.server.RefreshTokenHandler()).
		Methods("POST")

	r.router.HandleFunc(
		"/session",
		adapters.Adapt(r.server.RefreshTokensList(), adapters.AuthMw(r.server.ReqHelper))).
		Methods("GET")
	r.router.HandleFunc("/session/{token}", r.server.DeleteSession()).Methods("DELETE")
}

func NewRouter(s Server) *mux.Router {
	muxRouter := mux.NewRouter()
	routerInstance := router{
		server: s,
		router: muxRouter,
	}
	routerInstance.build()
	limiterOptions := adapters.LimiterOptions{
		RequestHelper: routerInstance.server.ReqHelper,
		Limit:         3,
		Namespace:     "test",
		AddHeaders:    true,
		Logger: s.Logger,
		Limiter: helpers.Limiter{
			Decay: 10 * time.Second,
			Limit: 3,
			RedisClient: routerInstance.server.Redis,
		},
	}
	routerInstance.router.HandleFunc("/test", adapters.Adapt(func(writer http.ResponseWriter, request *http.Request) {
		respond.With(writer, request, http.StatusOK, "success")
	}, adapters.Limit(limiterOptions)))
	return routerInstance.router
}
