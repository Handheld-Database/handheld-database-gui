package vars

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

type RepositoryDetails struct {
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Repositories []string `json:"repositories"`
}

type ConfigDefinition struct {
	Repositories []map[string]RepositoryDetails `json:"repositories"`
}

func LoadConfig(configFile []byte) (*ConfigDefinition, error) {
	var config ConfigDefinition
	err := json.Unmarshal(configFile, &config)
	if err != nil && Debug {
		return nil, fmt.Errorf("erro ao fazer o parse do arquivo embutido: %v", err)
	}

	return &config, nil
}
