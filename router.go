package crazy_nl_backend

import (
	"crazy_nl_backend/adapters"

	"github.com/gorilla/mux"
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
		Methods("PUT")

	r.router.HandleFunc(
		"/users", adapters.Adapt(r.server.ListUsers(),
			adapters.MustHavePermission(r.perm.User.List))).
		Methods("GET")

	r.router.Handle("/users/available", r.server.UserAvailability()).
		Methods("POST")
}

func (r router) permissionsResource() {
	r.router.HandleFunc(
		"/permissions", adapters.Adapt(r.server.ListPermissions(),
			adapters.MustHavePermission(r.perm.Permissions.List))).
		Methods("GET")
}

func (r router) authResource() {
	r.router.HandleFunc("/session/", r.server.LoginHandler()).
		Methods("POST")

	r.router.HandleFunc("/session/refresh", r.server.RefreshTokenHandler()).
		Methods("POST")

	// todo: make a middleware which requires for someone to be logged in (but no permission) and add it here.
	r.router.HandleFunc("/session", r.server.RefreshTokensList()).Methods("GET")
	r.router.HandleFunc("/session/{token}", r.server.DeleteSession()).Methods("DELETE")
}

func NewRouter(s Server) *mux.Router {
	muxRouter := mux.NewRouter()
	routerInstance := router{
		server: s,
		router: muxRouter,
	}
	routerInstance.build()
	return routerInstance.router
}
