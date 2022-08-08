package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/enroute-mobi/ara/logger"
	yaml "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     uint
}

type config struct {
	DB DatabaseConfig

	ApiKey                string
	Debug                 bool
	BigQueryProjectID     string
	BigQueryDatasetPrefix string
	BigQueryTest          string
	Sentry                string
	Syslog                bool
	ColorizeLog           bool
	LoadMaxInsert         int
	FakeUUIDRealFormat    bool
}

var Config = config{}

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

	// Test env variables
	bigQueryTestEnv := os.Getenv("ARA_BIGQUERY_TEST")
	if bigQueryTestEnv != "" {
		Config.BigQueryTest = bigQueryTestEnv
	}
	bigQueryPrefixEnv := os.Getenv("ARA_BIGQUERY_PREFIX")
	if bigQueryPrefixEnv != "" {
		Config.BigQueryDatasetPrefix = bigQueryPrefixEnv
	}
	fakeUUIDReal := os.Getenv("ARA_FAKEUUID_REAL")
	if strings.ToLower(fakeUUIDReal) == "true" {
		Config.FakeUUIDRealFormat = true
	}

	logger.Log.Syslog = Config.Syslog
	logger.Log.Debug = Config.Debug
	logger.Log.Color = Config.ColorizeLog

	return nil
}

func (c *config) ValidBQConfig() bool {
	return c.BigQueryTestMode() || (c.BigQueryProjectID != "" && c.BigQueryDatasetPrefix != "")
}

func (c *config) BigQueryTestMode() bool {
	return c.BigQueryTest != ""
}

var environment string

func SetEnvironment(env string) {
	environment = env
}

func Environment() string {
	if environment == "" {
		env := os.Getenv("ARA_ENV")
		if env == "" {
			logger.Log.Debugf("ARA_ENV not set, default environment is development")
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
	if databaseYaml["user"] != nil {
		Config.DB.User = databaseYaml["user"].(string)
	}
	if databaseYaml["password"] != nil {
		Config.DB.Password = databaseYaml["password"].(string)
	}
	if databaseYaml["port"] != nil {
		Config.DB.Port = uint(databaseYaml["port"].(int))
	}

	if databaseYaml["host"] != nil {
		Config.DB.Host = databaseYaml["host"].(string)
	}

	return nil
}

func getConfigDirectory(path string) (string, error) {
	paths := [3]string{
		path,
		os.Getenv("ARA_CONFIG"),
		fmt.Sprintf("%s/src/bitbucket.org/enroute-mobi/ara/config", os.Getenv("GOPATH")),
	}
	for _, directoryPath := range paths {
		if found := checkDirectory("config", directoryPath); found {
			return directoryPath, nil
		}
	}
	return "", errors.New("can't find config directory")
}

func GetTemplateDirectory() (string, error) {
	paths := [2]string{
		os.Getenv("ARA_ROOT"),
		fmt.Sprintf("%s/src/bitbucket.org/enroute-mobi/ara", os.Getenv("GOPATH")),
	}
	for _, directoryPath := range paths {
		templatePath := filepath.Join(directoryPath, "/siri/templates")
		if found := checkDirectory("template", templatePath); found {
			return templatePath, nil
		}
	}
	return "", errors.New("can't find template directory")
}

func checkDirectory(kind, path string) bool {
	if path == "" {
		return false
	}
	if _, err := os.Stat(path); err == nil {
		logger.Log.Debugf("Found %v directory at %s", kind, path)
		return true
	}
	logger.Log.Debugf("Can't find %v directory at %s", kind, path)
	return false
}

func getConfigFileContent(path, file string) ([]byte, error) {
	// Check file at location
	filePath := strings.Join([]string{path, file}, "/")
	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	logger.Log.Debugf("Can't find %s config file in %s", file, path)
	return nil, fmt.Errorf("can't find %s configuration file", file)
}
