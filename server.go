package starter

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/cyberhck/pushy"
	"github.com/fossapps/starter/config"
	"github.com/fossapps/starter/db"
	"github.com/fossapps/starter/jwt"
	"github.com/fossapps/starter/logger"
	"github.com/fossapps/starter/middleware"
	"github.com/fossapps/starter/redis"
	"github.com/globalsign/mgo"
	"github.com/gorilla/handlers"
	"gopkg.in/matryer/respond.v1"
)

// Server is a global struct which holds implementation of different things application depends on.
// One can think of this as dependency container
type Server struct {
	Logger    logger.Client
	Db        db.DB
	Redis     redis.Client
	Pushy     pushy.IPushyClient
	ReqHelper RequestHelper
	Jwt       jwt.Manager
}

// SimpleResponse is used when our handler wants to responds with simple boolean type information, like success, or fail
type SimpleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Init actually starts listening and server on a port
func Init() {
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Request-Headers", "Origin", "authorization"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	server := createServer()
	defer server.cleanup()
	router := middleware.DevMw(1000)(NewRouter(server))
	port := strconv.Itoa(config.GetAPIPort())
	server.Logger.Info("Attempting to listen on port " + port)
	err := http.
		ListenAndServe(":"+port, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router))
	if err != nil {
		server.Logger.Fatal(err)
		panic(err)
	}
}

// ErrorResponse util method to indicate failure
func (s *Server) ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	respond.With(w, r, statusCode, SimpleResponse{
		Success: false,
		Message: message,
	})
	return
}

// SuccessResponse util method to indicate success
func (s *Server) SuccessResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	respond.With(w, r, statusCode, SimpleResponse{
		Success: true,
		Message: message,
	})
}

func createServer() Server {
	session := getMongo()
	dbLayer := db.GetDbImplementation(session)
	requestHelper := requestHelper{}
	jwtManager := jwt.Client{
		Config: jwt.Config{
			Secret: config.GetApplicationConfig().JWTSecret,
			Expiry: config.GetApplicationConfig().JWTExpiryTime,
		},
	}
	return Server{
		Logger:    logger.GetClient(config.GetLogLevel()),
		Db:        dbLayer,
		Redis:     *getRedis(),
		Pushy:     getPushy(),
		ReqHelper: requestHelper,
		Jwt:       jwtManager,
	}
}

func getMongo() *mgo.Session {
	mongo, err := mgo.Dial(config.GetMongoConfig().Connection)
	if err != nil {
		panic(err)
	}
	return mongo
}

func getRedis() *redis.Client {
	client, err := redis.NewClient()
	if err != nil {
		panic(err)
	}
	return &client
}

func getPushy() pushy.IPushyClient {
	sdk := pushy.Create(config.GetPushyToken(), pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(5 * time.Second))
	return sdk
}

func (s *Server) cleanup() {
	s.Redis.Close()
}

// requestHelper simple requestHelper implementation
type requestHelper struct{}

// GetIPAddress util method to get IP address from request
func (requestHelper) GetIPAddress(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}

// RequestHelper interface to implement to satisfy as a Request Helper for this application
type RequestHelper interface {
	GetIPAddress(r *http.Request) string
}
