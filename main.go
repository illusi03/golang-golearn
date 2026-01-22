package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/illusi03/golearn/module_category"
	"github.com/illusi03/golearn/module_product"
	"github.com/illusi03/golearn/template"
	"github.com/illusi03/golearn/utils/middleware"
)

func main() {

	app := fiber.New(fiber.Config{
		AppName:      "API",
		ServerHeader: "Fiber",
		// ErrorHandler: customErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	// Health check endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "API is healthy",
			"status":  "up",
		})
	})

	// API routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Products routes
	productHandler := module_product.NewProductHandler()
	products := v1.Group("/products")
	products.Post("/", productHandler.Create)
	products.Get("/", productHandler.GetAll)
	products.Get("/:id", productHandler.GetDetail)
	products.Put("/:id", productHandler.Update)
	products.Delete("/:id", productHandler.Delete)

	// Categories routes
	categoryHandler := module_category.NewCategoryHandler()
	categorys := v1.Group("/categories")
	categorys.Post("/", categoryHandler.Create)
	categorys.Get("/", categoryHandler.GetAll)
	categorys.Get("/:id", categoryHandler.GetDetail)
	categorys.Put("/:id", categoryHandler.Update)
	categorys.Delete("/:id", categoryHandler.Delete)

	// Static files middleware (must be last)
	fsys := template.GetFileSystem("dist")
	app.Use(middleware.TryFilesHTML(fsys))
	app.Use(static.New("", static.Config{
		FS:         fsys,
		IndexNames: []string{"index.html"},
		Browse:     false,
		MaxAge:     5,
	}))

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	host := "0.0.0.0"
	port := 8080
	// Start server
	addr := fmt.Sprintf("%s:%s", host, strconv.Itoa(port))
	log.Printf("Server starting on %s", addr)

	if err := app.Listen(addr, fiber.ListenConfig{
		EnablePrefork: false,
	}); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Cleanup
	log.Println("Running cleanup tasks...")
	// Cleanup Logic ...
	log.Println("Server stopped gracefully")
}

// customErrorHandler handles errors globally
func customErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": err.Error(),
		"error":   nil,
	})
}
