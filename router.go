package starter

import (
	"github.com/fossapps/starter/middleware"

	"net/http"
	"time"

	"github.com/fossapps/starter/rate"
	"github.com/gorilla/mux"
	"gopkg.in/matryer/respond.v1"
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
	r.router.HandleFunc("/device/register", middleware.Adapt(r.server.RegisterHandler(), middleware.ResponseTime(r.server.Logger))).
		Methods("POST")
}

func (r router) userResource() {
	r.router.HandleFunc(
		"/users",
		middleware.Adapt(r.server.CreateUser(), middleware.MustHavePermission(r.perm.User.Create))).
		Methods("POST")

	r.router.HandleFunc(
		"/users", middleware.Adapt(r.server.ListUsers(),
			middleware.MustHavePermission(r.perm.User.List))).
		Methods("GET")

	r.router.Handle("/users/available", r.server.UserAvailability()).
		Methods("POST")
	r.router.Handle(
		"/users/{user}",
		middleware.Adapt(r.server.UpdateUser(), middleware.MustHavePermission(r.perm.User.Edit))).
		Methods("PUT")
	r.router.Handle("/users/{user}", r.server.GetUser()).Methods("GET")
}

func (r router) permissionsResource() {
	r.router.HandleFunc(
		"/permissions", middleware.Adapt(r.server.ListPermissions(),
			middleware.AuthMw(r.server.ReqHelper))).
		Methods("GET")
}

func (r router) authResource() {
	r.router.HandleFunc("/session/", r.server.LoginHandler()).
		Methods("POST")

	r.router.HandleFunc("/session/refresh", r.server.RefreshTokenHandler()).
		Methods("POST")

	r.router.HandleFunc(
		"/session",
		middleware.Adapt(r.server.RefreshTokensList(), middleware.AuthMw(r.server.ReqHelper))).
		Methods("GET")
	r.router.HandleFunc("/session/{token}", r.server.DeleteSession()).Methods("DELETE")
}

// NewRouter returns a configured router with path and handlers
func NewRouter(s Server) *mux.Router {
	muxRouter := mux.NewRouter()
	routerInstance := router{
		server: s,
		router: muxRouter,
	}
	routerInstance.build()
	limiterOptions := middleware.LimiterOptions{
		RequestHelper: routerInstance.server.ReqHelper,
		Limit:         3,
		Namespace:     "test",
		AddHeaders:    true,
		Logger:        s.Logger,
		Limiter: rate.Limiter{
			Decay:       10 * time.Second,
			Limit:       3,
			RedisClient: routerInstance.server.Redis,
		},
	}
	routerInstance.router.HandleFunc("/test", middleware.Adapt(func(writer http.ResponseWriter, request *http.Request) {
		respond.With(writer, request, http.StatusOK, "success")
	}, middleware.Limit(limiterOptions)))
	return routerInstance.router
}
