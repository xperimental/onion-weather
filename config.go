package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/spf13/pflag"
)

// Config contains the application configuration.
type Config struct {
	Netatmo        netatmo.Config `json:"netatmo"`
	UpdateInterval time.Duration  `json:"updateInterval"`
	Dummy          bool           `json:"useDummy"`
	ConfigFile     string         `json:"-"`
}

const (
	minimumUpdateInterval = 30 * time.Second
)

func parseConfig() (config Config, err error) {
	pflag.StringVarP(&config.Netatmo.ClientID, "client-id", "i", "", "Client ID for NetAtmo app.")
	pflag.StringVarP(&config.Netatmo.ClientSecret, "client-secret", "s", "", "Client secret for NetAtmo app.")
	pflag.StringVarP(&config.Netatmo.Username, "username", "u", "", "Username of NetAtmo account.")
	pflag.StringVarP(&config.Netatmo.Password, "password", "p", "", "Password of NetAtmo account.")
	pflag.DurationVarP(&config.UpdateInterval, "interval", "t", 15*time.Minute, "Interval between updates.")
	pflag.BoolVar(&config.Dummy, "dummy", false, "Use dummy display (output to STDOUT).")
	pflag.StringVarP(&config.ConfigFile, "config", "c", "", "Path to configuration file.")
	pflag.Parse()

	if len(config.ConfigFile) > 0 {
		config, err = readConfig(config.ConfigFile)
		if err != nil {
			return config, err
		}
	}

	if len(config.Netatmo.ClientID) == 0 {
		return config, errors.New("need a NetAtmo client ID")
	}

	if len(config.Netatmo.ClientSecret) == 0 {
		return config, errors.New("need a NetAtmo client secret")
	}

	if len(config.Netatmo.Username) == 0 {
		return config, errors.New("username can not be blank")
	}

	if len(config.Netatmo.Password) == 0 {
		return config, errors.New("password can not be blank")
	}

	if config.UpdateInterval < minimumUpdateInterval {
		return config, fmt.Errorf("update interval too short (min %s): %s", minimumUpdateInterval, config.UpdateInterval)
	}

	return config, nil
}

func readConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
