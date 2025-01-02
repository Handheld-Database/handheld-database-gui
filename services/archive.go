package services

import (
	"context"
	"encoding/json"
	"fmt"
	"handheldui/output"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

var (
	cacheLock sync.RWMutex
)

// Returns the cache file path specific to the given name
func getCacheFilePath(name string) string {
	return filepath.Join(".cache", "archive_metadata", output.Sprintf("cache_%s.json", name))
}

func loadCacheFromFile(name string) (map[string]string, error) {
	cacheFilePath := getCacheFilePath(name)
	cache := make(map[string]string)

	// Tries to open the cache file
	cacheFile, err := os.Open(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist, returns an empty cache
			return cache, nil
		}
		return nil, output.Errorf("error opening cache file: %v", err)
	}
	defer cacheFile.Close()

	// Reads the content of the file
	cacheData, err := ioutil.ReadAll(cacheFile)
	if err != nil {
		return nil, output.Errorf("error reading cache file: %v", err)
	}

	// Decodes the data into the cache
	if err := json.Unmarshal(cacheData, &cache); err != nil {
		return nil, output.Errorf("error decoding cache data: %v", err)
	}

	return cache, nil
}

func saveCacheToFile(name string, cache map[string]string) error {
	cacheFilePath := getCacheFilePath(name)

	// Creates necessary directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(cacheFilePath), os.ModePerm); err != nil {
		return output.Errorf("error creating cache directories: %v", err)
	}

	// Creates or overwrites the cache file
	cacheFile, err := os.Create(cacheFilePath)
	if err != nil {
		return output.Errorf("error creating or overwriting cache file: %v", err)
	}
	defer cacheFile.Close()

	// Encodes the cache data to the file
	cacheData, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return output.Errorf("error encoding cache data: %v", err)
	}

	// Writes the data to the file
	if _, err := cacheFile.Write(cacheData); err != nil {
		return output.Errorf("error writing to cache file: %v", err)
	}

	return nil
}

func FetchMetadata(name string) (map[string]string, error) {
	cacheLock.RLock()
	cache, err := loadCacheFromFile(name)
	cacheLock.RUnlock()
	if err != nil {
		return nil, err
	}

	// If the cache already contains data, return it
	if len(cache) > 0 {
		return cache, nil
	}

	// Downloads the metadata from the URL
	resp, err := http.Get(fmt.Sprintf("https://archive.org/metadata/%s", name))
	if err != nil {
		return nil, output.Errorf("error fetching metadata for %s: %v", name, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, output.Errorf("error fetching metadata for %s: %v", name, resp.Status)
	}

	// Decodes the metadata
	var metadata struct {
		Files []struct {
			Name string `json:"name"`
		} `json:"files"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, output.Errorf("error decoding response: %v", err)
	}

	// Processes the metadata
	metadataList := make(map[string]string)
	for _, file := range metadata.Files {
		metadataList[file.Name] = fmt.Sprintf("https://archive.org/download/%s/%s", name, file.Name)
	}

	// Sorts the files by name
	var sortedFileNames []string
	for fileName := range metadataList {
		sortedFileNames = append(sortedFileNames, fileName)
	}
	sort.Strings(sortedFileNames)

	// Creates a sorted list of metadata
	sortedMetadataList := make(map[string]string)
	for _, fileName := range sortedFileNames {
		sortedMetadataList[fileName] = metadataList[fileName]
	}

	// Updates the cache
	cacheLock.Lock()
	err = saveCacheToFile(name, sortedMetadataList)
	cacheLock.Unlock()
	if err != nil {
		return nil, err
	}

	return sortedMetadataList, nil
}

func FetchAllMetadata(names []string) ([]map[string]string, error) {
	var (
		allMetadata []map[string]string
		mu          sync.Mutex
		wg          sync.WaitGroup
		errs        []error
	)

	// Channel to synchronize results
	results := make(chan map[string]string, len(names))

	// Iterates over each name and fetches metadata
	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()

			metadata, err := FetchMetadata(name)
			if err != nil {
				mu.Lock()
				errs = append(errs, output.Errorf("error fetching metadata for %s: %v", name, err))
				mu.Unlock()
				return
			}

			results <- metadata
		}(name)
	}

	// Waits for all goroutines to finish
	wg.Wait()
	close(results)

	// Checks for accumulated errors
	if len(errs) > 0 {
		return nil, output.Errorf("errors while fetching metadata: %v", errs)
	}

	// Combines the results
	for metadata := range results {
		allMetadata = append(allMetadata, metadata)
	}

	return allMetadata, nil
}

func DownloadFile(ctx context.Context, path, filename, link string, progress func(int64, int64)) error {
	// Creates the destination directory if necessary
	if err := os.MkdirAll(path, 0755); err != nil {
		return output.Errorf("error creating directory %s: %v", path, err)
	}

	fullPath := filepath.Join(path, filename)

	// Starts the download
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return output.Errorf("error creating request for %s: %v", link, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return output.Errorf("error downloading file from %s: %v", link, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return output.Errorf("error downloading file from %s: %v", link, resp.Status)
	}

	totalSize := resp.ContentLength
	out, err := os.Create(fullPath)
	if err != nil {
		return output.Errorf("error creating file %s: %v", fullPath, err)
	}
	defer out.Close()

	var downloaded int64
	buf := make([]byte, 32*1024)
	for {
		select {
		case <-ctx.Done(): // Monitors cancellation
			return output.Errorf("download cancelled")
		default:
			n, err := resp.Body.Read(buf)
			if n > 0 {
				downloaded += int64(n)
				progress(downloaded, totalSize)

				if _, writeErr := out.Write(buf[:n]); writeErr != nil {
					return output.Errorf("error saving file %s: %v", fullPath, writeErr)
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return output.Errorf("error downloading file %s: %v", link, err)
			}
		}
	}
}
