package config

import "os"

var Filepath string = os.Getenv("PROMETHEUS_DYN_SD_FILEPATH")

type Hosts []Host

type Host struct {
	Hostname []string          `json:"targets"`
	Labels   map[string]string `json:"labels"`
}
