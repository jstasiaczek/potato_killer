package main

import (
	"flag"
	"potato_killer/config"
)

func initConfig() (config.Config, string) {
	conf := config.Config{}
	configFileLocation := flag.String("config", "config.json", "config file location")
	flag.Parse()
	config.Import(configFileLocation, &conf)
	return conf, *configFileLocation

}
