package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

type Config struct {
	McServerAddress       string
	McServerPort          string
	JavaCrashReportDir    string
	JavaCrashReportPrefix string
	McCrashReportDir      string
	McServerJarFileName   string
	PufferPanelServerUrl  string
	PufferPanelUserEmail  string
	PufferPanelUserPasswd string
	PufferPanelServerId   string
	RefreshTimeout        time.Duration
	LogLevel              int
}

func Import(configPath *string, config *Config) {
	raw, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(raw, config)
	if err != nil {
		log.Fatalln(err)
	}
}
