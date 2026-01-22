package middleware

import (
	"io/fs"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func TryFilesHTML(fsys fs.FS) fiber.Handler {
	return func(c fiber.Ctx) error {
		p := c.Path()

		if strings.HasPrefix(p, "/api") || p == "/health" {
			return c.Next()
		}

		return c.Next()
	}
}
