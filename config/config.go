package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"bitbucket.org/enroute-mobi/edwig/logger"
	yaml "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     uint
}

var Config = struct {
	DB DatabaseConfig

	ApiKey        string
	Debug         bool
	LogStash      string
	Syslog        bool
	ColorizeLog   bool
	LoadMaxInsert int
}{}

func LoadConfig(path string) error {
	// Default values
	Config.LoadMaxInsert = 100000

	env := Environment()
	logger.Log.Debugf("Loading %s configuration", env)

	configPath, err := getConfigDirectory(path)
	if err != nil {
		return err
	}

	// general config
	files := []string{"config.yml", fmt.Sprintf("%s.yml", env)}
	for _, file := range files {
		data, err := getConfigFileContent(configPath, file)
		if err != nil {
			return err
		}
		yaml.Unmarshal(data, &Config)
	}
	// database config
	LoadDatabaseConfig(configPath)

	logger.Log.Syslog = Config.Syslog
	logger.Log.Debug = Config.Debug
	logger.Log.Color = Config.ColorizeLog

	return nil
}

var environment string

func SetEnvironment(env string) {
	environment = env
}

func Environment() string {
	if environment == "" {
		env := os.Getenv("EDWIG_ENV")
		if env == "" {
			logger.Log.Debugf("EDWIG_ENV not set, default environment is development")
			env = "development"
		}
		environment = env
	}
	return environment
}

func LoadDatabaseConfig(configPath string) error {
	data, err := getConfigFileContent(configPath, "database.yml")
	if err != nil {
		return err
	}

	rawYaml := make(map[interface{}]interface{})

	err = yaml.Unmarshal(data, &rawYaml)
	if err != nil {
		return err
	}

	databaseYaml := rawYaml[Environment()].(map[interface{}]interface{})

	Config.DB.Name = databaseYaml["name"].(string)
	Config.DB.User = databaseYaml["user"].(string)
	Config.DB.Password = databaseYaml["password"].(string)
	Config.DB.Port = uint(databaseYaml["port"].(int))

	if databaseYaml["host"] != nil {
		Config.DB.Host = databaseYaml["host"].(string)
	}

	return nil
}

func getConfigDirectory(path string) (string, error) {
	paths := [3]string{
		path,
		os.Getenv("EDWIG_CONFIG"),
		fmt.Sprintf("%s/src/bitbucket.org/enroute-mobi/ara/config", os.Getenv("GOPATH")),
	}
	for _, directoryPath := range paths {
		if found := checkDirectory(directoryPath); found {
			return directoryPath, nil
		}
	}
	return "", errors.New("can't find config directory")
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
	return nil, fmt.Errorf("can't find %s configuration file", file)
}
