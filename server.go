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
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func Init() {
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Request-Headers", "Origin", "authorization"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	server := createServer()
	defer server.cleanup()
	router := NewRouter(server)
	server.Logger.Info("Attempting to listen on port " + strconv.Itoa(config.GetApiPort()))
	err := http.ListenAndServe(":"+strconv.Itoa(config.GetApiPort()), handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router))
	if err != nil {
		server.Logger.Fatal(err)
		panic(err)
	}
}

func (s *Server) ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	respond.With(w, r, statusCode, ErrorResponse{
		Message: message,
	})
	return
}

func createServer() Server {
	session := getMongo()
	dbLayer := db.GetDbImplementation(session)
	return Server{
		Logger: getLogger(),
		Db:     dbLayer,
		Redis:  *getRedis(),
		Pushy:  getPushy(),
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
