package main

import (
	"log"
	"potato_killer/config"
	"potato_killer/mcping"
	"potato_killer/puffer_client"
)

func doLoginOrPanic(client *puffer_client.PufferClient, conf config.Config) {
	err := client.Login(conf.PufferPanelUserEmail, conf.PufferPanelUserPasswd)
	if err != nil {
		log.Panicln(err)
	}
}

func isServerPingable(addr string) (bool, *mcping.Status) {
	status, _, err := mcping.PingAndList(addr, 753)

	if err != nil {
		return false, nil
	}

	return true, status
}
