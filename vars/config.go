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
	Width            int32 `json:"width"`
	Height           int32 `json:"height"`
	AspectRatio      string
	MaxLines         int
	MaxListItemWidth int
	MaxListItens     int
}

type ConfigDefinition struct {
	Logs         bool                       `json:"logs"`
	Control      map[string]string          `json:"control"`
	Screen       ScreenDetails              `json:"screen"`
	Repositories map[string]PlatformDetails `json:"repositories"`
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

	aspectRation := calculateAspectRatio(config.Screen.Width, config.Screen.Height)
	config.Screen.AspectRatio = aspectRation

	config.Screen.MaxLines = calculateMaxLines(aspectRation)
	config.Screen.MaxListItens = calculateMaxListItens(aspectRation)
	config.Screen.MaxListItemWidth = calculateMaxListItemWidth(aspectRation)

	return &config, nil
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

func calculateMaxLines(aspectRatio string) int {
	aspectRatios := map[string]int{
		"16:9": 15,
		"4:3":  15,
		"3:2":  15,
		"21:9": 15,
		"32:9": 15,
		"1:0":  15,
	}

	return aspectRatios[aspectRatio]
}

func calculateMaxLineWidth(aspectRatio string) int {
	aspectRatios := map[string]int{
		"16:9": 30,
		"4:3":  25,
		"3:2":  20,
		"21:9": 60,
		"32:9": 50,
		"1:0":  10,
	}

	return aspectRatios[aspectRatio]
}

func calculateMaxListItemWidth(aspectRatio string) int {
	aspectRatios := map[string]int{
		"16:9": 30,
		"4:3":  25,
		"3:2":  20,
		"21:9": 60,
		"32:9": 50,
		"1:0":  10,
	}

	return aspectRatios[aspectRatio]
}

func calculateMaxListItens(aspectRatio string) int {
	aspectRatios := map[string]int{
		"16:9": 10,
		"4:3":  11,
		"3:2":  15,
		"21:9": 10,
		"32:9": 10,
		"1:0":  10,
	}

	return aspectRatios[aspectRatio]
}
