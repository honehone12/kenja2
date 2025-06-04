package logs

import "github.com/sirupsen/logrus"

var __LOGGER *logrus.Logger

func FmtDefault() {
	__LOGGER = logrus.New()
	__LOGGER.Formatter = new(logrus.JSONFormatter)
}

func FmtPlain() {
	__LOGGER = logrus.New()
	__LOGGER.Formatter = new(logrus.TextFormatter)
	__LOGGER.Formatter.(*logrus.TextFormatter).FullTimestamp = true
}

func SetLevel(level logrus.Level) {
	__LOGGER.Level = level
}

func lazyDefault() {
	if __LOGGER == nil {
		FmtDefault()
	}
}

func Tracef(format string, args ...any) {
	lazyDefault()
	__LOGGER.Tracef(format, args...)
}

func Trace(args ...any) {
	lazyDefault()
	__LOGGER.Trace(args...)
}

func Debugf(format string, args ...any) {
	lazyDefault()
	__LOGGER.Debugf(format, args...)
}

func Debug(args ...any) {
	lazyDefault()
	__LOGGER.Debug(args...)
}

func Infof(format string, args ...any) {
	lazyDefault()
	__LOGGER.Infof(format, args...)
}

func Info(args ...any) {
	lazyDefault()
	__LOGGER.Info(args...)
}

func Warnf(format string, args ...any) {
	lazyDefault()
	__LOGGER.Warnf(format, args...)
}

func Warn(args ...any) {
	lazyDefault()
	__LOGGER.Warn(args...)
}

func Errorf(format string, args ...any) {
	lazyDefault()
	__LOGGER.Errorf(format, args...)
}

func Error(args ...any) {
	lazyDefault()
	__LOGGER.Error(args...)
}

func Panicf(format string, args ...any) {
	lazyDefault()
	__LOGGER.Panicf(format, args...)
}

func Panic(args ...any) {
	lazyDefault()
	__LOGGER.Panic(args...)
}

func Fatalf(format string, args ...any) {
	lazyDefault()
	__LOGGER.Fatalf(format, args...)
}

func Fatal(args ...any) {
	lazyDefault()
	__LOGGER.Fatal(args...)
}
