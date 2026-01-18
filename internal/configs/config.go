package configs

import (
	"fmt"
	"os"

	"github.com/Fiagram/gateway/configs"
	"gopkg.in/yaml.v3"
)

type ConfigFilePath string

type Config struct {
	Http  Http  `yaml:"http"`
	Auth  Auth  `yaml:"auth"`
	Log   Log   `yaml:"log"`
	Cache Cache `yaml:"cache"`
}

// Creates a new config instance by reading from a given YAML file.
// If the filePath is empty, it uses the default embedded configuration.
func NewConfig(filePath ConfigFilePath) (Config, error) {
	var (
		configBytes = configs.DefaultConfigBytes
		config      = Config{}
		err         error
	)

	if filePath != "" {
		configBytes, err = os.ReadFile(string(filePath))
		if err != nil {
			return Config{}, fmt.Errorf("Failed to read YAML file: %w", err)
		}
	}

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Failed to unmarshal YAML file: %w", err)
	}

	return config, nil
}

func GetConfigHttp(c Config) Http {
	return c.Http
}

func GetConfigLog(c Config) Log {
	return c.Log
}

func GetConfigCache(c Config) Cache {
	return c.Cache
}
