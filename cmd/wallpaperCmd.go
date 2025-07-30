package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	core "github.com/ssmns/goBingWall/core"

	"github.com/spf13/cobra"
)

// Wallpaper Subcommand
var wallpaperCmd = &cobra.Command{
	Use:   "wallpaper",
	Short: "Manage wallpapers",
}

// Set Wallpaper Subcommand
var setWallpaperCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the latest downloaded wallpaper",
	Run: func(cmd *cobra.Command, args []string) {
		// Use default save path if not specified
		if savePath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}
			savePath = filepath.Join(homeDir, "BingWall")
		}

		wallpaperManager := &core.WallpaperManager{
			SavePath: savePath,
		}

		err := wallpaperManager.SetLatestWallpaper()
		if err != nil {
			log.Fatalf("Failed to set wallpaper: %v", err)
		}
		fmt.Println(core.Green("Wallpaper set successfully!"))
	},
}

// Rotate Wallpaper Subcommand
var rotateWallpaperCmd = &cobra.Command{
	Use:   "random",
	Short: "Rotate to a random wallpaper",
	Run: func(cmd *cobra.Command, args []string) {
		// Use default save path if not specified
		if savePath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}
			savePath = filepath.Join(homeDir, "BingWall")
		}

		wallpaperManager := &core.WallpaperManager{
			SavePath: savePath,
		}

		err := wallpaperManager.SetRandomWallpapers()
		if err != nil {
			log.Fatalf("Failed to rotate wallpaper: %v", err)
		}
		fmt.Println(core.Green("Wallpaper rotated successfully!"))
	},
}
