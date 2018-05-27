package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/fossapps/starter/config"
	"github.com/multiplay/go-slack/lrhook"
	"github.com/multiplay/go-slack/chat"
)

// Client needs to be implemented for a logger to be used on this project
type Client interface {
	Info(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Print(args ...interface{})
	Panic(args ...interface{})
}

// GetClient returns implementation of Client interface
func GetClient(level string) Client {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(level)
	logger.AddHook(getSlackHook())
	if err != nil {
		panic(err)
	}
	logger.SetLevel(lvl)
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
