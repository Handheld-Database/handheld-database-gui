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
		16.0 / 9.0: "16_9",
		4.0 / 3.0:  "4_3",
		3.0 / 2.0:  "3_2",
		21.0 / 9.0: "21_9",
		32.0 / 9.0: "32_9",
		1.0:        "1_1",
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
		"16_9": 15,
		"4_3":  15,
		"3_2":  15,
		"21_9": 15,
		"32_9": 15,
		"1_0":  15,
	}

	return aspectRatios[aspectRatio]
}

func calculateMaxLineWidth(aspectRatio string) int {
	aspectRatios := map[string]int{
		"16_9": 30,
		"4_3":  25,
		"3_2":  20,
		"21_9": 60,
		"32_9": 50,
		"1_0":  10,
	}

	return aspectRatios[aspectRatio]
}

func calculateMaxListItemWidth(aspectRatio string) int {
	aspectRatios := map[string]int{
		"16_9": 60,
		"4_3":  50,
		"3_2":  20,
		"21_9": 60,
		"32_9": 50,
		"1_0":  10,
	}

	return aspectRatios[aspectRatio]
}

func calculateMaxListItens(aspectRatio string) int {
	aspectRatios := map[string]int{
		"16_9": 10,
		"4_3":  11,
		"3_2":  15,
		"21_9": 10,
		"32_9": 10,
		"1_0":  10,
	}

	return aspectRatios[aspectRatio]
}
