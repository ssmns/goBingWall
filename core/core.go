package core

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

// Color-coded logging functions for beautiful and informative logs
var (
	Green   = color.New(color.FgGreen, color.Bold).SprintFunc()
	Yellow  = color.New(color.FgYellow, color.Bold).SprintFunc()
	Red     = color.New(color.FgRed, color.Bold).SprintFunc()
	Magenta = color.New(color.FgMagenta, color.Bold).SprintFunc()
)

// VisitURL fetches and parses an HTML document from the given URL
//
// Parameters:
//   - url: The web address to fetch and parse
//
// Returns:
//   - A *goquery.Document representing the parsed HTML document
//
// Behavior:
//   - Sends an HTTP GET request to the specified URL
//   - Parses the response body into a goquery Document
//   - Logs a fatal error if the request or parsing fails
//
// Example:
//
//	doc := visitURL("https://example.com")
func VisitURL(url string) *goquery.Document {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

// SaveText writes the given text content to a file at the specified path
//
// Parameters:
//   - content: The text to be written to the file
//   - path: The file path where the content will be saved
//
// Returns:
//   - An error if the file writing fails, nil otherwise
//
// Permissions:
//   - Creates file with read/write permissions for the owner (0644)
//
// Example:
//
//	err := saveText("Hello, World!", "./output.txt")
func SaveText(content, path string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// SaveImage downloads an image from the given URL and saves it to the specified path
//
// Parameters:
//   - url: The web address of the image to download
//   - path: The file path where the image will be saved
//
// Returns:
//   - An error if downloading or saving the image fails, nil otherwise
//
// Behavior:
//   - Sends an HTTP GET request to fetch the image
//   - Creates a new file at the specified path
//   - Copies the image data to the file
//
// Example:
//   err := saveImage("https://example.com/image.jpg", "./downloaded_image.jpg")

func SaveImage(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// GetCurrentImages retrieves a map of existing image filenames in a specified directory
//
// Parameters:
//   - savePath: The directory path where images are stored
//
// Returns:
//   - A map where keys are the base names of image files (before the first "-")
//     and values are always true, indicating the image exists
//
// Behavior:
//   - Searches for all .jpg files in the specified directory
//   - Extracts the base name of each file (part before the first "-")
//   - Creates a map with these base names as keys
//
// Error Handling:
//   - If there's an error finding files, returns an empty map
//
// Example:
//
//	// Finds all .jpg files in "./downloads" directory
//	existingImages := getCurrentImages("./downloads")
//	// Check if a specific image exists
//	if _, exists := existingImages["image123"]; exists {
//	    fmt.Println("Image already downloaded")
//	}
//
// Notes:
//   - Only considers .jpg files
//   - Uses filepath.Base to extract filename
//   - Splits filename by "-" to get the base name
func GetCurrentImages(savePath string) map[string]bool {
	currentImages := make(map[string]bool)
	matches := GetImagesList(savePath)
	for _, match := range matches {
		filename := filepath.Base(match)
		name := strings.Split(filename, "-")[0]
		currentImages[name] = true
	}
	return currentImages
}

func GetImagesList(savePath string) []string {
	var images []string
	for _, ext := range ImageExtensions {
		matches, err := filepath.Glob(filepath.Join(savePath, "*"+ext))
		if err != nil {
			continue
		}
		images = append(images, matches...)
	}
	return images
}

func GetImagesWithStat(savePath string) []ImageFullInfo {
	var images []ImageFullInfo
	imgs := GetImagesList(savePath)
	for _, f := range imgs {
		fileInfo, err := os.Stat(f)
		if err != nil {
			continue
		}
		fabs, err := filepath.Abs(f)
		if err != nil {
			continue
		}
		images = append(images, ImageFullInfo{Path: fabs, Info: fileInfo})
	}
	return images
}
