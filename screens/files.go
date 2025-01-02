package screens

import (
	"context"
	"handheldui/components"
	"handheldui/helpers/sdlutils"
	"handheldui/input"
	"handheldui/output"
	"handheldui/services"
	"handheldui/vars"
	"sort"

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
	listComponent := components.NewListComponent(renderer, 19, func(index int, item map[string]interface{}) string {
		return item["name"].(string)
	})

	progressBar := components.NewProgressBarComponent(renderer, 300, 20, 490, 320, sdl.Color{R: 200, G: 200, B: 200, A: 255}, sdl.Color{R: 0, G: 255, B: 0, A: 255})

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

	var repositoriesList []string

	for _, repo := range vars.Config.Repositories {
		if repoDetails, ok := repo[vars.CurrentRepo]; ok {
			f.repoName = repoDetails.Name
			f.repoPath = repoDetails.Path

			repositoriesList = repoDetails.Repositories

			// Calls FetchAndSortAllMetadata to retrieve metadata from all repositories
			allMetadata, err := services.FetchAllMetadata(repositoriesList)
			if err != nil {
				panic(output.Sprintf("Error fetching metadata: %v", err))
			}

			// Initializes an items slice
			var items []map[string]interface{}

			// Processes the fetched metadata
			for _, metadata := range allMetadata {
				for fileName, fileURL := range metadata {
					items = append(items, map[string]interface{}{
						"name":  fileName,
						"value": fileURL,
					})
				}
			}

			// Sorts items before updating the list
			sort.Slice(items, func(i, j int) bool {
				return items[i]["name"].(string) < items[j]["name"].(string)
			})

			// Updates the list of items in the component
			f.listComponent.SetItems(items)

			break
		}
	}

	f.initialized = true
}

func (f *FilesScreen) HandleInput(event input.InputEvent) {
	if len(f.listComponent.GetItems()) == 0 {
		return
	}

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
		url := selectedItem["value"].(string)
		name := selectedItem["name"].(string)
		go f.download(f.repoPath, name, url)
	case "B":
		if f.isDownloading {
			if f.cancelDownload != nil {
				f.cancelDownload() // Cancels the current download
			}
			f.isDownloading = false
			f.progressBar.SetProgress(0.0)
			f.cancelDownload = nil
		} else {
			f.initialized = false
			vars.CurrentScreen = "home_screen"
		}
	}
}

func (f *FilesScreen) Draw() {
	f.InitRepositories()

	f.renderer.SetDrawColor(255, 255, 255, 255)
	f.renderer.Clear()

	// Checks if the download is in progress
	if f.isDownloading {

		// Displays only the progress bar
		sdlutils.RenderTexture(f.renderer, "assets/textures/bg.bmp", "Q2", "Q4")
		f.progressBar.Draw()
	} else {

		// Displays the list and other screen elements
		sdlutils.RenderTexture(f.renderer, "assets/textures/bg.bmp", "Q2", "Q4")

		// Draws the current title
		sdlutils.DrawText(f.renderer, f.repoName, sdl.Point{X: 25, Y: 25}, vars.Colors.PRIMARY, vars.HeaderFont)

		// Draws the list component
		f.listComponent.Draw(vars.Colors.WHITE, vars.Colors.SECONDARY)

		sdlutils.RenderTexture(f.renderer, "assets/textures/ui_controls_1280_720.bmp", "Q3", "Q4")
	}

	f.renderer.Present()
}

func (f *FilesScreen) download(path string, fileName string, uri string) {
	// Creates a context to cancel the download
	ctx, cancel := context.WithCancel(context.Background())
	f.cancelDownload = cancel

	f.progressBar.SetProgress(0.0)
	f.isDownloading = true

	err := services.DownloadFile(ctx, path, fileName, uri, func(downloaded, total int64) {
		// Updates the progress
		f.progressBar.SetProgress(float64(downloaded) / float64(total) * 100)
		if f.progressBar.GetProgress() >= 100 {
			f.isDownloading = false
			f.cancelDownload = nil
		}
	})

	if err != nil {
		output.Errorf("Error during download", err)
		f.isDownloading = false
		f.cancelDownload = nil
	}
}
