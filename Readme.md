
# BingWallpaper Downloader For All OS


## Install and Build

Clone the project and install dependencies:

```sh
 git clone ....
 cd BingwallpalerGo
 go mod tidy 
```

For development:

```sh
go mod run
```

For building:

```sh
mkdir build
go mod build -o build/Bingwallpaper.exe
```

Copy `BingWallpaper.exe` to a directory in your system's PATH.


## Usage

### Download Bing Wallpapers

- Download Bing wallpapers from [anerg.com](https://bingwallpaper.anerg.com:)

```sh
Bingwallpaper.exe -p "D:\BingWall" -r FHD
```

Resolution (`-r`) options include `4K`, `2K`, and `FHD`, defaulting to `all`.

### Set Wallpaper

- Set the latest image from a directory as background:

```sh
Bingwallpaper.exe -p "D:\BingWall" wallpaper set
```

- Select a random image from a directory as background:

```sh
Bingwallpaper.exe -p "D:\BingWall" wallpaper random
```

#### Notes

- Ensure you have the necessary permissions to set wallpapers
- The `-p` flag specifies the directory to save or select images from
