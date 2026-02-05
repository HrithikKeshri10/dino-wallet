package server

import (
	"dino-wallet/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Setup() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Dino Wallet Service",
	})

	app.Use(logger.New())
	app.Use(recover.New())

	setupRoutes(app)

	return app
}

func setupRoutes(app *fiber.App) {
	api := app.Group("/api")
	routes.WalletRoutes(api)
}
