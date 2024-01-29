package configuration_manager

import (
	"encoding/json"
	"os"
)

// TODO: Use DB with permission + resource id + path/file name

type Config struct {
	ListenAddress             string `json:"listenAddress"`
	CertificatePath           string `json:"certificatePath"`
	CertificatePrivateKeyPath string `json:"certificatePrivateKeyPath"`
	DatabaseConnectionString  string `json:"databaseConnectionString"`
	GlobalRequestsPerSecond   int    `json:"globalRequestsPerSecond"`
}

var Configuration *Config

func LoadConfiguration(file string) (*Config, error) {
	if Configuration != nil {
		return Configuration, nil
	}

	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return &config, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	Configuration = &config
	return &config, err
}
