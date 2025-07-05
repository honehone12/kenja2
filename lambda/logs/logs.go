package logs

import (
	"kenja2/engine"

	"github.com/sirupsen/logrus"
)

var __logger *logrus.Logger

func FmtDefault() {
	__logger = logrus.New()
	__logger.Formatter = new(logrus.JSONFormatter)
}

func FmtPlain() {
	__logger = logrus.New()
	__logger.Formatter = new(logrus.TextFormatter)
	__logger.Formatter.(*logrus.TextFormatter).FullTimestamp = true
}

func SetLevel(level Level) {
	__logger.SetLevel(level)
}

func Logger() engine.Logger {
	return __logger
}

func lazyDefault() {
	if __logger == nil {
		FmtDefault()
	}
}

func Tracef(format string, args ...any) {
	lazyDefault()
	__logger.Tracef(format, args...)
}

func Trace(args ...any) {
	lazyDefault()
	__logger.Trace(args...)
}

func Debugf(format string, args ...any) {
	lazyDefault()
	__logger.Debugf(format, args...)
}

func Debug(args ...any) {
	lazyDefault()
	__logger.Debug(args...)
}

func Infof(format string, args ...any) {
	lazyDefault()
	__logger.Infof(format, args...)
}

func Info(args ...any) {
	lazyDefault()
	__logger.Info(args...)
}

func Warnf(format string, args ...any) {
	lazyDefault()
	__logger.Warnf(format, args...)
}

func Warn(args ...any) {
	lazyDefault()
	__logger.Warn(args...)
}

func Errorf(format string, args ...any) {
	lazyDefault()
	__logger.Errorf(format, args...)
}

func Error(args ...any) {
	lazyDefault()
	__logger.Error(args...)
}

func Panicf(format string, args ...any) {
	lazyDefault()
	__logger.Panicf(format, args...)
}

func Panic(args ...any) {
	lazyDefault()
	__logger.Panic(args...)
}

func Fatalf(format string, args ...any) {
	lazyDefault()
	__logger.Fatalf(format, args...)
}

func Fatal(args ...any) {
	lazyDefault()
	__logger.Fatal(args...)
}
