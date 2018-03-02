package main

import (
	"testing"
	"github.com/Sirupsen/logrus"
	"github.com/Bestfeel/hole/logger"
)

func TestLog(t *testing.T) {

	//logger.AddHook()

	logger.SetFormatter(&logrus.TextFormatter{ForceColors: true, DisableColors: false, DisableTimestamp: false, FullTimestamp: true, DisableSorting: true})

	logger.Info("log info")
	logger.InfoWithFields("log info", logrus.Fields{"name": "hole"})

	logger.Error("log error")
	logger.ErrorWithFields("log error", logrus.Fields{"name": "hole"})

	logger.Debug("log debug")
	logger.DebugWithFields("log debug", logrus.Fields{"name": "hole"})

	logger.Warn("log warn")
	logger.WarnWithFields("log warn", logrus.Fields{"name": "hole"})

}
