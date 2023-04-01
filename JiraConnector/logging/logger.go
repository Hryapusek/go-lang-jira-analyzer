package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	logger        *logrus.Logger
	logsFile      *os.File
	errorLogsFile *os.File
}

func NewLogger() *Logger {
	logger := logrus.New()
	level, _ := logrus.ParseLevel("trace")
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})

	logsFile, _ := os.OpenFile("./logs/logs.log", os.O_APPEND, 0666)
	errorLogsFile, _ := os.OpenFile("./logs/err_logs.log", os.O_APPEND, 0666)

	return &Logger{
		logger:        logger,
		logsFile:      logsFile,
		errorLogsFile: errorLogsFile,
	}
}

func (l *Logger) Log(logLevel LogLevel, message string) {
	l.logger.Out = l.logsFile
	switch logLevel {
	case DEBUG:
		l.logger.Debug(message)
	case INFO:
		l.logger.Info(message)
	case WARNING:
		l.logger.Warning(message)
		l.logger.Out = l.errorLogsFile
		l.logger.Warning(message)
		fmt.Println(message)
	case ERROR:
		l.logger.Error(message)
		l.logger.Out = l.errorLogsFile
		l.logger.Error(message)
		fmt.Println(message)
	}
}
