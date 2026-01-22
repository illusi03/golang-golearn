package middleware

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func TryFilesHTML(fsys fs.FS) fiber.Handler {
	return func(c fiber.Ctx) error {
		p := c.Path()
		fmt.Printf("[TryFilesHTML] Incoming path: %s\n", p)

		if strings.HasPrefix(p, "/api") || p == "/health" {
			return c.Next()
		}

		if strings.Contains(filepath.Base(p), ".") {
			return c.Next()
		}

		cleanPath := filepath.Clean(strings.TrimPrefix(p, "/"))

		if p != "/" && cleanPath != "." {
			testPath := cleanPath + ".html"
			fmt.Printf("[TryFilesHTML] Testing path: %s\n", testPath)

			if fileInfo, err := fs.Stat(fsys, testPath); err == nil && !fileInfo.IsDir() {
				c.Path("/" + testPath)
			} else {
				indexPath := filepath.Join(cleanPath, "index.html")
				if fileInfo, err := fs.Stat(fsys, indexPath); err == nil && !fileInfo.IsDir() {
					c.Path("/" + indexPath)
				}
			}
		}

		return c.Next()
	}
}
