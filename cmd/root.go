package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	savePath     string
	resolution   string
	setWallpaper bool
	baseURL      string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "bingwallpaper",
	Short: "Bing Wallpaper Downloader",
	Run:   runScrapperCmd,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	baseURL = os.Getenv("BingWallPaper_URL")
	fmt.Println(baseURL)
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {

	// Add flags
	RootCmd.PersistentFlags().StringVarP(&savePath, "path", "p", "", "Path to save images")
	RootCmd.Flags().StringVarP(&resolution, "resolution", "r", "FHD", "Preferred resolution (all, 4K, 2K, FHD)")
	// RootCmd.Flags().BoolVarP(&setWallpaper, "set", "s", false, "Set a random image as wallpaper")
	// If baseURL was previously a global variable, you might want to add a flag for it
	RootCmd.Flags().StringVar(&baseURL, "url", os.Getenv("BingWallPaper_URL"), "Base URL for Bing wallpapers")

	// Add wallpaper subcommands
	wallpaperCmd.AddCommand(setWallpaperCmd)
	wallpaperCmd.AddCommand(rotateWallpaperCmd)
	RootCmd.AddCommand(wallpaperCmd)
}
