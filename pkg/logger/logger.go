package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// InitLogger initializes the logger
func InitLogger() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// Get the file name and line number
			filename := filepath.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

	// Set output to stdout by default
	log.SetOutput(os.Stdout)

	// Include the calling method as a field
	log.SetReportCaller(true)
}

// SetOutput sets the output destination for the logger
func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Debugf logs a debug message with format
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	log.Info(args...)
}

// Infof logs an info message with format
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Warnf logs a warning message with format
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	log.Error(args...)
}

// Errorf logs an error message with format
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatal logs a fatal message and exits the application
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Fatalf logs a fatal message with format and exits the application
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// WithField adds a single field to the log entry
func WithField(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, value)
}

// WithFields adds multiple fields to the log entry
func WithFields(fields map[string]interface{}) *logrus.Entry {
	return log.WithFields(fields)
}

// WithError adds an error to the log entry
func WithError(err error) *logrus.Entry {
	return log.WithError(err)
}

// GetLogger returns the underlying logger instance
func GetLogger() *logrus.Logger {
	return log
}

// LogRequest logs an HTTP request
func LogRequest(method, path string, status int, latency time.Duration, clientIP string) {
	entry := log.WithFields(logrus.Fields{
		"method":   method,
		"path":     path,
		"status":   status,
		"latency":  latency,
		"clientIP": clientIP,
	})

	if status >= 500 {
		entry.Error("Server error")
	} else if status >= 400 {
		entry.Warn("Client error")
	} else {
		entry.Info("Request processed")
	}
}

// LogSQLError logs an SQL error with context
func LogSQLError(err error, query string, args ...interface{}) {
	// Limit the length of the query and args in logs
	safeQuery := query
	if len(safeQuery) > 1000 {
		safeQuery = safeQuery[:1000] + "..."
	}

	// Convert args to string representation
	argsStr := make([]string, 0, len(args))
	for _, arg := range args {
		argsStr = append(argsStr, fmt.Sprintf("%v", arg))
	}

	log.WithFields(logrus.Fields{
		"query": safeQuery,
		"args":  argsStr,
	}).Errorf("SQL error: %v", err)
}
