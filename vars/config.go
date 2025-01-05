package vars

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"math"
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

type ScreenDetails struct {
	Width       int32 `json:"width"`
	Height      int32 `json:"height"`
	AspectRatio string
}

type ConfigDefinition struct {
	Logs         bool                       `json:"logs"`
	Control      map[string]string          `json:"control"`
	Screen       ScreenDetails              `json:"screen"`
	Repositories map[string]PlatformDetails `json:"repositories"`
}

func calculateAspectRatio(width, height int32) string {
	if height == 0 {
		return "Unknown"
	}

	ratio := float64(width) / float64(height)

	aspectRatios := map[float64]string{
		16.0 / 9.0: "16:9",
		4.0 / 3.0:  "4:3",
		3.0 / 2.0:  "3:2",
		21.0 / 9.0: "21:9",
		32.0 / 9.0: "32:9",
		1.0:        "1:1",
	}

	tolerance := 0.01
	for k, v := range aspectRatios {
		if math.Abs(ratio-k) < tolerance {
			return v
		}
	}

	return fmt.Sprintf("%.2f:1", ratio)
}

func LoadConfig(configFile []byte) (*ConfigDefinition, error) {
	var config ConfigDefinition
	err := json.Unmarshal(configFile, &config)
	if err != nil {
		if config.Logs {
			log.Fatalf("Erro ao fazer o parse do arquivo embutido: %v", err)
		}
		return nil, err
	}

	config.Screen.AspectRatio = calculateAspectRatio(config.Screen.Width, config.Screen.Height)
	return &config, nil
}
