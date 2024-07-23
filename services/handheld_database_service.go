package services

import (
	"encoding/json"
	"fmt"
	"handheldui/helpers"
	"io"
	"net/http"
	"strings"
)

const baseURL = "https://handheld-database.github.io/handheld-database"

// GetRankColor returns the color for a given rank.
func GetRankColor(key string) string {
	colors := map[string]string{
		"PLATINUM": "rgb(180, 199, 220)",
		"GOLD":     "rgb(207, 181, 59)",
		"SILVER":   "rgb(166, 166, 166)",
		"BRONZE":   "rgb(205, 127, 50)",
		"FAULTY":   "red",
	}
	return colors[key]
}

// FetchPlatformsIndex fetches the index of platforms.
func FetchPlatformsIndex() ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/platforms/index.json", baseURL))
	if err != nil {
		return nil, fmt.Errorf("error fetching popular platforms: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching popular platforms: %v", resp.Status)
	}

	var result struct {
		Platforms []string `json:"platforms"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.Platforms, nil
}

// FetchPlatform fetches data for a given platform.
func FetchPlatform(platformKey string) (map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/platforms/%s/index.json", baseURL, platformKey))
	if err != nil {
		return nil, fmt.Errorf("error fetching systems from %s: %v", platformKey, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching systems from %s: %v", platformKey, resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result, nil
}

// FetchGames fetches games for a given platform and system.
func FetchGames(platformKey, systemKey string) ([]map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/platforms/%s/systems/%s/index.json", baseURL, platformKey, systemKey))
	if err != nil {
		return nil, fmt.Errorf("error fetching games from %s/%s: %v", platformKey, systemKey, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching games from %s/%s: %v", platformKey, systemKey, resp.Status)
	}

	var result struct {
		Games []map[string]interface{} `json:"games"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.Games, nil
}

// FetchGameDetails fetches details for a given game.
func FetchGameDetails(platformKey, systemKey, gameKey string) (map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/platforms/%s/systems/%s/%s/%s.json", baseURL, platformKey, systemKey, gameKey, gameKey))
	if err != nil {
		return nil, fmt.Errorf("error fetching game details from %s/%s/%s: %v", platformKey, systemKey, gameKey, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching game details from %s/%s/%s: %v", platformKey, systemKey, gameKey, resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result, nil
}

// FetchGameOverview fetches the overview for a given game.
func FetchGameOverview(gameKey string) (string, error) {
	fmt.Printf("%s/commons/overviews/%s.overview.md", baseURL, gameKey)
	resp, err := http.Get(fmt.Sprintf("%s/commons/overviews/%s.overview.md", baseURL, gameKey))
	if err != nil {
		return "", fmt.Errorf("error fetching game overview: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error fetching game overview: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return helpers.MarkdownToPlaintext(string(body)), nil
}

// FetchGameMarkdown fetches the markdown content for a given game.
func FetchGameMarkdown(platformKey, systemKey, gameKey string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/platforms/%s/systems/%s/%s/%s.md", baseURL, platformKey, systemKey, gameKey, gameKey))
	if err != nil {
		return "", fmt.Errorf("error fetching game markdown from %s/%s/%s: %v", platformKey, systemKey, gameKey, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error fetching game markdown from %s/%s/%s: %v", platformKey, systemKey, gameKey, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

// FilterGames filters games based on a search term and rank filter.
func FilterGames(games []map[string]interface{}, searchTerm, rankFilter string) []map[string]interface{} {
	var filteredGames []map[string]interface{}

	for _, game := range games {
		name := game["name"].(string)
		rank := game["rank"].(string)

		if searchTerm != "" && !containsIgnoreCase(name, searchTerm) {
			continue
		}

		if rankFilter != "ALL" && rank != rankFilter {
			continue
		}

		filteredGames = append(filteredGames, game)
	}

	return filteredGames
}

// FetchCollaborators fetches the list of collaborators.
func FetchCollaborators() (map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/commons/collaborators/collaborators.json", baseURL))
	if err != nil {
		return nil, fmt.Errorf("error fetching collaborators: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching collaborators: %v", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result, nil
}

// Helper function to check if a string contains another string, case-insensitive.
func containsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}
