package core

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// WallpaperManager provides methods to manage wallpapers
type WallpaperManager struct {
	SavePath string
}

// SetLatestWallpaper sets the most recently downloaded image as desktop background
func (wm *WallpaperManager) SetLatestWallpaper() error {

	var latestFile string
	var latestTime time.Time

	images := GetImagesWithStat(wm.SavePath)
	for _, f := range images {
		if f.Info.ModTime().After(latestTime) {
			latestTime = f.Info.ModTime()
			latestFile = f.Path
		}
	}

	// If no image found
	if latestFile == "" {
		return fmt.Errorf("no images found in %s", wm.SavePath)
	}

	// Set the wallpaper
	return SetWallpaper(latestFile)
}

// Optional: Rotate through wallpapers
func (wm *WallpaperManager) SetRandomWallpapers() error {

	imageFiles := GetImagesList(wm.SavePath)

	if len(imageFiles) == 0 {
		return fmt.Errorf("no images found in %s", wm.SavePath)
	}

	// Randomly select an image
	rand.New(rand.NewSource(125))
	randomIndex := rand.Intn(len(imageFiles))

	return SetWallpaper(imageFiles[randomIndex])
}

// SetWallpaper sets the desktop background for different operating systems
func SetWallpaper(imagePath string) error {

	switch runtime.GOOS {
	case "windows":
		return setWindowsWallpaper(imagePath)
	case "darwin":
		return setMacWallpaper(imagePath)
	case "linux":
		return setLinuxWallpaper(imagePath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// setWindowsWallpaper sets the wallpaper on Windows
func setWindowsWallpaper(imagePath string) error {
	// Ensure the image path is absolute
	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		return err
	}

	if !fileExists(absPath) {
		return errors.New("file not exist")
	}
	// Use PowerShell to set wallpaper
	cmd := exec.Command("powershell",
		"-Command",
		fmt.Sprintf(`
    Add-Type -TypeDefinition @"
    using System;
    using System.Runtime.InteropServices;
    public class Wallpaper { 
        [DllImport("user32.dll", CharSet=CharSet.Auto)] 
        public static extern int SystemParametersInfo(int uAction, int uParam, string lpvParam, int fuWinIni); 
        public static void Set(string path) { 
            SystemParametersInfo(0x0014, 0, path, 0x01 | 0x02); 
        } 
    } 
"@ -Language CSharp; 
    [Wallpaper]::Set("%s")
    `, absPath))

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error setting wallpaper: %v\nOutput: %s", err, string(output))
	}

	return err
}

// setMacWallpaper sets the wallpaper on macOS
func setMacWallpaper(imagePath string) error {
	// Ensure the image path is absolute
	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		return err
	}

	// Use osascript to set wallpaper
	cmd := exec.Command("osascript", "-e",
		fmt.Sprintf(`tell application "System Events" to set picture of every desktop to "%s"`, absPath))

	return cmd.Run()
}

// setLinuxWallpaper sets the wallpaper on Linux (supports multiple desktop environments)
func setLinuxWallpaper(imagePath string) error {
	// Ensure the image path is absolute
	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		return err
	}

	// List of desktop environment wallpaper setting commands
	desktopCommands := [][]string{
		{"gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://" + absPath},
		{"plasma-desktop", "setWallpaper", absPath},
		{"feh", "--bg-fill", absPath},
		{"nitrogen", "--set-zoom-fill", absPath},
		{"xfconf-query", "-c", "xfce4-desktop", "-p", "/backdrop/screen0/monitor0/image-path", "-s", absPath},
	}

	// Try each command until one succeeds
	for _, cmdArgs := range desktopCommands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("could not set wallpaper on Linux")
}

// Optional: Utility function to check if file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
