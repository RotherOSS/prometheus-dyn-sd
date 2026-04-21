package config

import "os"

var Filepath string = os.Getenv("PROMETHEUS_DYN_SD_FILEPATH")

type Hosts struct {
	Hosts []Host `json:"hosts"`
}

type Host struct {
	Hostname string `json:"hostname"`
	Url      string `json:"url"`
}
