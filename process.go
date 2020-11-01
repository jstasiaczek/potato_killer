package main

import (
	"errors"
	"io/ioutil"
	"potato_killer/config"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
)

func findServerProcess(conf config.Config) (int, error) {
	processes, _ := ps.Processes()
	for _, process := range processes {
		if process.Executable() == "java" {
			bytes, err := ioutil.ReadFile("/proc/" + strconv.Itoa(process.Pid()) + "/cmdline")
			if err != nil {
				return 0, err
			}
			if strings.Contains(string(bytes), conf.McServerJarFileName) {
				return process.Pid(), nil
			}
		}
	}
	return 0, errors.New("Proccess not found")
}
