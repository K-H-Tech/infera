package logger

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

type log struct {
	logger *slog.Logger
}

var (
	Log  *log
	once sync.Once
)

func InitLogger() {
	once.Do(func() {
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		handler := slog.NewTextHandler(os.Stdout, opts)
		Log = &log{
			logger: slog.New(handler),
		}
	})
}

func (l *log) Debug(args ...any) {
	if l == nil || l.logger == nil {
		return
	}
	if len(args) == 1 {
		l.logger.Debug(fmt.Sprint(args[0]))
	} else {
		l.logger.Debug(fmt.Sprint(args...))
	}
}

func (l *log) Debugf(format string, args ...any) {
	if l == nil || l.logger == nil {
		return
	}
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *log) Info(args ...any) {
	if l == nil || l.logger == nil {
		return
	}
	if len(args) == 1 {
		l.logger.Info(fmt.Sprint(args[0]))
	} else {
		l.logger.Info(fmt.Sprint(args...))
	}
}

func (l *log) Infof(format string, args ...any) {
	if l == nil || l.logger == nil {
		return
	}
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *log) Warn(args ...any) {
	if l == nil || l.logger == nil {
		return
	}
	if len(args) == 1 {
		l.logger.Warn(fmt.Sprint(args[0]))
	} else {
		l.logger.Warn(fmt.Sprint(args...))
	}
}

func (l *log) Warnf(format string, args ...any) {
	if l == nil || l.logger == nil {
		return
	}
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *log) Error(args ...any) {
	if l == nil || l.logger == nil {
		return
	}
	if len(args) == 1 {
		l.logger.Error(fmt.Sprint(args[0]))
	} else {
		l.logger.Error(fmt.Sprint(args...))
	}
}

func (l *log) Errorf(format string, args ...any) {
	if l == nil || l.logger == nil {
		return
	}
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *log) Fatal(args ...any) {
	if l == nil || l.logger == nil {
		os.Exit(1)
	}
	if len(args) == 1 {
		l.logger.Error(fmt.Sprint(args[0]))
	} else {
		l.logger.Error(fmt.Sprint(args...))
	}
	os.Exit(1)
}

func (l *log) Fatalf(format string, args ...any) {
	if l == nil || l.logger == nil {
		os.Exit(1)
	}
	l.logger.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}
