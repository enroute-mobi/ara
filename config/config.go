package config

import (
	"errors"
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

	configPath, err := getConfigDirectory(path)
	if err != nil {
		return err
	}
	// general config
	data, err := getConfigFileContent(configPath, "config.yml")
	if err != nil {
		return err
	}
	yaml.Unmarshal(data, &Config)
	// database
	data, err = getConfigFileContent(configPath, "database.yml")
	if err != nil {
		return err
	}
	yaml.Unmarshal(data, &Config.DB)
	// environement
	data, err = getConfigFileContent(configPath, fmt.Sprintf("%s.yml", env))
	if err != nil {
		return err
	}
	yaml.Unmarshal(data, &Config)

	return nil
}

func getConfigDirectory(path string) (string, error) {
	paths := [3]string{
		path,
		os.Getenv("EDWIG_CONFIG"),
		fmt.Sprintf("%s/src/github.com/af83/edwig/config", os.Getenv("GOPATH")),
	}
	for _, directoryPath := range paths {
		if found := checkDirectory(directoryPath); found {
			return directoryPath, nil
		}
	}
	return "", errors.New("Cant find config directory")
}

func checkDirectory(path string) bool {
	if path == "" {
		return false
	}
	if _, err := os.Stat(path); err == nil {
		logger.Log.Debugf("Found config directory at %s", path)
		return true
	}
	logger.Log.Debugf("Can't find config directory at %s", path)
	return false
}

func getConfigFileContent(path, file string) ([]byte, error) {
	// Check file at location
	filePath := strings.Join([]string{path, file}, "/")
	if _, err := os.Stat(filePath); err == nil {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	logger.Log.Debugf("Can't find %s config file in %s", file, path)
	return nil, fmt.Errorf("Can't find %s configuration file", file)
}
