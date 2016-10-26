package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/af83/edwig/logger"

	"gopkg.in/yaml.v2"
)

var Config = struct {
	LogStash string
	Syslog   bool

	DB struct {
		Name     string
		User     string
		Password string
		Port     uint
	}
}{}

func LoadConfig(path string) error {
	// Default values
	Config.Syslog = false
	// ...

	env := os.Getenv("EDWIG_ENV")
	if env == "" {
		logger.Log.Debugf("EDWIG_ENV not set")
		env = "development"
	}
	logger.Log.Debugf("Loading %s configuration", env)

	// general config
	data, err := getConfigFileContent(path, "config.yml")
	if err != nil {
		return err
	}
	yaml.Unmarshal(data, &Config)
	// database
	data, err = getConfigFileContent(path, "database.yml")
	if err != nil {
		return err
	}
	yaml.Unmarshal(data, &Config.DB)
	// environement
	data, err = getConfigFileContent(path, fmt.Sprintf("%s.yml", env))
	if err != nil {
		return err
	}
	yaml.Unmarshal(data, &Config)

	return nil
}

func getConfigFileContent(path, file string) ([]byte, error) {
	// Check file at location
	filePath := strings.Join([]string{path, file}, "/")
	if _, err := os.Stat(filePath); err == nil {
		logger.Log.Debugf("Found %s config file at %s", file, filePath)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	// Check file in default location
	gopath := os.Getenv("GOPATH")
	filePath = strings.Join([]string{gopath, "src/github.com/af83/edwig/config", file}, "/")
	if _, err := os.Stat(filePath); err == nil {
		logger.Log.Debugf("Found %s config file at %s", file, filePath)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	logger.Log.Debugf("Can't find %s config file", file)
	return nil, fmt.Errorf("Can't find %s configuration file", file)
}
