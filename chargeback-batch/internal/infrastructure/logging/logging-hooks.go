package logging

import (
	"context"
	"github.com/sirupsen/logrus"
)

type TraceIDHook struct {
	TraceIDKey string
	Context    context.Context
}

func (hook *TraceIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *TraceIDHook) Fire(entry *logrus.Entry) error {
	if traceID, ok := hook.Context.Value(hook.TraceIDKey).(string); ok {
		entry.Data[hook.TraceIDKey] = traceID
	}
	return nil
}

type CountryHook struct {
	CountryKey string
	Context    context.Context
}

func (hook *CountryHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *CountryHook) Fire(entry *logrus.Entry) error {
	if country, ok := hook.Context.Value(hook.CountryKey).(string); ok {
		entry.Data[hook.CountryKey] = country
	}
	return nil
}
