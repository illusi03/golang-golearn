package template

import (
	"embed"
	"io/fs"
)

//go:embed dist
var embeddedFiles embed.FS

func GetFileSystem(path string) fs.FS {
	fs, err := fs.Sub(embeddedFiles, path)
	if err != nil {
		panic(err)
	}
	return fs
}
