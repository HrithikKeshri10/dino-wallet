package app

import (
	"dino-wallet/internals/database"
	"dino-wallet/internals/server"
	"dino-wallet/models/wallet"
	"log"
	"time"
)

func Setup() {
	database.Connect()

	log.Println("Running database migrations...")
	if err := database.Client().AutoMigrate(
		&wallet.Account{},
		&wallet.Transaction{},
		&wallet.LedgerEntry{},
		&wallet.IdempotencyKey{},
	); err != nil {
		log.Fatal("Migration failed:", err)
	}

	seedDatabase()

	app := server.Setup()

	log.Println("Server starting on port 3000...")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

func seedDatabase() {
	db := database.Client()
	var count int64
	db.Model(&wallet.Account{}).Count(&count)

	if count == 0 {
		log.Println("Database empty. Seeding initial data...")
		now := time.Now()

		treasury := []wallet.Account{
			{OwnerID: "SYSTEM_TREASURY", AssetType: "GOLD_COIN", Balance: 1000000000, CreatedAt: now, UpdatedAt: now},
			{OwnerID: "SYSTEM_TREASURY", AssetType: "DIAMOND", Balance: 1000000000, CreatedAt: now, UpdatedAt: now},
			{OwnerID: "SYSTEM_TREASURY", AssetType: "LOYALTY_POINT", Balance: 1000000000, CreatedAt: now, UpdatedAt: now},
		}
		db.Create(&treasury)

		users := []wallet.Account{
			{OwnerID: "USER_1", AssetType: "GOLD_COIN", Balance: 500, CreatedAt: now, UpdatedAt: now},
			{OwnerID: "USER_1", AssetType: "DIAMOND", Balance: 10, CreatedAt: now, UpdatedAt: now},
			{OwnerID: "USER_1", AssetType: "LOYALTY_POINT", Balance: 100, CreatedAt: now, UpdatedAt: now},
			{OwnerID: "USER_2", AssetType: "GOLD_COIN", Balance: 100, CreatedAt: now, UpdatedAt: now},
		}
		db.Create(&users)

		log.Println("Auto-seeding complete.")
	}
}
