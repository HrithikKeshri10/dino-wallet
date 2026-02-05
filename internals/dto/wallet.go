package dto

import (
	"dino-wallet/models/wallet"
)

type CreditRequest struct {
	UserID         string `json:"user_id"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	IdempotencyKey string `json:"idempotency_key"`
}

type SpendRequest struct {
	UserID         string `json:"user_id"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	IdempotencyKey string `json:"idempotency_key"`
}

type TransactionHistoryResponse struct {
	ID        string `json:"transaction_id"`
	Type      string `json:"type"`
	Amount    int64  `json:"amount"`
	AssetType string `json:"asset_type"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type UserTransactionsResponse struct {
	UserID       string                       `json:"user_id"`
	Transactions []TransactionHistoryResponse `json:"transactions"`
}

type UserBalanceResponse struct {
	UserID   string           `json:"user_id"`
	Balances []wallet.Account `json:"balances"`
}
