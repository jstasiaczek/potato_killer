package main

import (
	"log"
	"net"
	"potato_killer/puffer_client"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	conf, confPath := initConfig()
	log.Println("Start with config:", confPath)
	log.Println("Refresh timeout:", conf.RefreshTimeout*time.Second)

	client := puffer_client.NewPufferClient(conf.PufferPanelServerUrl, conf.PufferPanelServerId)
	doLoginOrPanic(client, conf)
	log.Println("Loged to PufferPanel")

	addr := net.JoinHostPort(conf.McServerAddress, conf.McServerPort)

	state := &ServerState{}
	state.LastMcCrashCount, _ = CheckLogs(conf.McCrashReportDir, "")
	state.McCrashCount = state.LastMcCrashCount
	state.LastJavaCrashCount, _ = CheckLogs(conf.JavaCrashReportDir, conf.JavaCrashReportPrefix)
	state.JavaCrashCount = state.LastJavaCrashCount

	for {
		// check if server can be asked for basic data
		isPing, st := isServerPingable(addr)
		state.IsPingable = isPing
		state.UsersCount = 0
		if isPing {
			state.UsersCount = st.Players.Online
		}

		// find server process id
		pid, _ := findServerProcess(conf)
		state.ServerPid = pid

		// check logs count for java and mc
		state.LastJavaCrashCount, _ = CheckLogs(conf.JavaCrashReportDir, conf.JavaCrashReportPrefix)
		state.LastMcCrashCount, _ = CheckLogs(conf.McCrashReportDir, "")

		// relogin to pufferpanel if token is about to expire
		if client.WillTokenExpire() {
			doLoginOrPanic(client, conf)
		}
		//validate state
		validateState(state)
		// detect if crash happened
		detectCrash(state, client)
		// log.Printf("%+v\n", state)
		// wait before next check
		time.Sleep(conf.RefreshTimeout * time.Second)
	}
}

func detectCrash(state *ServerState, client *puffer_client.PufferClient) {
	crashProbability := 0
	// if we have more crash reports than last time
	if state.LastMcCrashCount > state.McCrashCount || state.LastJavaCrashCount > state.JavaCrashCount {
		crashProbability += 2
	}
	// if server is not pingable
	if !state.IsPingable {
		crashProbability++
	}
	// if there is no pid found
	if state.ServerPid == 0 {
		crashProbability++
	}

	//so if we have less than 2 crash symptom, we are not sure if server crasched, lets wait for more
	if crashProbability <= 2 {
		return
	}
	log.Println("[ERROR] Detected server crash! Crash probability score:", crashProbability)
	log.Println("[ERROR] Summary: ")
	log.Println("[ERROR]\tMC crash count: ", state.LastMcCrashCount, ", previous:", state.McCrashCount)
	log.Println("[ERROR]\tJava crash count: ", state.LastJavaCrashCount, ", previous:", state.JavaCrashCount)
	log.Println("[ERROR]\tServer is pingable:", state.IsPingable)
	log.Println("[ERROR]\tServer process id:", state.ServerPid)

	// if pid is bigger than 0, we need to kill server, and wait some time
	if state.ServerPid > 0 {
		err := client.KillServer()
		if err != nil {
			log.Println(err)
		}
		time.Sleep(5 * time.Second)
	}

	// send start server command
	err := client.StartServer()
	log.Println("STARTING SERVER")
	if err != nil {
		log.Println(err)
	}

	// reset log state
	state.McCrashCount = state.LastMcCrashCount
	state.JavaCrashCount = state.LastJavaCrashCount

}
