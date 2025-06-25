package logs

import "github.com/sirupsen/logrus"

type Level = logrus.Level

const (
	PANIC Level = logrus.PanicLevel
	FATAL Level = logrus.FatalLevel
	ERROR Level = logrus.ErrorLevel
	WARN  Level = logrus.WarnLevel
	INFO  Level = logrus.InfoLevel
	DEBUG Level = logrus.DebugLevel
	TRACE Level = logrus.TraceLevel
)
