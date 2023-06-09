package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type Logger struct {
	logger        *logrus.Logger
	logsFile      *io.Writer
	errorLogsFile *io.Writer
}

func NewLogger() *Logger {
	logger := logrus.New()
	level, _ := logrus.ParseLevel("trace")
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})

	logs, _ := os.OpenFile("../logs/logs.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	errors, _ := os.OpenFile("../logs/err_logs.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	logsFile := io.MultiWriter(logs)
	errorLogsFile := io.MultiWriter(os.Stdout, errors)

	return &Logger{
		logger:        logger,
		logsFile:      &logsFile,
		errorLogsFile: &errorLogsFile,
	}
}

func (l *Logger) Log(logLevel LogLevel, message string) {
	l.logger.Out = *l.logsFile
	switch logLevel {
	case DEBUG:
		l.logger.Debug(message)
	case INFO:
		l.logger.Info(message)
	case WARNING:
		l.logger.Warning(message)
		l.logger.Out = *l.errorLogsFile
		l.logger.Warning(message)
		fmt.Println(message)
	case ERROR:
		l.logger.Error(message)
		l.logger.Out = *l.errorLogsFile
		l.logger.Error(message)
		fmt.Println(message)
	}
}
