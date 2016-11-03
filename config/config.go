package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"

	"gopkg.in/yaml.v2"
)

var Config = struct {
	DB     model.DatabaseConfig
	TestDB model.DatabaseConfig

	LogStash string
	Syslog   bool
}{}

func LoadConfig(path string) error {
	// Default values
	Config.Syslog = false
	// ...

	env := os.Getenv("EDWIG_ENV")
	if env == "" {
		logger.Log.Debugf("EDWIG_ENV not set, default environment is development")
		env = "development"
	}
	logger.Log.Debugf("Loading %s configuration", env)

	configPath, err := getConfigDirectory(path)
	if err != nil {
		return err
	}
	// general config
	files := []string{"config.yml", "database.yml", fmt.Sprintf("%s.yml", env)}
	for _, file := range files {
		data, err := getConfigFileContent(configPath, file)
		if err != nil {
			return err
		}
		yaml.Unmarshal(data, &Config)
	}

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
