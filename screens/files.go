package screens

import (
	"context"
	"handheldui/components"
	"handheldui/helpers/sdlutils"
	"handheldui/helpers/wrappers"
	"handheldui/input"
	"handheldui/output"
	"handheldui/services"
	"handheldui/vars"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

type FilesScreen struct {
	initialized    bool
	renderer       *sdl.Renderer
	listComponent  *components.ListComponent
	repoName       string
	repoPath       string
	progressBar    *components.ProgressBarComponent
	isDownloading  bool
	cancelDownload context.CancelFunc
}

func NewFilesScreen(renderer *sdl.Renderer) (*FilesScreen, error) {
	listComponent := components.NewListComponent(
		renderer,
		vars.Config.Screen.MaxListItens,
		vars.Config.Screen.MaxListItemWidth,
		func(index int, item map[string]interface{}) string {
			return item["name"].(string)
		})

	progressBar := components.NewProgressBarComponent(renderer, 300, 20, 490, 320, vars.Colors.WHITE, vars.Colors.SECONDARY)

	return &FilesScreen{
		renderer:      renderer,
		listComponent: listComponent,
		progressBar:   progressBar,
	}, nil
}

func (f *FilesScreen) InitRepositories() {
	if f.initialized {
		return
	}

	if currentRepoDetails, ok := vars.Config.Repositories[vars.CurrentRepo]; ok {
		f.repoName = currentRepoDetails.Name
		f.repoPath = currentRepoDetails.Path

		collections := currentRepoDetails.Collections
		extList := currentRepoDetails.ExtList

		// Initializes an items slice
		var items []map[string]interface{}

		// Calls FetchAndSortAllMetadata to retrieve metadata from all repositories
		for _, collection := range collections {

			allMetadata, err := services.FetchMetadata(collection.Name)
			if err != nil {
				panic(output.Sprintf("Error fetching metadata: %v", err))
			}

			// Process the fetched metadata
			for fileName, fileURL := range allMetadata {
				// If extList is empty, add all files
				if len(extList) == 0 {
					items = append(items, map[string]interface{}{
						"name":  fileName,
						"value": fileURL,
						"unzip": collection.Unzip,
					})
				} else {
					// Check if the file has one of the specified extensions
					for _, ext := range extList {
						if strings.HasSuffix(fileName, ext) {
							items = append(items, map[string]interface{}{
								"name":  fileName,
								"value": fileURL,
								"unzip": collection.Unzip,
							})
							break
						}
					}
				}
			}
		}

		// Sorts items before updating the list
		sort.Slice(items, func(i, j int) bool {
			return items[i]["name"].(string) < items[j]["name"].(string)
		})

		// Updates the list of items in the component
		f.listComponent.SetItems(items)
	}

	f.initialized = true
}

func (f *FilesScreen) HandleInput(event input.InputEvent) {
	// Handle the B button regardless of the list state
	if event.KeyCode == "B" {
		if f.isDownloading {
			if f.cancelDownload != nil {
				f.cancelDownload() // Cancels the current download
			}
			f.isDownloading = false
			f.progressBar.SetProgress(0.0)
			f.cancelDownload = nil
		} else {
			f.initialized = false
			vars.CurrentScreen = "repositories_screen"
		}
		return
	}

	// Skip other input handling if the list is empty
	if len(f.listComponent.GetItems()) == 0 {
		return
	}

	// Handle other inputs
	switch event.KeyCode {
	case "DOWN":
		f.listComponent.ScrollDown()
	case "UP":
		f.listComponent.ScrollUp()
	case "L1":
		f.listComponent.PageUp()
	case "R1":
		f.listComponent.PageDown()
	case "A":
		selectedItem := f.listComponent.GetItems()[f.listComponent.GetSelectedIndex()]
		go f.downloadFile(f.repoPath, selectedItem)
	}
}

func (f *FilesScreen) Draw() {
	f.InitRepositories()

	f.renderer.SetDrawColor(255, 255, 255, 255)
	f.renderer.Clear()

	// Checks if the download is in progress
	if f.isDownloading {

		// Displays only the progress bar
		sdlutils.RenderTextureCover(f.renderer, "assets/textures/bg.bmp")

		f.progressBar.Draw()

		sdlutils.RenderTextureCartesian(f.renderer, "assets/textures/$aspect_ratio/ui_controls_download.bmp", "Q3", "Q4")

	} else {

		// Displays the list and other screen elements
		sdlutils.RenderTextureCartesian(f.renderer, "assets/textures/bg.bmp", "Q2", "Q4")

		// Draws the current title
		sdlutils.DrawText(f.renderer, f.repoName, sdl.Point{X: 25, Y: 25}, vars.Colors.WHITE, vars.HeaderFont)

		// Draws the list component
		f.listComponent.Draw(vars.Colors.SECONDARY, vars.Colors.WHITE)

		sdlutils.RenderTextureCartesian(f.renderer, "assets/textures/$aspect_ratio/ui_controls.bmp", "Q3", "Q4")
	}

	f.renderer.Present()
}

func (f *FilesScreen) downloadFile(path string, selectedItem map[string]interface{}) {
	// get variables
	uri := selectedItem["value"].(string)
	fileName := selectedItem["name"].(string)
	unzip := selectedItem["unzip"].(bool)

	// Creates a context to cancel the download
	ctx, cancel := context.WithCancel(context.Background())
	f.cancelDownload = cancel

	f.progressBar.SetProgress(0.0)
	f.isDownloading = true

	// Download file
	err := services.DownloadFile(ctx, path, fileName, uri, func(downloaded, total int64) {
		// Update progress
		f.progressBar.SetProgress(float64(downloaded) / float64(total) * 100)
		if f.progressBar.GetProgress() >= 100 {
			f.isDownloading = false
			if unzip {
				splitedPath := strings.Split(fileName, "/")
				deepSize := len(splitedPath) - 1
				f.unzipFile(path, splitedPath[deepSize])
			}
		}
	})

	if err != nil {
		output.Errorf("Error during download", err)
		f.isDownloading = false
		f.cancelDownload = nil
		return
	}
}

func (f *FilesScreen) unzipFile(path string, fileName string) {
	// Caminho completo do arquivo ZIP
	zipFilePath := filepath.Join(path, fileName)

	// Destino da extração (diretório atual sem subpastas adicionais)
	destDir := path

	// Executa a descompactação em uma goroutine
	go func() {
		// Chama a função de descompactação com os caminhos corretos
		err := wrappers.UnzipFile(zipFilePath, destDir)
		if err != nil {
			output.Printf("Error extracting file: %v\n", err)
			return
		}

		// Remove o arquivo ZIP após a extração
		err = os.Remove(zipFilePath)
		if err != nil {
			output.Printf("Error removing zip file: %v\n", err)
		} else {
			output.Printf("Successfully removed zip file: %s\n", zipFilePath)
		}
	}()
}
