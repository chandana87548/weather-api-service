package util

import (
	"fmt"
	"github.com/go-akka/configuration"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func InitApp(configFileLocation string) *configuration.Config {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Info("current path is " + path)
	configFile := fmt.Sprintf("application.conf")
	conf := configuration.LoadConfig(configFileLocation + configFile)
	initLogging(conf)
	log.Infof("successfully loaded config [%s]", configFile)
	log.Debugf("loaded conf %s", conf)

	return conf
}

func initLogging(config *configuration.Config) {
	logLevel := config.GetString("log-level")
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Warnf("failed to parse log level: %v", err)
		level = log.InfoLevel
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(level)
	log.Infof("log level set to [%s]", logLevel)
}

func GetCurrentTime(m time.Duration) string {
	loc, _ := time.LoadLocation("UTC")
	t := time.Now().In(loc).Add(m * time.Second)
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return formatted
}
