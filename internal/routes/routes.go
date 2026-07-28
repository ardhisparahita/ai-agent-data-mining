package routes

import (
	"github.com/ardhisparahita/ecommerce-api/internal/handler"
	"github.com/ardhisparahita/ecommerce-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, authHandler *handler.AuthHandler, categoryHandler *handler.CategoryHandler, productHandler *handler.ProductHandler, addressHandler *handler.AddressHandler, cartHandler *handler.CartHandler, checkoutHandler *handler.CheckoutHandler, orderHandler *handler.OrderHandler) {
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	users := api.Group("/users", middleware.JWT())
	users.Get("/profile", authHandler.GetProfile)
	users.Put("/profile", authHandler.UpdateProfile)
	users.Patch("/change-password", authHandler.ChangePassword)

	category := api.Group("/categories", middleware.JWT())
	category.Get("/", categoryHandler.FindAll)
	category.Get("/:id", categoryHandler.FindByID)
	category.Post("/", middleware.AdminOnly(), categoryHandler.Create)
	category.Put("/:id", middleware.AdminOnly(), categoryHandler.Update)

	product := api.Group("/products", middleware.JWT())
	product.Post("/", middleware.AdminOnly(), productHandler.Create)
	product.Get("/", productHandler.FindAll)
	product.Get("/:id", productHandler.FindByID)
	product.Put("/:id", middleware.AdminOnly(), productHandler.Update)
	product.Delete("/:id", middleware.AdminOnly(), productHandler.Delete)
	product.Post("/:id/image", middleware.AdminOnly(), productHandler.UploadImage)

	address := api.Group("/addresses", middleware.JWT())
	address.Post("/", addressHandler.Create)
	address.Get("/", addressHandler.FindAll)
	address.Get("/:id", addressHandler.FindByID)
	address.Put("/:id", addressHandler.Update)
	address.Delete("/:id", addressHandler.Delete)

	cart := api.Group("/carts", middleware.JWT())
	cart.Post("/", cartHandler.AddToCart)
	cart.Get("/", cartHandler.FindAll)
	cart.Put("/:id", cartHandler.Update)
	cart.Delete("/:id", cartHandler.Delete)

	checkout := api.Group("/checkouts", middleware.JWT())
	checkout.Post("/", checkoutHandler.Checkout)

	order := api.Group("/orders", middleware.JWT())
	order.Get("/", orderHandler.FindAll)
	order.Get("/:id", orderHandler.FindByID)
	order.Patch("/:id/pay", middleware.AdminOnly(), orderHandler.MarkAsPaid)
	order.Patch("/:id/fail", orderHandler.MarkAsFailed)
	order.Patch("/:id/cancel", orderHandler.Cancel)
	order.Patch("/:id/ship", middleware.AdminOnly(), orderHandler.MarkAsShipped)
	order.Patch("/:id/complete", orderHandler.MarkAsCompleted)

}
