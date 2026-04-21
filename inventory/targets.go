package inventory

import (
	"encoding/json"
	"log"
	"os"

	"git.otobo.org/rotheross/intern/prometheus-dyn-sd/config"
)

func getTargets() (config.Hosts, error) {
	byteValue, err := os.ReadFile(config.Filepath)
	if err != nil {
		log.Println("Error while reading " + config.Filepath + ". Skipping this entry. Error: \n" + err.Error())
		return config.Hosts{}, err
	}
	var hosts config.Hosts
	json.Unmarshal(byteValue, &hosts)
	return hosts, err
}

func GetTarget(hostname string) (config.Host, error) {
	var hosts, err = getTargets()
	if err != nil {
		log.Println("Failed to fetch target because config " + config.Filepath + " could not be parsed. Skip getting and returning targets.")
		return config.Host{}, err
	}
	for _, host := range hosts.Hosts {
		if host.Hostname == hostname {
			return host, err
		}
	}
	return config.Host{}, err
}

func AddTarget(host config.Host) error {
	hosts, err := getTargets()
	if err != nil {
		log.Println("Failed to add target " + host.Hostname + " because config " + config.Filepath + " could not be parsed. Skip adding target.")
		return err
	}
	hosts.Hosts = append(hosts.Hosts, host)
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

func RemoveTarget(hostname string) (error, bool) {
	hosts, err := getTargets()
	var found bool = false
	if err != nil {
		log.Println("Failed to remove target " + hostname + " because config " + config.Filepath + " could not be parsed. Skip removing target.")
		return err, found
	}
	for i, host := range hosts.Hosts {
		if host.Hostname == hostname {
			hosts.Hosts = append(hosts.Hosts[:i], hosts.Hosts[i+1:]...)
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

func UpdateTarget(target config.Host) error {
	hosts, err := getTargets()
	if err != nil {
		log.Println("Failed to update target " + target.Hostname + " because config " + config.Filepath + " could not be parsed. Skip updating target.")
		return err
	}
	for i, host := range hosts.Hosts {
		if host.Hostname == target.Hostname {
			hosts.Hosts[i] = target
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
