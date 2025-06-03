package swissknife

import (
	"log"
	"os"
)

// LogLevel represents different logging levels
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger interface for customizable logging

// Default logger implementation

type DefaultLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	level       LogLevel
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		infoLogger:  log.New(os.Stdout, "[TLS-RPC] ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "[TLS-RPC] ", log.LstdFlags),
		level:       LogLevelInfo, // Default to Info level
	}
}

func (dl *DefaultLogger) SetLevel(level LogLevel) {
	dl.level = level
}

func (dl *DefaultLogger) GetLevel() LogLevel {
	return dl.level
}

func (dl *DefaultLogger) shouldLog(level LogLevel) bool {
	return level >= dl.level
}

func (dl *DefaultLogger) Debug(v ...interface{}) {
	if dl.shouldLog(LogLevelDebug) {
		dl.infoLogger.SetPrefix("[TLS-RPC] DEBUG: ")
		dl.infoLogger.Println(v...)
		dl.infoLogger.SetPrefix("[TLS-RPC] ")
	}
}

func (dl *DefaultLogger) Debugf(format string, v ...interface{}) {
	if dl.shouldLog(LogLevelDebug) {
		dl.infoLogger.SetPrefix("[TLS-RPC] DEBUG: ")
		dl.infoLogger.Printf(format, v...)
		dl.infoLogger.SetPrefix("[TLS-RPC] ")
	}
}

func (dl *DefaultLogger) Info(v ...interface{}) {
	if dl.shouldLog(LogLevelInfo) {
		dl.infoLogger.SetPrefix("[TLS-RPC] INFO: ")
		dl.infoLogger.Println(v...)
		dl.infoLogger.SetPrefix("[TLS-RPC] ")
	}
}

func (dl *DefaultLogger) Infof(format string, v ...interface{}) {
	if dl.shouldLog(LogLevelInfo) {
		dl.infoLogger.SetPrefix("[TLS-RPC] INFO: ")
		dl.infoLogger.Printf(format, v...)
		dl.infoLogger.SetPrefix("[TLS-RPC] ")
	}
}

func (dl *DefaultLogger) Warn(v ...interface{}) {
	if dl.shouldLog(LogLevelWarn) {
		dl.errorLogger.SetPrefix("[TLS-RPC] WARN: ")
		dl.errorLogger.Println(v...)
		dl.errorLogger.SetPrefix("[TLS-RPC] ")
	}
}

func (dl *DefaultLogger) Warnf(format string, v ...interface{}) {
	if dl.shouldLog(LogLevelWarn) {
		dl.errorLogger.SetPrefix("[TLS-RPC] WARN: ")
		dl.errorLogger.Printf(format, v...)
		dl.errorLogger.SetPrefix("[TLS-RPC] ")
	}
}

func (dl *DefaultLogger) Error(v ...interface{}) {
	if dl.shouldLog(LogLevelError) {
		dl.errorLogger.SetPrefix("[TLS-RPC] ERROR: ")
		dl.errorLogger.Println(v...)
		dl.errorLogger.SetPrefix("[TLS-RPC] ")
	}
}

func (dl *DefaultLogger) Errorf(format string, v ...interface{}) {
	if dl.shouldLog(LogLevelError) {
		dl.errorLogger.SetPrefix("[TLS-RPC] ERROR: ")
		dl.errorLogger.Printf(format, v...)
		dl.errorLogger.SetPrefix("[TLS-RPC] ")
	}
}

func (dl *DefaultLogger) Fatal(v ...interface{}) {
	if dl.shouldLog(LogLevelFatal) {
		dl.errorLogger.SetPrefix("[TLS-RPC] FATAL: ")
		dl.errorLogger.Println(v...)
		os.Exit(1)
	}
}

func (dl *DefaultLogger) Fatalf(format string, v ...interface{}) {
	if dl.shouldLog(LogLevelFatal) {
		dl.errorLogger.SetPrefix("[TLS-RPC] FATAL: ")
		dl.errorLogger.Printf(format, v...)
		os.Exit(1)
	}
}
