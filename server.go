package crazy_nl_backend

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"gopkg.in/matryer/respond.v1"
	"crazy_nl_backend/helpers"
	"crazy_nl_backend/config"
	"strconv"
	"github.com/cyberhck/pushy"
	"time"
)

type Server struct {
	Logger logrus.Logger
	Mongo  helpers.IMongoClient
	Redis  helpers.IRedisClient
	Pushy  pushy.Pushy
}

func Init() {
	server := createServer()
	defer server.cleanup()
	router := NewRouter(server)
	server.Logger.Info("Attempting to listen on port " + strconv.Itoa(config.GetApiPort()))
	err := http.ListenAndServe(":"+strconv.Itoa(config.GetApiPort()), router)
	if err != nil {
		panic(err)
	}
}

func Decode(r *http.Request, data interface{}) {
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&data)
}

func (s *Server) ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	respond.With(w, r, statusCode, struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	return
}

func createServer() Server {
	db := getMongo()
	return Server{
		Logger: getLogger(),
		Mongo:  db,
		Redis:  *getRedis(),
		Pushy:  getPushy(),
	}
}

func getLogger() logrus.Logger {
	logger := logrus.New()
	level, err := logrus.ParseLevel(config.GetLogLevel())
	if err != nil {
		panic(err)
	}
	logger.SetLevel(level)
	return *logger
}

func getMongo() helpers.IMongoClient {
	mongo, err := helpers.GetMongo(config.GetMongoConfig())
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

func getPushy() pushy.Pushy {
	sdk := pushy.Create(config.GetPushyToken(), pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(5 * time.Second))
	return *sdk
}

func (s *Server) cleanup() {
	s.Redis.Close()
}
