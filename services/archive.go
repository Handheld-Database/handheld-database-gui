package services

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"handheldui/output"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

var (
	cacheLock sync.RWMutex
)

// File represents the structure of a file in the XML metadata.
type File struct {
	Name string `xml:"name,attr"`
}

// Files represents the root structure of the XML metadata.
type Files struct {
	File []File `xml:"file"`
}

// Returns the cache file path specific to the given name
func getCacheFilePath(name string) string {
	return filepath.Join(".cache", "archive_metadata", fmt.Sprintf("cache_%s.json", name))
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

	// Decodes the data into the cache
	if err := json.NewDecoder(cacheFile).Decode(&cache); err != nil {
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
	if err := json.NewEncoder(cacheFile).Encode(cache); err != nil {
		return output.Errorf("error encoding cache data: %v", err)
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
	resp, err := http.Get(fmt.Sprintf("https://archive.org/download/%s/%s_files.xml", name, name))
	if err != nil {
		return nil, output.Errorf("error fetching metadata for %s: %v", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, output.Errorf("error fetching metadata for %s: %v", name, resp.Status)
	}

	// Decodes the metadata
	var metadata Files
	if err := xml.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, output.Errorf("error decoding response: %v", err)
	}

	// Processes the metadata
	metadataList := make(map[string]string)
	for _, file := range metadata.File {
		escapedFileName := strings.ReplaceAll(file.Name, " ", "%20")
		metadataList[file.Name] = fmt.Sprintf("https://archive.org/download/%s/%s", name, escapedFileName)
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

func DownloadFile(ctx context.Context, path, filename, link string, progress func(int64, int64)) error {
	// Ensure the destination directory is created
	if err := os.MkdirAll(path, 0755); err != nil {
		return output.Errorf("error creating directory %s: %v", path, err)
	}

	// Sanitize the filename to prevent issues with special characters
	sanitizedFilename := filepath.Base(filename)
	fullPath := filepath.Join(path, sanitizedFilename)

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
