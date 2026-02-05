package routes

import (
	"dino-wallet/controllers/wallet"

	"github.com/gofiber/fiber/v2"
)

func WalletRoutes(router fiber.Router) {
	w := router.Group("/wallet")

	w.Post("/topup", wallet.TopUp)
	w.Post("/bonus", wallet.Bonus)
	w.Post("/spend", wallet.Spend)
	w.Get("/balance/:id", wallet.GetBalance)
	w.Get("/transactions/:id", wallet.GetTransactions)
	w.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "wallet",
		})
	})

}
