# Dino Wallet

Dino Wallet is a wallet service that tracks user credits and spending inside the application, ensuring data integrity through atomic and consistent transactions.

---

## How to Use

## Build Manually

- Start a PostgreSQL database server.
- Clone the repository:

```sh
git clone https://github.com/Ahmed-Armaan/Dino-Wallet.git
cd Dino-Wallet
```

- Create a `.env` file and configure the environment variables:

```env
DATABASE_URL='<database connection string>'
PORT=<desired port to run the server> # defaults to 8080
```

- Migrate the schema, seed the database, and run the application.

There are two options to perform migration and seeding:

---

## Option 1 — GORM Migrator Script

1. Run the GORM migrator:

```sh
go run cmd/migrator/.
```

2. After migration completes, start the application:

```sh
go run .
```

---

## Option 2 — GORM AutoMigrate + SQL Seed

1. Start the server:

```sh
go run .
```

This executes `AutoMigrate`, which builds the database schema automatically.

2. In a separate terminal, run the seed file manually:

```sh
psql "$DATABASE_URL" -f seed.sql
```

Since the application is already running, it does not need to be started again.

---

## Concurrency & Idempotency

### Concurrency

All balance-changing operations run inside database transactions.

To prevent race conditions, a PostgreSQL advisory transaction lock is acquired per user:

```sql
SELECT pg_advisory_xact_lock(hashtext(user_id))
```

Locks are acquired in a deterministic order (per user), ensuring consistent lock ordering and preventing deadlocks.

This guarantees:

- Only one transaction can modify a user’s associated data at a time
- No double spending
- No lost updates

Retries are performed to handle temporary lock contention.

---

### Idempotency

Each operation includes a unique UUID idempotency key.

The `idempotency_key` column is enforced as `UNIQUE` in the database.

If the same request is retried with the same key:

- Duplicate ledger transactions are rejected
- No duplicate balance updates occur

This provides safe retry behavior and exactly-once execution semantics.

---

# API Documentation

Base URL: `http://localhost:<PORT>` if .env is empty it is `http://localhost:8080`

---

## 1. Top Up

**POST** `/top_up`

Credits a user's account.

### Request Body

```json
{
  "user": "alice",
  "asset": "gold", // assets ca be "gold", "coins" or "gems"
  "amount": 100,
  "idempotency_key": "uuid"
}
```

### Response

- `200 OK` → Success
- `500` → Transaction failed

---

## 2. Random Bonus

**POST** `/bonus`

Grants a random bonus to the user.

### Request Body

```json
{
  "user": "alice",
  "idempotency_key": "uuid"
}
```

### Response

```json
{
  "status": "bonus processed"
}
```

- `400` → Invalid request
- `500` → Database error

---

## 3. Purchase

**POST** `/purchase`

Debits a user's account.

### Request Body

```json
{
  "user": "alice",
  "asset": "gold", // assets ca be "gold", "coins" or "gems"
  "amount": 50,
  "idempotency_key": "uuid"
}
```

### Response

```json
{
  "status": "purchase successful"
}
```

- `400` → Invalid request / insufficient balance
- `500` → Database error

---

## 4. Get Balance

**GET** `/balance?user=alice`

### Response

```json
{
  "balance": [
    {
      "asset": "gold",
      "balance": 100
    }
  ]
}
```

- `500` → Database error

---

## 5. Get Ledger

**GET** `/ledger?pageSize=10&pageToken=0`

### Query Parameters

| Parameter  | Type | Description |
|------------|------|------------|
| pageSize   | int  | Number of records per page |
| pageToken  | int  | Page index (0-based) |

### Response

```json
[
  {
    "id": "...",
    "account": "normal",
    "user": "alice",
    "amount": 100,
    "created_at": "timestamp"
  }
]
```

- `400` → Invalid pagination parameters
- `500` → Database error

---

## Notes

- All state-changing operations require a unique `idempotency_key` (UUID).
- Concurrency is controlled via transactional advisory locks.
- Ledger entries follow double-entry accounting principles.


