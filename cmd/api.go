package cmd

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
	"github.com/illusi03/golearn/internal/database"
	"github.com/illusi03/golearn/internal/handler"
	"github.com/illusi03/golearn/internal/repository"
	"github.com/illusi03/golearn/internal/service"
	"github.com/illusi03/golearn/template"
	"github.com/illusi03/golearn/utils/middleware"
)

type ApiDeamon struct {
	config             *ApiConfig
	db                 *database.DB
	categoryService    *service.CategoryService
	categoryRepository *repository.CategoryRepository
}

type ApiConfig struct {
	Port   int    `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func InitApi(config *ApiConfig) *ApiDeamon {
	api := &ApiDeamon{
		config: config,
	}
	api.runApi()
	return api
}

func (a *ApiDeamon) runApi() {
	a.registerProvider()
}

func (a *ApiDeamon) registerProvider() {
	a.registerDb()
	a.registerRepository()
	a.registerService()
	a.registerHandler()
}

func (a *ApiDeamon) registerDb() {
	db, err := database.NewDatabaseConnection(a.config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	a.db = db
}

func (a *ApiDeamon) registerRepository() {
	a.categoryRepository = repository.NewCategoryRepository(a.db.Pool)
}

func (a *ApiDeamon) registerService() {
	a.categoryService = service.NewCategoryService(a.categoryRepository)
}

func (a *ApiDeamon) registerHandler() {
	config := a.config
	host := "0.0.0.0"
	port := config.Port

	app := fiber.New(fiber.Config{
		AppName:      "API",
		ServerHeader: "Fiber",
		ErrorHandler: customErrorHandler,
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
	// productHandler := module_product.NewProductHandler()
	// products := v1.Group("/products")
	// products.Post("/", productHandler.Create)
	// products.Get("/", productHandler.GetAll)
	// products.Get("/:id", productHandler.GetDetail)
	// products.Put("/:id", productHandler.Update)
	// products.Delete("/:id", productHandler.Delete)

	// Categories routes
	categoryHandler := handler.NewCategoryHandler(a.categoryService)
	categories := v1.Group("/categories")
	categories.Get("/", categoryHandler.GetAll)
	categories.Get("/:id", categoryHandler.GetDetail)
	categories.Post("/", categoryHandler.Create)
	categories.Put("/:id", categoryHandler.Update)
	categories.Delete("/:id", categoryHandler.Delete)

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
		a.db.Close()
	}()

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
