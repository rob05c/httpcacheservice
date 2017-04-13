package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	// RFCCompliant determines whether `Cache-Control: no-cache` requests are honored. The ability to ignore `no-cache` is necessary to protect origin servers from DDOS attacks.
	RFCCompliant bool `json:"rfc_compliant"`
	// Port is the HTTP port to serve on
	Port int `json:"port"`
	// CacheSizeBytes is the size of the memory cache, in bytes.
	CacheSizeBytes int `json:"cache_size_bytes"`
}

// DefaultConfig is the default configuration for the application, if no configuration file is given, or if a given config setting doesn't exist in the config file.
var DefaultConfig = Config{
	RFCCompliant:   true,
	Port:           80,
	CacheSizeBytes: bytesPerGibibyte,
}

// Load loads the given config file. If an empty string is passed, the default config is returned.
func LoadConfig(fileName string) (Config, error) {
	cfg := DefaultConfig
	if fileName == "" {
		return cfg, nil
	}
	configBytes, err := ioutil.ReadFile(fileName)
	if err == nil {
		err = json.Unmarshal(configBytes, &cfg)
	}
	return cfg, err
}

const bytesPerGibibyte = 1024 * 1024 * 1024
