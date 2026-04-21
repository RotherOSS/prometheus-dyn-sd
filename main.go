package main

import (
	"errors"
	"log"
	"os"

	"git.otobo.org/rotheross/intern/prometheus-dyn-sd/api"
	"git.otobo.org/rotheross/intern/prometheus-dyn-sd/config"
)

func checkConfigFile() {
	var path = config.Filepath
	_, err := os.Stat(path)
	if err == nil {
		log.Println("Config file " + path + " exists. Skipping file creation.")
	} else if errors.Is(err, os.ErrNotExist) {
		os.Create(path)
		log.Println("Config file " + path + " has been created successfully.")
	} else if errors.Is(err, os.ErrPermission) {
		log.Fatalln("Insuficient permissions to inspect " + path + ". Canceling execution.")
		return
	} else {
		log.Fatalln("Error while trying to read " + path + ". Canceling execution.")
		return
	}
}

func main() {
	checkConfigFile()
	api.Initialize()
}
