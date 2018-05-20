package starter

import (
	"github.com/fossapps/starter/config"
	"github.com/fossapps/starter/db"
	"github.com/fossapps/starter/helpers"
	"net/http"
	"strconv"
	"time"

	"github.com/cyberhck/pushy"
	"github.com/globalsign/mgo"
	"github.com/gorilla/handlers"
	"github.com/multiplay/go-slack/chat"
	"github.com/multiplay/go-slack/lrhook"
	"github.com/sirupsen/logrus"
	"gopkg.in/matryer/respond.v1"
	"github.com/fossapps/starter/adapters"
	"github.com/dgrijalva/jwt-go/request"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"net"
)

// ILogger needs to be implemented for a logger to be used on this project
type ILogger interface {
	Info(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Print(args ...interface{})
	Panic(args ...interface{})
}

// Server is a global struct which holds implementation of different things application depends on.
// One can think of this as dependency container
type Server struct {
	Logger    ILogger
	Db        db.Db
	Redis     helpers.IRedisClient
	Pushy     pushy.IPushyClient
	ReqHelper IRequestHelper
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
	router := adapters.DevMw(1000)(NewRouter(server))
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
	requestHelper := RequestHelper{}
	return Server{
		Logger:    getLogger(),
		Db:        dbLayer,
		Redis:     *getRedis(),
		Pushy:     getPushy(),
		ReqHelper: requestHelper,
	}
}

func getLogger() ILogger {
	logger := logrus.New()
	level, err := logrus.ParseLevel(config.GetLogLevel())
	logger.AddHook(getSlackHook())
	if err != nil {
		panic(err)
	}
	logger.SetLevel(level)
	return logger
}

func getSlackHook() *lrhook.Hook {
	cfg := lrhook.Config{
		MinLevel: logrus.WarnLevel,
		Message: chat.Message{
			Channel:   "#general",
			IconEmoji: ":gopher:",
		},
	}
	return lrhook.New(cfg, config.GetApplicationConfig().SlackLoggingAppConfig)
}

func getMongo() *mgo.Session {
	mongo, err := mgo.Dial(config.GetMongoConfig().Connection)
	if err != nil {
		panic(err)
	}
	return mongo
}

func getRedis() *helpers.IRedisClient {
	redis, err := helpers.GetRedis()
	if err != nil {
		panic(err)
	}
	return &redis
}

func getPushy() pushy.IPushyClient {
	sdk := pushy.Create(config.GetPushyToken(), pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(5 * time.Second))
	return sdk
}

func (s *Server) cleanup() {
	s.Redis.Close()
}

// RequestHelper simple IRequestHelper implementation
type RequestHelper struct{}

// GetJwtData util method to get data from request
func (RequestHelper) GetJwtData(r *http.Request) (*adapters.Claims, error) {
	var claims adapters.Claims
	token, parseErr := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, &claims, signingFunc)
	err := claims.Valid()
	if parseErr != nil {
		return nil, parseErr
	}
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &claims, nil
}

// GetIPAddress util method to get IP address from request
func (RequestHelper) GetIPAddress(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}

// IRequestHelper interface to implement to satisfy as a Request Helper for this application
type IRequestHelper interface {
	GetJwtData(r *http.Request) (*adapters.Claims, error)
	GetIPAddress(r *http.Request) string
}

func signingFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(config.GetApplicationConfig().JWTSecret), nil
}
