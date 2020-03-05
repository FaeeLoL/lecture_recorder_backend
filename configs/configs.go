package configs

import (
	"encoding/json"
	"os"
)

type DBConfig struct {
	Address  string `yaml:"address"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

var DBConfigs DBConfig

const ConfigFilename = "db_configs.json"

func init() {
	configFile, err := os.Open(ConfigFilename)
	if err != nil {
		return
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	_ = jsonParser.Decode(&DBConfigs)
}
