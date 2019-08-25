//Gist from https://gist.github.com/gonzaloserrano/d360c1b195f7313b6d19

package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// Fields wraps logrus.Fields, which is a map[string]interface{}
type Fields logrus.Fields

// SetLogLevel sets the log lever for this logger
func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

// SetLogFormatter sets the format for the logs
func SetLogFormatter(formatter logrus.Formatter) {
	logger.Formatter = formatter
}

// Debug Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Debug(args...)
	}
}

// DebugWithFields Debug logs a message with fields at level Debug on the standard logger.
func DebugWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Debug(l)
	}
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Info(args...)
	}
}

// InfoWithFields Info logs a message with fields at level Info on the standard logger.
func InfoWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Info(l)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Warn(args...)
	}
}

// WarnWithFields Warn logs a message with fields at level Warn on the standard logger.
func WarnWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Warn(l)
	}
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Error(args...)
	}
}

// ErrorWithFields Error logs a message with fields at level Error on the standard logger.
func ErrorWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Error(l)
	}
}

// Fatal Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(args...)
	}
}

// FatalWithFields Fatal logs a message with fields at level Fatal on the standard logger.
func FatalWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(l)
	}
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	if logger.Level >= logrus.PanicLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Panic(args...)
	}
}

// PanicWithFields Panic logs a message with fields at level Panic on the standard logger.
func PanicWithFields(l interface{}, f Fields) {
	if logger.Level >= logrus.PanicLevel {
		entry := logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Panic(l)
	}
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
