package main

import (
	"log"
	"os"
	"time"

	_ "github.com/ardhisparahita/ecommerce-api/docs"
	"github.com/ardhisparahita/ecommerce-api/internal/handler"
	"github.com/ardhisparahita/ecommerce-api/internal/repository"
	"github.com/ardhisparahita/ecommerce-api/internal/routes"
	"github.com/ardhisparahita/ecommerce-api/internal/seeders"
	"github.com/ardhisparahita/ecommerce-api/internal/service"
	"github.com/ardhisparahita/ecommerce-api/pkg/config"
	"github.com/ardhisparahita/ecommerce-api/pkg/database"
	"github.com/ardhisparahita/ecommerce-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// @title Ecommerce API
// @version 1.0
// @description Ecommerce Backend API using Golang Fiber
// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @tag.name Auth
// @tag.description Authentication APIs

// @tag.name Categories
// @tag.description Category APIs

// @tag.name Products
// @tag.description Product APIs

// @tag.name Addresses
// @tag.description Address APIs

// @tag.name Carts
// @tag.description Shopping Cart APIs

// @tag.name Checkout
// @tag.description Checkout APIs

// @tag.name Orders
// @tag.description Order APIs

func main() {
	config.LoadEnv()

	if err := os.MkdirAll("./uploads/products", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	var (
		db  *gorm.DB
		err error
	)

	const (
		maxRetries = 60
		retryDelay = 5 * time.Second
	)

	for i := 1; i <= maxRetries; i++ {

		db, err = database.Connect()

		if err == nil {

			sqlDB, errPing := db.DB()
			if errPing == nil {

				if errPing = sqlDB.Ping(); errPing == nil {
					log.Println("Database connected successfully")

					if err := seeders.SeedAdmin(db); err != nil {
						log.Fatalf("failed to seed admin: %v", err)
					}
					break
				}

				err = errPing
			}
		}

		log.Printf("Waiting for database... (%d/%d)", i, maxRetries)
		log.Printf("Database Error: %v", err)

		time.Sleep(retryDelay)
	}

	if err != nil {
		log.Fatalf("Failed to connect database after %d attempts: %v", maxRetries, err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
	})

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	orderItemRepo := repository.NewOrderItemRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	userService := service.NewAuthService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo)
	addressService := service.NewAddressService(addressRepo)
	cartService := service.NewCartService(cartRepo, productRepo)

	checkoutService := service.NewCheckoutService(
		db,
		cartRepo,
		productRepo,
		addressRepo,
		orderRepo,
		orderItemRepo,
		paymentRepo,
	)

	orderService := service.NewOrderService(
		db,
		orderRepo,
		productRepo,
		paymentRepo,
	)

	userHandler := handler.NewAuthHandler(userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	addressHandler := handler.NewAddressHandler(addressService)
	cartHandler := handler.NewCartHandler(cartService)
	checkoutHandler := handler.NewCheckoutHandler(checkoutService)
	orderHandler := handler.NewOrderHandler(orderService)

	app.Static("/uploads", "./uploads")

	routes.SetupRoutes(
		app,
		userHandler,
		categoryHandler,
		productHandler,
		addressHandler,
		cartHandler,
		checkoutHandler,
		orderHandler,
	)

	port := config.Get("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on port %s", port)

	log.Fatal(app.Listen(":" + port))
}
