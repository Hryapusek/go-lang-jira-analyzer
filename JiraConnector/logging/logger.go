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
	logger.SetFormatter(&logrus.TextFormatter{})

	logs, _ := os.OpenFile("./log/logs.log", os.O_APPEND, 0666)
	errors, _ := os.OpenFile("./log/err_logs.log", os.O_APPEND, 0666)

	logsFile := io.MultiWriter(os.Stdout, logs)
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
