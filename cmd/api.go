package cmd

import (
	"errors"
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
	config                *ApiConfig
	db                    *database.DB
	categoryService       *service.CategoryService
	categoryRepository    *repository.CategoryRepository
	productService        *service.ProductService
	productRepository     *repository.ProductRepository
	transactionService    *service.TransactionService
	transactionRepository *repository.TransactionRepository
	reportService         *service.ReportService
	reportRepository      *repository.ReportRepository
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
	a.productRepository = repository.NewProductRepository(a.db.Pool)
	a.transactionRepository = repository.NewTransactionRepository(a.db.Pool)
	a.reportRepository = repository.NewReportRepository(a.db.Pool)
}

func (a *ApiDeamon) registerService() {
	a.categoryService = service.NewCategoryService(a.categoryRepository)
	a.productService = service.NewProductService(a.productRepository)
	a.transactionService = service.NewTransactionService(a.transactionRepository, a.productRepository)
	a.reportService = service.NewReportService(a.reportRepository)
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
	productHandler := handler.NewProductHandler(a.productService)
	products := v1.Group("/products")
	products.Get("/", productHandler.GetAll)
	products.Get("/:id", productHandler.GetDetail)
	products.Post("/", productHandler.Create)
	products.Put("/:id", productHandler.Update)
	products.Delete("/:id", productHandler.Delete)

	// Categories routes
	categoryHandler := handler.NewCategoryHandler(a.categoryService)
	categories := v1.Group("/categories")
	categories.Get("/", categoryHandler.GetAll)
	categories.Get("/:id", categoryHandler.GetDetail)
	categories.Post("/", categoryHandler.Create)
	categories.Put("/:id", categoryHandler.Update)
	categories.Delete("/:id", categoryHandler.Delete)

	// Checkout/Transactions routes
	checkoutHandler := handler.NewCheckoutHandler(a.transactionService)
	transactions := v1.Group("/transactions")
	transactions.Get("/", checkoutHandler.GetAllTransactions)
	transactions.Post("/checkout", checkoutHandler.Checkout)

	// Report routes
	reportHandler := handler.NewReportHandler(a.reportService)
	report := api.Group("/report")
	report.Get("/hari-ini", reportHandler.GetTodayReport)
	report.Get("/", reportHandler.GetReport)

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
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": err.Error(),
		"error":   nil,
	})
}
