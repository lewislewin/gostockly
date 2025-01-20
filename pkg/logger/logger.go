package logger

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

var (
	instance *Logger
	once     sync.Once
)

// GetLogger ensures a single instance of Logger.
func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{
			infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
			errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
			debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		}
	})
	return instance
}

// Info logs general information.
func (l *Logger) Info(message string, args ...interface{}) {
	l.infoLogger.Printf(message, args...)
}

// Error logs errors.
func (l *Logger) Error(message string, args ...interface{}) {
	l.errorLogger.Printf(message, args...)
}

// Debug logs debug information
func (l *Logger) Debug(message string, args ...interface{}) {
	l.debugLogger.Printf(message, args...)
}
