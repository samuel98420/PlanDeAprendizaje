package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"net/http"
)

var Logger *logrus.Logger

func Init() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
		},
	})
	Logger.SetOutput(os.Stdout)
}

func LogRequest(r *http.Request) {
	Logger.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"ip":     r.RemoteAddr,
	}).Info("incoming_request")
}