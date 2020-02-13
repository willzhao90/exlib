package logging

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func SetupLogging(v *viper.Viper) {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
}

func Logger(event, topic, key, code string) *log.Entry {
	return log.WithFields(log.Fields{
		"event": event,
		"topic": topic,
		"key":   key,
		"code":  code,
	})
}
