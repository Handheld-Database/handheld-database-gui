package vars

import (
	_ "embed"
	"encoding/json"
	"log"
)

type CollectionDetails struct {
	Name  string `json:"name"`
	Unzip bool   `json:"unzip"`
}

type PlatformDetails struct {
	Name        string              `json:"name"`
	Path        string              `json:"path"`
	ExtList     []string            `json:"extlist"`
	Collections []CollectionDetails `json:"collections"`
}

type ConfigDefinition struct {
	Logs         bool                       `json:"logs"`
	Control      map[string]string          `json:"control"`
	Screen       map[string]int32           `json:"screen"`
	Repositories map[string]PlatformDetails `json:"repositories"`
}

func LoadConfig(configFile []byte) (*ConfigDefinition, error) {
	var config ConfigDefinition
	err := json.Unmarshal(configFile, &config)
	if err != nil && config.Logs {
		log.Fatalf("erro ao fazer o parse do arquivo embutido: %v", err)
		return nil, err
	}

	return &config, nil
}
