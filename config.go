package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	ConnectionString string `json:"ConnectionString"`
	DatabaseName     string `json:"DatabaseName"`
	CollectionName   string `json:"CollectionName"`
}

func ReadConfig(path string) (*Config, error) {
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
