package main

import (
	"github.com/koding/multiconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel                   string `default:"info"`
	MetricsPort                string `default:":9876"`
	MetricsPath                string `default:"/metrics"`
	OpenhabEndpoint            string `required:"true"`
	OpenhabPollIntervalSeconds int    `default:"30"`
	OpenhabClientTimoutSeconds int    `default:"10"`
}

func loadConfig() *Config {
	logrus.Info("Loading Config")
	cnf := &Config{}
	mc := multiconfig.New()
	mc.MustLoad(cnf)
	logrus.Info("Config loaded successfully")
	return cnf
}
