package wallet

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uint   `gorm:"primaryKey"`
	OwnerID   string `gorm:"uniqueIndex:idx_owner_currency"`
	AssetType string `gorm:"uniqueIndex:idx_owner_currency"`
	Balance   int64  `gorm:"default:0;check:balance >= 0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transaction struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ReferenceID string    `gorm:"uniqueIndex"`
	Type        string
	Status      string
	CreatedAt   time.Time
}

type LedgerEntry struct {
	ID            uint      `gorm:"primaryKey"`
	TransactionID uuid.UUID `gorm:"type:uuid;index"`
	AccountID     uint      `gorm:"index"`
	Amount        int64
	CreatedAt     time.Time
}

type IdempotencyKey struct {
	Key        string `gorm:"primaryKey"`
	Response   string
	StatusCode int
	CreatedAt  time.Time
}
