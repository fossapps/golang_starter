package crazy_nl_backend

import (
	"crazy_nl_backend/config"
	"crazy_nl_backend/db"
	"crazy_nl_backend/helpers"
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
	"crazy_nl_backend/adapters"
	"github.com/dgrijalva/jwt-go/request"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"net"
)

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

type Server struct {
	Logger ILogger
	Db     db.Db
	Redis  helpers.IRedisClient
	Pushy  pushy.IPushyClient
	ReqHelper IRequestHelper
}

type SimpleResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

func Init() {
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Request-Headers", "Origin", "authorization"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	server := createServer()
	defer server.cleanup()
	router := adapters.DevMw(1000)(NewRouter(server))
	port := strconv.Itoa(config.GetApiPort())
	server.Logger.Info("Attempting to listen on port " + port)
	err := http.
		ListenAndServe(":" + port, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router))
	if err != nil {
		server.Logger.Fatal(err)
		panic(err)
	}
}

func (s *Server) ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	respond.With(w, r, statusCode, SimpleResponse{
		Success:false,
		Message: message,
	})
	return
}

func (s *Server) SuccessResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	respond.With(w, r, statusCode, SimpleResponse{
		Success:true,
		Message: message,
	})
}

func createServer() Server {
	session := getMongo()
	dbLayer := db.GetDbImplementation(session)
	requestHelper := RequestHelper{}
	return Server{
		Logger: getLogger(),
		Db:     dbLayer,
		Redis:  *getRedis(),
		Pushy:  getPushy(),
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

type RequestHelper struct {}

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

func (RequestHelper) GetIpAddress(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}

type IRequestHelper interface {
	GetJwtData(r *http.Request) (*adapters.Claims, error)
	GetIpAddress(r *http.Request) string
}

func signingFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
	}
	return []byte(config.GetApplicationConfig().JWTSecret), nil
}
