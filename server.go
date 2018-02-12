package crazy_nl_backend

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"gopkg.in/matryer/respond.v1"
	"crazy_nl_backend/helpers"
)

type Server struct {
	Logger *logrus.Logger
	Db     *helpers.Mongo
	Redis  *helpers.Redis
}

func Init() {
	helpers.InitDotEnv()
	server := createServer()
	defer server.cleanup()
	router := NewRouter(server)
	server.Logger.Info("Listening on port " + helpers.DotEnv{}.GetServerPort())
	http.ListenAndServe(":" + helpers.DotEnv{}.GetServerPort(), router)
}

func Decode(r *http.Request, data interface{}) {
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&data)
}

func (s *Server) ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	respond.With(w, r, http.StatusBadRequest, struct {
		Message string `json:"message"`
	} {
		Message: message,
	})
	return
}

func createServer() Server{
	db := getMongo()
	return Server{
		Logger: getLogger(),
		Db:     db,
		Redis:  getRedis(),
	}
}

func getLogger() *logrus.Logger {
	logger := logrus.New()
	level, err := logrus.ParseLevel(helpers.DotEnv{}.GetLogLevel())
	if err != nil {
		panic(err)
	}
	logger.SetLevel(level)
	return logger
}

func getMongo() *helpers.Mongo {
	mongo, err := helpers.GetMongo()
	if err != nil {
		panic(err)
	}
	return mongo
}

func getRedis() *helpers.Redis{
	redis, err := helpers.GetRedis()
	if err != nil {
		panic(err)
	}
	return redis
}

func (s *Server) cleanup() {
	s.Db.Close()
	s.Redis.Close()
}
