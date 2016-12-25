package main

import (
	"encoding/json"
	"os"
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/spf13/pflag"
)

type Config struct {
	Netatmo        netatmo.Config `json:"netatmo"`
	UpdateInterval time.Duration  `json:"updateInterval"`
	Dummy          bool           `json:"useDummy"`
	ConfigFile     string         `json:"-"`
}

func parseConfig() (Config, error) {
	var config Config
	pflag.StringVarP(&config.Netatmo.ClientID, "client-id", "i", "", "Client ID for NetAtmo app.")
	pflag.StringVarP(&config.Netatmo.ClientSecret, "client-secret", "s", "", "Client secret for NetAtmo app.")
	pflag.StringVarP(&config.Netatmo.Username, "username", "u", "", "Username of NetAtmo account.")
	pflag.StringVarP(&config.Netatmo.Password, "password", "p", "", "Password of NetAtmo account.")
	pflag.DurationVarP(&config.UpdateInterval, "interval", "t", 15*time.Minute, "Interval between updates.")
	pflag.BoolVar(&config.Dummy, "dummy", false, "Use dummy display (output to STDOUT).")
	pflag.StringVarP(&config.ConfigFile, "config", "c", "", "Path to configuration file.")
	pflag.Parse()

	if len(config.ConfigFile) > 0 {
		return readConfig(config.ConfigFile)
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
