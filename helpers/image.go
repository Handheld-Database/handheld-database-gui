package helpers

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"os"

	"golang.org/x/image/webp"
)

func FetchGameImage(gameName string) string {
	imagePath := fmt.Sprintf("./tmp/%s.png", gameName)

	if _, err := os.Stat(imagePath); err == nil {
		fmt.Printf("Imagem encontrada em disco: %s\n", imagePath)
		return imagePath
	}

	imageURL := fmt.Sprintf("https://handheld-database.github.io/handheld-database/commons/images/games/%s.icon.webp", gameName)
	response, err := http.Get(imageURL)
	if err != nil {
		fmt.Printf("Erro na requisição HTTP: %v\n", err)
		return ""
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		imageData, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Erro ao ler dados da imagem: %v\n", err)
			return ""
		}

		if response.Header.Get("Content-Type") == "image/webp" {
			fmt.Println("Imagem no formato webp, convertendo para PNG...")
			convertedImageData, err := ConvertWebpToPNG(imageData)
			if err != nil {
				fmt.Printf("Erro ao converter imagem webp para PNG: %v\n", err)
				return ""
			}
			imageData = convertedImageData
		}

		err = os.WriteFile(imagePath, imageData, 0644)
		if err != nil {
			fmt.Printf("Erro ao salvar imagem no disco: %v\n", err)
			return ""
		}

		fmt.Printf("Imagem salva em: %s\n", imagePath)
		return imagePath
	}

	fmt.Printf("Erro ao recuperar a imagem do jogo: %s, status code: %d\n", gameName, response.StatusCode)
	return ""
}

func ConvertWebpToPNG(webpData []byte) ([]byte, error) {
	// Crie um buffer de bytes para armazenar a imagem PNG
	var pngBuffer bytes.Buffer

	// Decodifique os dados WebP
	img, err := webp.Decode(bytes.NewReader(webpData))
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar imagem WebP: %v", err)
	}

	// Escreva a imagem no buffer como PNG
	err = png.Encode(&pngBuffer, img)
	if err != nil {
		return nil, fmt.Errorf("falha ao codificar imagem PNG: %v", err)
	}

	return pngBuffer.Bytes(), nil
}
