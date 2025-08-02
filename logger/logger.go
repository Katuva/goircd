package logger

import (
	"fmt"
	"goircd/utils"
	"log"
	"os"
	"strings"
	"sync"

	"goircd/config"
)

var (
	defaultLogger *Logger
	once          sync.Once
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelFatal
)

type Logger struct {
	logger     *log.Logger
	level      LogLevel
	mu         sync.Mutex
	useColours bool
}

func getDefaultLogger() *Logger {
	once.Do(func() {
		defaultLogger = NewLogger()
	})
	return defaultLogger
}

func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warning", "warn":
		return LogLevelWarning
	case "error":
		return LogLevelError
	case "fatal":
		return LogLevelFatal
	default:
		return LogLevelInfo
	}
}

func (l *Logger) getColoredLevelString(level LogLevel) string {
	if !l.useColours {
		switch level {
		case LogLevelDebug:
			return "DEBUG"
		case LogLevelInfo:
			return "INFO"
		case LogLevelWarning:
			return "WARNING"
		case LogLevelError:
			return "ERROR"
		case LogLevelFatal:
			return "FATAL"
		}
		return "UNKNOWN"
	}

	switch level {
	case LogLevelDebug:
		return utils.ColourCyan + "DEBUG" + utils.ColourReset
	case LogLevelInfo:
		return utils.ColourGreen + "INFO" + utils.ColourReset
	case LogLevelWarning:
		return utils.ColourYellow + "WARNING" + utils.ColourReset
	case LogLevelError:
		return utils.ColourRed + "ERROR" + utils.ColourReset
	case LogLevelFatal:
		return utils.BgRed + utils.ColourWhite + utils.ColourBold + "FATAL" + utils.ColourReset
	}
	return "UNKNOWN"
}

func isTerminal(file *os.File) bool {
	fileInfo, _ := file.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func NewLogger() *Logger {
	cfg := config.Get()
	level := parseLogLevel(cfg.Logging.Level)

	logger := &Logger{
		level: level,
	}

	if cfg.Logging.File != "" && cfg.Logging.Console {
		logger.logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
		logger.useColours = isTerminal(os.Stdout)
	} else if cfg.Logging.File != "" {
		file, err := os.OpenFile(cfg.Logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
			logger.useColours = isTerminal(os.Stdout)
		} else {
			logger.logger = log.New(file, "", log.Ldate|log.Ltime)
			logger.useColours = false
		}
	} else {
		logger.logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
		logger.useColours = isTerminal(os.Stdout)
	}

	return logger
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	levelStr := l.getColoredLevelString(level)
	message := fmt.Sprintf(format, args...)

	if level >= LogLevelFatal {
		l.logger.Fatalf("[%s] %s", levelStr, message)
	} else {
		l.logger.Printf("[%s] %s", levelStr, message)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(LogLevelWarning, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LogLevelError, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(LogLevelFatal, format, args...)
}

func (l *Logger) Command(client, command, params string) {
	l.Info("COMMAND: %s executed %s", client, command)
}

func (l *Logger) ChannelOp(client, action, target, channel string) {
	l.Info("CHANOP: %s %s %s in %s", client, action, target, channel)
}

func (l *Logger) IRCOp(client, action, target, info string) {
	l.Info("IRCOP: %s %s %s (%s)", client, action, target, info)
}

func (l *Logger) Ban(client, mask, channel, reason string) {
	l.Info("BAN: %s set ban on %s in %s (%s)", client, mask, channel, reason)
}

func (l *Logger) Kick(client, target, channel, reason string) {
	l.Info("KICK: %s kicked %s from %s (%s)", client, target, channel, reason)
}

func (l *Logger) Connection(client, event string) {
	l.Info("CONNECTION: %s %s", client, event)
}

func Debug(format string, args ...interface{}) {
	getDefaultLogger().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	getDefaultLogger().Info(format, args...)
}

func Warning(format string, args ...interface{}) {
	getDefaultLogger().Warning(format, args...)
}

func Error(format string, args ...interface{}) {
	getDefaultLogger().Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	getDefaultLogger().Fatal(format, args...)
}

func Command(client, command, params string) {
	getDefaultLogger().Command(client, command, params)
}

func ChannelOp(client, action, target, channel string) {
	getDefaultLogger().ChannelOp(client, action, target, channel)
}

func IRCOp(client, action, target, channel string) {
	getDefaultLogger().IRCOp(client, action, target, channel)
}

func Ban(client, mask, channel, reason string) {
	getDefaultLogger().Ban(client, mask, channel, reason)
}

func Kick(client, target, channel, reason string) {
	getDefaultLogger().Kick(client, target, channel, reason)
}

func Connection(client, event string) {
	getDefaultLogger().Connection(client, event)
}
