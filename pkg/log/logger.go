package log

import (
	"github.com/alejoacosta74/gologger"
)

var root gologger.ILogger

func SetLogger(logger *gologger.Logger) {
	root = logger
}

func Debugf(format string, args ...interface{}) {
	root.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	root.Infof(format, args...)
}

func Tracef(format string, args ...interface{}) {
	root.Tracef(format, args...)
}

func With(key string, value interface{}) gologger.ILogger {
	return root.With(key, value)
}

func IsDebug() bool {
	return root.IsDebug()
}
