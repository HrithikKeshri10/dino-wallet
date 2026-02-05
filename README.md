# Dino Wallet Service

Internal wallet service for managing virtual currencies (Gold Coins, Diamonds, Loyalty Points) in gaming/rewards platforms.

## Tech Stack

- **Go + Fiber** — Fast, lightweight, and I'm comfortable with Go's concurrency model
- **PostgreSQL** — Needed ACID transactions and row-level locking for a financial system
- **GORM** — Makes working with transactions and locking straightforward
- **Docker** — Easy setup, no "works on my machine" issues

## Getting Started

```bash
# Start everything
docker-compose up --build

# App runs on http://localhost:3000
```

That's it. The app auto-creates the database, runs migrations, and seeds test data.

### Manual seed (if needed)

```bash
psql -h localhost -U manager -d wallet_db -f seed.sql
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/wallet/topup` | Credit wallet (simulates purchase) |
| POST | `/api/wallet/bonus` | Issue free credits |
| POST | `/api/wallet/spend` | Deduct credits |
| GET | `/api/wallet/balance/:id` | Get user balances |
| GET | `/api/wallet/transactions/:id` | Get transaction history |

All POST endpoints require: `user_id`, `amount`, `currency`, `idempotency_key`

### Example Request

POST `http://localhost:3000/api/wallet/topup`

```json
{
  "user_id": "USER_1",
  "amount": 100,
  "currency": "GOLD_COIN",
  "idempotency_key": "unique-key-123"
}
```

## Concurrency Handling

Used pessimistic locking (`SELECT FOR UPDATE`) — when a transaction reads a wallet row, it locks it until the transaction completes. Simple, guarantees consistency, and fits the "never lose money" requirement.

```go
tx.Clauses(clause.Locking{Strength: "UPDATE"}).
    First(&account, "owner_id = ? AND asset_type = ?", userID, currency)
```

Deadlocks are avoided by always locking accounts in a consistent order (user account first, treasury second) across all flows. If a future flow needed the reverse order, the fix would be to sort accounts by ID before acquiring locks.

## Notes

- **Idempotency**: Every transaction stores its idempotency key, so retries return the same response without double-processing
- **Ledger entries**: Instead of just updating balances, every transaction creates debit/credit entries for audit trails
- **Balance protection**: A database-level `CHECK (balance >= 0)` constraint prevents negative balances even if application logic has a bug
- **Seeded data**: Treasury starts with 1B of each currency; USER_1 and USER_2 have small balances for testing
