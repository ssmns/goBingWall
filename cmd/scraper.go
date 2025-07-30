package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"

	core "github.com/ssmns/goBingWall/core"

	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var runScrapperCmd = func(cmd *cobra.Command, args []string) {
	// If no path provided, use default
	if savePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		savePath = filepath.Join(homeDir, "BingWall")
		fmt.Println(core.Magenta("Image Store Path:") + savePath)
	}

	// Create save directory if it doesn't exist
	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// Optional: Set wallpaper if flag is set
	if setWallpaper {
		// Get current images
		currentImages := core.GetImagesList(savePath)
		rand.New(rand.NewSource(125))
		randomIndex := rand.Intn(len(currentImages))
		fmt.Println(currentImages[randomIndex])
		core.SetWallpaper(currentImages[randomIndex])

	} else {
		// Get current images
		currentImages := core.GetCurrentImages(savePath)

		// Visit main page
		doc := core.VisitURL(baseURL)

		// Get archives
		links := getArchives(doc)

		// Download images
		err := getImages(baseURL, links, currentImages, savePath, resolution)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func getArchives(doc *goquery.Document) core.Archives {
	archives := make(map[string]string)
	doc.Find(".container.mt-3.pb-3 a").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		href, _ := s.Attr("href")
		archives[text] = href
	})

	// Convert map to a slice and sort
	var sortedArchives core.Archives
	for name, url := range archives {
		sortedArchives = append(sortedArchives, struct {
			Name string
			URL  string
		}{Name: name, URL: url})
	}

	sort.Slice(sortedArchives, func(i, j int) bool {
		return sortedArchives[i].Name > sortedArchives[j].Name
	})

	return sortedArchives
}

func getImageList(doc *goquery.Document) []core.Image {
	var images []core.Image
	doc.Find(".row.align-items-start a").Each(func(i int, s *goquery.Selection) {
		detail, _ := s.Attr("data-bs-title")
		href, _ := s.Attr("href")
		if detail != "" {
			images = append(images, core.Image{
				Detail: detail,
				Name:   filepath.Base(href),
				URL:    href,
			})
		}
	})
	return images
}

func getImage(baseURL string, image core.Image) core.ImageDownload {
	doc := core.VisitURL(baseURL + image.URL)
	imgDownload := core.ImageDownload{
		Name: image.Name,
		URLs: make(map[string]string),
	}

	doc.Find(".row.align-items-end a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		size := strings.Split(strings.Split(href, "w:")[1], "/")[0]
		switch size {
		case "3840":
			imgDownload.URLs["4K"] = href
		case "2560":
			imgDownload.URLs["2K"] = href
		case "1920":
			imgDownload.URLs["FHD"] = href
		default:
			imgDownload.URLs["default"] = href
		}
	})
	return imgDownload
}

func getImages(baseURL string, links core.Archives, currentImages map[string]bool, savePath, resolution string) error {

	// Get Absolute Path
	absPath, err := filepath.Abs(savePath)
	if err != nil {
		log.Printf("Error getting absolute path: %v", err)
		absPath = savePath // fallback to original path if abs fails
	}
	fmt.Println(core.Magenta("Save Path: ") + absPath)

	// Create a progress bar for total archives
	// archiveBar := progressbar.Default(int64(len(links)), "Processing Archives")

	for _, archiveStruct := range links {

		doc := core.VisitURL(baseURL + archiveStruct.URL)
		images := getImageList(doc)

		// Create a separate progress bar for images in this archive
		imageBar := progressbar.NewOptions(int(len(images)),
			progressbar.OptionSetDescription(fmt.Sprintf("Images in %s", archiveStruct.Name)),
			progressbar.OptionShowCount(),
			progressbar.OptionShowBytes(false),
		)

		for _, img := range images {
			// Update image progress bar
			imageBar.Describe(core.Green("Processing ") + formatProgressDescription(archiveStruct.Name, img.Name, 6, 25))
			// imageBar.Describe(fmt.Sprintf("Processing %s %s", archiveName, img.Name))

			if _, exists := currentImages[img.Name]; !exists {
				imgDownload := getImage(baseURL, img)

				// Save image based on resolution
				if resolution == "all" {
					for res, url := range imgDownload.URLs {
						savePath := filepath.Join(savePath, imgDownload.Name+"-"+res+".jpg")
						err := core.SaveImage(url, savePath)
						if err != nil {
							// log.Printf(core.Red("Error")+" saving image %s: %v", savePath, err)
							continue
						}
						// fmt.Println(core.Green("Save:") + imgDownload.Name)
					}
				} else {
					if url, ok := imgDownload.URLs[resolution]; ok {
						savePath := filepath.Join(savePath, imgDownload.Name+"-"+resolution+".jpg")
						err := core.SaveImage(url, savePath)
						if err != nil {
							// log.Printf(core.Red("Error")+" saving image %s: %v", savePath, err)
							continue
						}
						// fmt.Println(core.Green("Save:") + imgDownload.Name)
					}
				}

				// Save image details
				detailPath := filepath.Join(savePath, img.Name+".txt")
				err := core.SaveText(img.Detail, detailPath)
				if err != nil {
					log.Printf("Error saving details %s: %v", detailPath, err)
				}
			} else {
				// fmt.Println(core.Yellow("Ignore:") + img.Name)

			}

			// Increment image progress bar
			imageBar.Add(1)
		}

	}
	return nil
}

// Utility function for fixed-width description
func formatProgressDescription(archiveName string, imgName string, maxArchiveWidth int, maxImgWidth int) string {
	// Truncate or pad archive name
	if len(archiveName) > maxArchiveWidth {
		archiveName = archiveName[:maxArchiveWidth-3] + "..."
	} else {
		archiveName = fmt.Sprintf("%-*s", maxArchiveWidth, archiveName)
	}

	// Truncate or pad image name
	if len(imgName) > maxImgWidth {
		imgName = imgName[:maxImgWidth-3] + "..."
	} else {
		imgName = fmt.Sprintf("%-*s", maxImgWidth, imgName)
	}

	return fmt.Sprintf("%s | %s", archiveName, imgName)
}
