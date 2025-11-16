package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type log struct {
	logger *slog.Logger
}

var Log *log

func InitLogger() {
	// Create a new slog logger with TextHandler that mimics the logrus formatting
	// logrus.SetFormatter(&logrus.TextFormatter{
	//     DisableColors: false,
	//     FullTimestamp: true,
	//     ForceColors:   true,
	// })
	// logrus.SetLevel(logrus.DebugLevel)

	// Configure slog with similar settings to logrus
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	Log = &log{
		logger: slog.New(handler),
	}
}

func (l *log) Debug(args ...interface{}) {
	if len(args) == 1 {
		l.logger.Debug(fmt.Sprint(args[0]))
	} else {
		l.logger.Debug(fmt.Sprint(args...))
	}
}

func (l *log) Info(args ...interface{}) {
	if len(args) == 1 {
		l.logger.Info(fmt.Sprint(args[0]))
	} else {
		l.logger.Info(fmt.Sprint(args...))
	}
}

func (l *log) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *log) Warn(args ...interface{}) {
	if len(args) == 1 {
		l.logger.Warn(fmt.Sprint(args[0]))
	} else {
		l.logger.Warn(fmt.Sprint(args...))
	}
}

func (l *log) Warnf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *log) Error(args ...interface{}) {
	if len(args) == 1 {
		l.logger.Error(fmt.Sprint(args[0]))
	} else {
		l.logger.Error(fmt.Sprint(args...))
	}
}

func (l *log) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *log) Fatal(args ...interface{}) {
	if len(args) == 1 {
		l.logger.Error(fmt.Sprint(args[0]))
	} else {
		l.logger.Error(fmt.Sprint(args...))
	}
	os.Exit(1)
}

func (l *log) Fatalf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}
