package logger

import "github.com/sirupsen/logrus"

type log struct {
}

var Log *log

func InitLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		ForceColors:   true,
	})
	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)
}

func (l *log) Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func (l *log) Info(args ...interface{}) {
	logrus.Info(args...)
}

func (l *log) Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func (l *log) Warn(args ...interface{}) {
	logrus.Warning(args...)
}

func (l *log) Warnf(format string, args ...interface{}) {
	logrus.Warningf(format, args...)
}

func (l *log) Error(args ...interface{}) {
	logrus.Error(args...)
}

func (l *log) Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func (l *log) Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func (l *log) Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
