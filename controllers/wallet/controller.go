package wallet

import (
	"dino-wallet/internals/dto"
	"dino-wallet/services/wallet"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func validateRequest(userID, currency, idempotencyKey string, amount int64) string {
	if strings.TrimSpace(userID) == "" {
		return "user_id is required"
	}
	if strings.TrimSpace(currency) == "" {
		return "currency is required"
	}
	if strings.TrimSpace(idempotencyKey) == "" {
		return "idempotency_key is required"
	}
	if amount <= 0 {
		return "amount must be positive"
	}
	return ""
}

func TopUp(c *fiber.Ctx) error {
	var req dto.CreditRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	if errStr := validateRequest(req.UserID, req.Currency, req.IdempotencyKey, req.Amount); errStr != "" {
		return c.Status(400).JSON(fiber.Map{"error": errStr})
	}

	code, resp, err := wallet.ProcessCredit(req, "TOPUP")
	if err != nil {
		return c.Status(code).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(code).JSON(resp)
}

func Bonus(c *fiber.Ctx) error {
	var req dto.CreditRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	if errStr := validateRequest(req.UserID, req.Currency, req.IdempotencyKey, req.Amount); errStr != "" {
		return c.Status(400).JSON(fiber.Map{"error": errStr})
	}

	code, resp, err := wallet.ProcessCredit(req, "BONUS")
	if err != nil {
		return c.Status(code).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(code).JSON(resp)
}

func Spend(c *fiber.Ctx) error {
	var req dto.SpendRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	if errStr := validateRequest(req.UserID, req.Currency, req.IdempotencyKey, req.Amount); errStr != "" {
		return c.Status(400).JSON(fiber.Map{"error": errStr})
	}

	code, resp, err := wallet.ProcessSpend(req)
	if err != nil {
		return c.Status(code).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(code).JSON(resp)
}

func GetBalance(c *fiber.Ctx) error {
	userID := c.Params("id")
	if strings.TrimSpace(userID) == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user_id is required"})
	}

	accounts, err := wallet.GetBalance(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	if len(accounts) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "No wallets found for this user"})
	}

	response := dto.UserBalanceResponse{
		UserID:   userID,
		Balances: accounts,
	}

	return c.JSON(response)
}

func GetTransactions(c *fiber.Ctx) error {
	userID := c.Params("id")
	if strings.TrimSpace(userID) == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user_id is required"})
	}

	history, err := wallet.GetUserTransactions(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch transactions"})
	}

	response := dto.UserTransactionsResponse{
		UserID:       userID,
		Transactions: history,
	}

	return c.JSON(response)
}
