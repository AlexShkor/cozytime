package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	ConnectionString string
	DatabaseName     string
	CollectionName   string
	Port             string
	TwilioSID        string
	TwilioToken      string
}

var current *Config

func readConfig(path string) (*Config, error) {
	configJson, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config: %s", err)
	}

	var conf Config
	err = json.Unmarshal(configJson, &conf)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse config: %s", err)
	}
	return &conf, nil
}

//Get returns current config
func Get() *Config {
	if current == nil {
		configPath := "config.json"
		conf, err := readConfig(configPath)
		current = conf
		if err != nil {
			fmt.Printf("Failed to read config: %s\n", err)
		}
		fmt.Println(current)
	}
	return current
}
