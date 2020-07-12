package logs

import (
	log "github.com/sirupsen/logrus"
	"eventspal/pkg/config"
)

func Init(cfg config.Config) {
	debugLevel, err := log.ParseLevel(cfg.App.DebugLevel)
	if err != nil {
		log.Warnf("Unknown debug level %s defaulting to warning level", cfg.App.DebugLevel)
		debugLevel = log.WarnLevel
	}
	log.SetLevel(debugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
}
