package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

func Init(out io.Writer, env string) {
	if env == "dev" {
		logrus.SetLevel(logrus.TraceLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	logrus.SetOutput(out)
}

func Trace(msg ...interface{}) {
	logrus.Trace(msg...)
}
func Tracef(format string, msg ...interface{}) {
	logrus.Tracef(format, msg...)
}

func Debug(msg ...interface{}) {
	logrus.Debug(msg...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Info(msg ...interface{}) {
	logrus.Info(msg...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Error(msg ...interface{}) {
	logrus.Error(msg...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Fatal(msg ...interface{}) {
	logrus.Fatal(msg...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
