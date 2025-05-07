package logging

import (
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"log"
)

var Logger = logrus.New()

type CustomFormatter struct {
	logrus.TextFormatter
}

// InitializeLogger configura o logger para usar o New Relic e o logrus
func InitializeLogger(newRelicApp *newrelic.Application) {
	if newRelicApp != nil {
		nrlogrusFormatter := nrlogrus.NewFormatter(newRelicApp, &logrus.TextFormatter{
			TimestampFormat: "2006/01/02 15:04:05.000",
		})
		Logger.SetFormatter(nrlogrusFormatter)
		Logger.SetLevel(logrus.InfoLevel)

		log.SetOutput(Logger.Writer())
		log.SetFlags(0)
	}
}

func SetFields(fields logrus.Fields) {
	Logger.WithFields(fields)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}

func Infoln(args ...interface{}) {
	Logger.Infoln(args...)
}

func Errorln(args ...interface{}) {
	Logger.Errorln(args...)
}

func Warnln(args ...interface{}) {
	Logger.Warnln(args...)
}

func Debugln(args ...interface{}) {
	Logger.Debugln(args...)
}

func Fatalln(args ...interface{}) {
	Logger.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	Logger.Panicln(args...)
}
