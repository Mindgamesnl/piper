package server

import (
	"github.com/sirupsen/logrus"
	"os"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port  int `yaml:"port"`
	Password  string `yaml:"password"`
}

var LoadedInstance Config;

func LoadConfiguration() Config {
	if len(os.Args) < 3 {
		logrus.Error("Pleace specify a config file")
		os.Exit(1)
	}
	f, err := os.Open(os.Args[2])
	if err != nil {
		logrus.Error("Could not find config file or no file specified")
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
