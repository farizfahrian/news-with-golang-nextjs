package app

import (
	"context"
	"log"
	"news-with-golang/config"
	"news-with-golang/internal/adapter/cloudflare"
	"news-with-golang/internal/adapter/handler"
	"news-with-golang/internal/adapter/repository"
	"news-with-golang/internal/core/service"
	"news-with-golang/lib/auth"
	"news-with-golang/lib/middleware"
	"news-with-golang/lib/pagination"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RunServer() {
	cfg := config.NewConfig()
	db, err := cfg.ConnectionPostgres()

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}

	err = os.MkdirAll("./temp/content", 0755)
	if err != nil {
		log.Fatalf("Error creating temp directory: %v", err)
		return
	}

	// Cloudflare R2
	cfdR2 := cfg.LoadAwsConfig()
	s3Client := s3.NewFromConfig(cfdR2)
	r2Adapter := cloudflare.NewCloudflareR2Adapter(s3Client, cfg)

	// JWT
	jwt := auth.NewJwt(cfg)

	// Middleware
	middlewareAuth := middleware.NewMiddleware(cfg)

	// Pagination
	_ = pagination.NewPagination()

	authRepo := repository.NewAuthRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	contentRepo := repository.NewContentRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)

	// Service
	authService := service.NewAuthService(authRepo, cfg, jwt)
	categoryService := service.NewCategoryService(categoryRepo)
	contentService := service.NewContentService(contentRepo, cfg, r2Adapter)
	userService := service.NewUserService(userRepo)

	// Handler
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	contentHandler := handler.NewContentHandler(contentService)
	userHandler := handler.NewUserHandler(userService)

	//initialization server with Fiber
	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} ${ip} ${status} - ${latency} ${method} ${path}\n",
	}))

	if os.Getenv("ENV") != "production" {
		cfg := swagger.Config{
			BasePath: "/api",
			FilePath: "./docs/swagger.json",
			Path:     "docs",
			Title:    "News API",
		}
		app.Use(swagger.New(cfg))
	}

	api := app.Group("/api")
	api.Post("/login", authHandler.Login)

	adminApi := api.Group("/admin")
	adminApi.Use(middlewareAuth.CheckToken())

	// Category
	categoryApi := adminApi.Group("/categories")
	categoryApi.Get("/", categoryHandler.GetCategories)
	categoryApi.Post("/", categoryHandler.CreateCategory)
	categoryApi.Get("/:categoryID", categoryHandler.GetCategoryById)
	categoryApi.Put("/:categoryID", categoryHandler.EditCategoryById)
	categoryApi.Delete("/:categoryID", categoryHandler.DeleteCategory)

	// Content
	contentApi := adminApi.Group("/contents")
	contentApi.Get("/", contentHandler.GetContents)
	contentApi.Post("/", contentHandler.CreateContent)
	contentApi.Get("/:contentID", contentHandler.GetContentById)
	contentApi.Put("/:contentID", contentHandler.UpdateContent)
	contentApi.Delete("/:contentID", contentHandler.DeleteContent)
	contentApi.Post("/upload-image", contentHandler.UploadImageR2)

	// User
	userApi := adminApi.Group("/users")
	userApi.Get("/profile", userHandler.GetUserByID)
	userApi.Put("/update-password", userHandler.UpdatePassword)

	// FE
	feApi := api.Group("/fe")
	feApi.Get("/categories", categoryHandler.GetCategoryFE)
	feApi.Get("/contents", contentHandler.GetContentWithQuery)
	feApi.Get("/contents/:contentID", contentHandler.GetContentDetail)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := app.Listen(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatalf("Error running server %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit

	log.Println("server shutdown of 5 seconds")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}
