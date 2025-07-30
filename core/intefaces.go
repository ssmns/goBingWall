package core

import "os"

type ImageFullInfo struct {
	Path string
	Info os.FileInfo
}

type Image struct {
	Detail string
	Name   string
	URL    string
}

type ImageDownload struct {
	Name string
	URLs map[string]string
}

type Archives []struct {
	Name string
	URL  string
}
