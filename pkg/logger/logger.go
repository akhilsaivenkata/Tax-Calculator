package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log is the shared logger instance used across the application.
var Log *logrus.Logger

// Init initializes the global Log instance with structured JSON logging.
func Init() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{}) // Structured JSON logs
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
}
