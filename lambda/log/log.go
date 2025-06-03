package log

import "github.com/sirupsen/logrus"

var l *logrus.Logger

func FmtDefault() {
	l = logrus.New()
	l.Formatter = new(logrus.JSONFormatter)
}

func FmtPlain() {
	l = logrus.New()
	l.Formatter = new(logrus.TextFormatter)
	l.Formatter.(*logrus.TextFormatter).FullTimestamp = true
}

func SetLevel(level logrus.Level) {
	l.Level = level
}

func lazyDefault() {
	if l == nil {
		FmtDefault()
	}
}

func Tracef(format string, args ...any) {
	lazyDefault()
	l.Tracef(format, args...)
}

func Trace(args ...any) {
	lazyDefault()
	l.Trace(args...)
}

func Debugf(format string, args ...any) {
	lazyDefault()
	l.Debugf(format, args...)
}

func Debug(args ...any) {
	lazyDefault()
	l.Debug(args...)
}

func Infof(format string, args ...any) {
	lazyDefault()
	l.Infof(format, args...)
}

func Info(args ...any) {
	lazyDefault()
	l.Info(args...)
}

func Warnf(format string, args ...any) {
	lazyDefault()
	l.Warnf(format, args...)
}

func Warn(args ...any) {
	lazyDefault()
	l.Warn(args...)
}

func Errorf(format string, args ...any) {
	lazyDefault()
	l.Errorf(format, args...)
}

func Error(args ...any) {
	lazyDefault()
	l.Error(args...)
}

func Panicf(format string, args ...any) {
	lazyDefault()
	l.Panicf(format, args...)
}

func Panic(args ...any) {
	lazyDefault()
	l.Panic(args...)
}

func Fatalf(format string, args ...any) {
	lazyDefault()
	l.Fatalf(format, args...)
}

func Fatal(args ...any) {
	lazyDefault()
	l.Fatal(args...)
}
