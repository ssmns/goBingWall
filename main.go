package main

import (
	"os"

	cmd "github.com/ssmns/goBingWall/cmd"
)

const baseURL = "https://bingwallpaper.anerg.com"

func main() {
	os.Setenv("BingWallPaper_URL", baseURL)
	cmd.Execute()
}
