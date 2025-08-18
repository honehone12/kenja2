package quicklog

import "log"

type QuickLogger struct{}

func NewQuickLogger() QuickLogger {
	return QuickLogger{}
}

func (l QuickLogger) Debug(args ...any) {
	log.SetPrefix("[DEBUG]")
	log.Println(args...)
}

func (l QuickLogger) Debugf(format string, args ...any) {
	log.SetPrefix("[DEBUG]")
	log.Printf(format, args...)
}

func (l QuickLogger) Info(args ...any) {
	log.SetPrefix("[INFO]")
	log.Println(args...)
}

func (l QuickLogger) Infof(format string, args ...any) {
	log.SetPrefix("[INFO]")
	log.Printf(format, args...)
}

func (l QuickLogger) Warn(args ...any) {
	log.SetPrefix("[WARNING]")
	log.Println(args...)
}

func (l QuickLogger) Warnf(format string, args ...any) {
	log.SetPrefix("[WARNING]")
	log.Printf(format, args...)
}

func (l QuickLogger) Error(args ...any) {
	log.SetPrefix("[ERROR]")
	log.Println(args...)
}

func (l QuickLogger) Errorf(format string, args ...any) {
	log.SetPrefix("[ERROR]")
	log.Printf(format, args...)
}

func (l QuickLogger) Fatal(args ...any) {
	log.SetPrefix("[FATAL]")
	log.Fatalln(args...)
}

func (l QuickLogger) Fatalf(format string, args ...any) {
	log.SetPrefix("[FATAL]")
	log.Fatalf(format, args...)
}

func (l QuickLogger) Panic(args ...any) {
	log.SetPrefix("[PANIC]")
	log.Panicln(args...)
}

func (l QuickLogger) Panicf(format string, args ...any) {
	log.SetPrefix("[PANIC]")
	log.Panicf(format, args...)
}
