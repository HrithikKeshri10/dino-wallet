package wallet

import (
	"dino-wallet/internals/database"
	"dino-wallet/internals/dto"
	"dino-wallet/models/wallet"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ProcessCredit(req dto.CreditRequest, txType string) (int, interface{}, error) {
	db := database.Client()

	var existingKey wallet.IdempotencyKey
	if err := db.First(&existingKey, "key = ?", req.IdempotencyKey).Error; err == nil {
		return existingKey.StatusCode, map[string]string{"status": "success", "message": existingKey.Response}, nil
	}

	err := db.Transaction(func(tx *gorm.DB) error {

		var userAccount wallet.Account
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&userAccount, "owner_id = ? AND asset_type = ?", req.UserID, req.Currency).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				userAccount = wallet.Account{
					OwnerID:   req.UserID,
					AssetType: req.Currency,
					Balance:   0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				if err := tx.Create(&userAccount).Error; err != nil {
					return errors.New("failed to create new currency wallet")
				}
			} else {
				return errors.New("database error fetching user account")
			}
		}

		var systemAccount wallet.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&systemAccount, "owner_id = ? AND asset_type = ?", "SYSTEM_TREASURY", req.Currency).Error; err != nil {
			return errors.New("currency not supported (system treasury missing)")
		}

		txn := wallet.Transaction{
			ID:          uuid.New(),
			ReferenceID: req.IdempotencyKey,
			Type:        txType,
			Status:      "COMPLETED",
			CreatedAt:   time.Now(),
		}
		if err := tx.Create(&txn).Error; err != nil {
			return err
		}

		if err := tx.Create(&wallet.LedgerEntry{TransactionID: txn.ID, AccountID: userAccount.ID, Amount: req.Amount}).Error; err != nil {
			return err
		}
		if err := tx.Create(&wallet.LedgerEntry{TransactionID: txn.ID, AccountID: systemAccount.ID, Amount: -req.Amount}).Error; err != nil {
			return err
		}

		userAccount.Balance += req.Amount
		systemAccount.Balance -= req.Amount

		if err := tx.Save(&userAccount).Error; err != nil {
			return err
		}
		if err := tx.Save(&systemAccount).Error; err != nil {
			return err
		}

		if err := tx.Create(&wallet.IdempotencyKey{Key: req.IdempotencyKey, Response: txType + " successful", StatusCode: 200}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 500, nil, err
	}
	return 200, map[string]string{"status": "success", "message": txType + " successful"}, nil
}

func ProcessSpend(req dto.SpendRequest) (int, interface{}, error) {
	db := database.Client()

	var existingKey wallet.IdempotencyKey
	if err := db.First(&existingKey, "key = ?", req.IdempotencyKey).Error; err == nil {
		return existingKey.StatusCode, map[string]string{"status": "success", "message": existingKey.Response}, nil
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		var userAccount wallet.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&userAccount, "owner_id = ? AND asset_type = ?", req.UserID, req.Currency).Error; err != nil {
			return errors.New("user account not found or insufficient funds")
		}

		if userAccount.Balance < req.Amount {
			return errors.New("insufficient funds")
		}

		var systemAccount wallet.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&systemAccount, "owner_id = ? AND asset_type = ?", "SYSTEM_TREASURY", req.Currency).Error; err != nil {
			return errors.New("currency not supported")
		}

		txn := wallet.Transaction{ID: uuid.New(), ReferenceID: req.IdempotencyKey, Type: "SPEND", Status: "COMPLETED", CreatedAt: time.Now()}
		if err := tx.Create(&txn).Error; err != nil {
			return err
		}

		if err := tx.Create(&wallet.LedgerEntry{TransactionID: txn.ID, AccountID: userAccount.ID, Amount: -req.Amount}).Error; err != nil {
			return err
		}
		if err := tx.Create(&wallet.LedgerEntry{TransactionID: txn.ID, AccountID: systemAccount.ID, Amount: req.Amount}).Error; err != nil {
			return err
		}

		userAccount.Balance -= req.Amount
		systemAccount.Balance += req.Amount

		if err := tx.Save(&userAccount).Error; err != nil {
			return err
		}
		if err := tx.Save(&systemAccount).Error; err != nil {
			return err
		}

		if err := tx.Create(&wallet.IdempotencyKey{Key: req.IdempotencyKey, Response: "Purchase successful", StatusCode: 200}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err.Error() == "insufficient funds" || err.Error() == "user account not found or insufficient funds" {
			return 400, nil, err
		}
		return 500, nil, err
	}
	return 200, map[string]string{"status": "success", "message": "Purchase successful"}, nil
}

func GetBalance(userID string) ([]wallet.Account, error) {
	var accounts []wallet.Account
	err := database.Client().Find(&accounts, "owner_id = ?", userID).Error
	return accounts, err
}

func GetUserTransactions(userID string) ([]dto.TransactionHistoryResponse, error) {
	var results []dto.TransactionHistoryResponse

	err := database.Client().Table("transactions").
		Select("transactions.id, transactions.type, transactions.status, transactions.created_at, ledger_entries.amount, accounts.asset_type").
		Joins("JOIN ledger_entries ON ledger_entries.transaction_id = transactions.id").
		Joins("JOIN accounts ON ledger_entries.account_id = accounts.id").
		Where("accounts.owner_id = ?", userID).
		Order("transactions.created_at desc").
		Scan(&results).Error

	return results, err
}
