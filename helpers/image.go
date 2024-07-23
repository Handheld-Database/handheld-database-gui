package helpers

import (
	"bytes"
	"fmt"
	"handheldui/output"
	"io"
	"net/http"
	"os"

	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
)

func FetchGameImage(gameName string) string {
	imagePath := fmt.Sprintf("./tmp/%s.bmp", gameName)

	if _, err := os.Stat(imagePath); err == nil {
		output.Printf("Image found on disk: %s\n", imagePath)
		return imagePath
	}

	imageURL := fmt.Sprintf("https://handheld-database.github.io/handheld-database/commons/images/games/%s.icon.webp", gameName)
	response, err := http.Get(imageURL)
	if err != nil {
		output.Printf("HTTP request error: %v\n", err)
		return ""
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		imageData, err := io.ReadAll(response.Body)
		if err != nil {
			output.Printf("Error reading image data: %v\n", err)
			return ""
		}

		if response.Header.Get("Content-Type") == "image/webp" {
			output.Printf("Image in webp format, converting to BMP...")
			convertedImageData, err := ConvertWebpToBMP(imageData)
			if err != nil {
				output.Printf("Error converting webp image to BMP: %v\n", err)
				return ""
			}
			imageData = convertedImageData
		}

		err = os.WriteFile(imagePath, imageData, 0644)
		if err != nil {
			output.Printf("Error saving image to disk: %v\n", err)
			return ""
		}

		output.Printf("Image saved at: %s\n", imagePath)
		return imagePath
	}

	output.Printf("Error retrieving game image: %s, status code: %d\n", gameName, response.StatusCode)
	return ""
}

func ConvertWebpToBMP(webpData []byte) ([]byte, error) {
	// Create a byte buffer to store the BMP image
	var bmpBuffer bytes.Buffer

	// Decode the WebP data
	img, err := webp.Decode(bytes.NewReader(webpData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode WebP image: %v", err)
	}

	// Write the image to the buffer as BMP
	err = bmp.Encode(&bmpBuffer, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode BMP image: %v", err)
	}

	return bmpBuffer.Bytes(), nil
}
