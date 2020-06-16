package client

import (
	"github.com/sirupsen/logrus"
	"os"
	"gopkg.in/yaml.v2"
)

type Config struct {
	IgnoredDirectories []string `yaml:"ignored-directories"`
	WatchedExtensions  []string `yaml:"watched-extensions"`
}

var LoadedInstance Config;

func LoadConfiguration() Config {
	f, err := os.Open("client.yml")
	if err != nil {
		logrus.Error("Could not find client.yml")
		os.Exit(1)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		logrus.Error(err)
	}

	LoadedInstance = cfg

	return cfg
}