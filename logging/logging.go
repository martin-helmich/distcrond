package logging

import (
	"os"
	logging "github.com/op/go-logging"
)

var Logger *logging.Logger

func Setup() {
	format := logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{module} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}")

	defaultBackend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(defaultBackend, format)

	logging.SetBackend(backendFormatter)

	Logger = logging.MustGetLogger("distcrond")
	Logger.Debug("Set up logging")
}

// Critical logs a message using CRITICAL as log level.
func Critical(format string, args ...interface{}) {
	Logger.Critical(format, args...)
}

// Error logs a message using ERROR as log level.
func Error(format string, args ...interface{}) {
	Logger.Error(format, args...)
}

// Warning logs a message using WARNING as log level.
func Warning(format string, args ...interface{}) {
	Logger.Warning(format, args...)
}

// Notice logs a message using NOTICE as log level.
func Notice(format string, args ...interface{}) {
	Logger.Notice(format, args...)
}

// Info logs a message using INFO as log level.
func Info(format string, args ...interface{}) {
	Logger.Info(format, args...)
}

// Debug logs a message using DEBUG as log level.
func Debug(format string, args ...interface{}) {
	Logger.Debug(format, args...)
}
