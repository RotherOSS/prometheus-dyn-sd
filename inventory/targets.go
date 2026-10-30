package inventory

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"git.otobo.org/rotheross/intern/prometheus-dyn-sd/config"
)

type Config struct {
	Filepath string
	mux      sync.Mutex
}

// This function is called without mux.Lock() because it is only called inside other functions that already implement mutex
func (conf *Config) getTargets() (config.Hosts, error) {
	byteValue, err := os.ReadFile(config.Filepath)
	if err != nil {
		log.Println("Error while reading " + config.Filepath + ". Skipping this entry. Error: \n" + err.Error())
		return config.Hosts{}, err
	}
	var hosts config.Hosts
	err = json.Unmarshal(byteValue, &hosts)
	if err != nil {
		log.Println("Error parsing JSON in " + config.Filepath + ". Error: \n" + err.Error())
	}
	return hosts, err
}

func (conf *Config) GetTarget(hostname string) (config.Host, error) {
	conf.mux.Lock()
	defer conf.mux.Unlock()
	var hosts, err = conf.getTargets()
	if err != nil {
		log.Println("Failed to fetch target because config " + config.Filepath + " could not be parsed. Skip getting and returning targets.")
		return config.Host{}, err
	}
	for _, host := range hosts {
		if len(host.Hostname) > 0 && host.Hostname[0] == hostname {
			return host, err
		}
	}
	return config.Host{}, err
}

func (conf *Config) AddTarget(host config.Host) error {
	conf.mux.Lock()
	defer conf.mux.Unlock()
	hosts, err := conf.getTargets()
	if err != nil {
		log.Println("Failed to add target " + host.Hostname[0] + " because config " + config.Filepath + " could not be parsed. Skip adding target.")
		return err
	}
	hosts = append(hosts, host)
	data, err := json.MarshalIndent(hosts, "", "\t")
	if err != nil {
		log.Println("Failed to marshal host struct to byteCode. Error: \n " + err.Error())
		return err
	}
	err = os.WriteFile(config.Filepath, data, 0644)
	if err != nil {
		log.Println("Failed to write added host to config file. Error: \n " + err.Error())
	}
	return err
}

func (conf *Config) RemoveTarget(hostname string) (error, bool) {
	conf.mux.Lock()
	defer conf.mux.Unlock()
	hosts, err := conf.getTargets()
	var found bool = false
	if err != nil {
		log.Println("Failed to remove target " + hostname + " because config " + config.Filepath + " could not be parsed. Skip removing target.")
		return err, found
	}
	for i, host := range hosts {
		if len(host.Hostname) > 0 && host.Hostname[0] == hostname {
			hosts = append(hosts[:i], hosts[i+1:]...)
			found = true
			break
		}
	}
	data, err := json.MarshalIndent(hosts, "", "\t")
	if err != nil {
		log.Println("Failed to marshal host struct to byteCode. Error: \n " + err.Error())
		return err, found
	}
	err = os.WriteFile(config.Filepath, data, 0644)
	if err != nil {
		log.Println("Failed to write removed host to config file. Error: \n " + err.Error())
	}
	return err, found
}

func (conf *Config) UpdateTarget(target config.Host) error {
	conf.mux.Lock()
	defer conf.mux.Unlock()
	hosts, err := conf.getTargets()
	if err != nil {
		log.Println("Failed to update target " + target.Hostname[0] + " because config " + config.Filepath + " could not be parsed. Skip updating target.")
		return err
	}
	if len(target.Hostname) == 0 {
		return nil
	}
	for i, host := range hosts {
		if len(host.Hostname) > 0 && host.Hostname[0] == target.Hostname[0] {
			hosts[i] = target
		}
	}
	data, err := json.MarshalIndent(hosts, "", "\t")
	if err != nil {
		log.Println("Failed to marshal host struct to byteCode. Error: \n " + err.Error())
		return err
	}
	err = os.WriteFile(config.Filepath, data, 0644)
	if err != nil {
		log.Println("Failed to write updated host to config file. Error: \n " + err.Error())
	}
	return err
}
