package image

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
	imagePath := fmt.Sprintf(".cache/images/%s.bmp", gameName)

	// Verificar se o diretório existe, se não, criar
	dir := fmt.Sprintf(".cache/images")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755) // Cria o diretório e os subdiretórios necessários
		if err != nil {
			output.Errorf("Error creating directory: %v\n", err)
			return ""
		}
	}

	if _, err := os.Stat(imagePath); err == nil {
		output.Errorf("Image found on disk: %s\n", imagePath)
		return imagePath
	}

	imageURL := fmt.Sprintf("https://handheld-database.github.io/handheld-database/commons/images/games/%s.icon.webp", gameName)
	response, err := http.Get(imageURL)
	if err != nil {
		output.Errorf("HTTP request error: %v\n", err)
		return ""
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		imageData, err := io.ReadAll(response.Body)
		if err != nil {
			output.Errorf("Error reading image data: %v\n", err)
			return ""
		}

		if response.Header.Get("Content-Type") == "image/webp" {
			output.Printf("Image in webp format, converting to BMP...")
			convertedImageData, err := ConvertWebpToBMP(imageData)
			if err != nil {
				output.Errorf("Error converting webp image to BMP: %v\n", err)
				return ""
			}

			imageData = convertedImageData
		}

		err = os.WriteFile(imagePath, imageData, 0644)

		if err != nil {
			output.Errorf("Error saving image to disk: %v\n", err)
			return ""
		}

		output.Printf("Image saved at: %s\n", imagePath)

		return imagePath
	}

	output.Errorf("Error retrieving game image: %s, status code: %d\n", gameName, response.StatusCode)
	return ""
}

func ConvertWebpToBMP(webpData []byte) ([]byte, error) {
	// Create a byte buffer to store the BMP image
	var bmpBuffer bytes.Buffer

	// Decode the WebP data
	img, err := webp.Decode(bytes.NewReader(webpData))
	if err != nil {
		output.Errorf("failed to decode WebP image: %v", err)
		return nil, err
	}

	// Write the image to the buffer as BMP
	err = bmp.Encode(&bmpBuffer, img)
	if err != nil {
		output.Errorf("failed to encode BMP image: %v", err)
		return nil, err
	}

	return bmpBuffer.Bytes(), nil
}
