package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/illusi03/golearn/module_product"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "User Management API",
		ServerHeader: "Fiber",
		ErrorHandler: customErrorHandler,
	})

	// Health check endpoint
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "User Management API is running",
			"version": "1.0.0",
		})
	})

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

	// User routes
	productHandler := module_product.NewProductHandler()
	products := v1.Group("/products")
	products.Post("/", productHandler.Create)
	products.Get("/", productHandler.GetAll)
	products.Put("/:id", productHandler.Update)
	products.Delete("/:id", productHandler.Delete)

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
