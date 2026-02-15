# Dino Wallet

Dino Wallet is a wallet service that tracks user credits and spending inside the application, ensuring data integrity through atomic and consistent transactions.

---

# How to Use

## Using Docker

Ensure Docker and Docker Compose are installed.

Clone the repository:

```sh
git clone https://github.com/Ahmed-Armaan/Dino-Wallet.git
cd Dino-Wallet
```

Run Docker Compose:

```sh
docker compose up
```

The application will start on port `8080`.

Base URL:

```
http://localhost:8080
```

---

## Build Manually

### 1. Start PostgreSQL

Ensure a PostgreSQL server is running.

### 2. Clone the Repository

```sh
git clone https://github.com/Ahmed-Armaan/Dino-Wallet.git
cd Dino-Wallet
```

### 3. Configure Environment Variables

Create a `.env` file:

```env
DATABASE_URL='<database connection string>'
PORT=<desired port> # defaults to 8080 if not set
```

Base URL:

```
http://localhost:<PORT>
```

If `PORT` is not set, it defaults to `8080`.

---

## Database Migration & Seeding

There are two options to migrate the schema and seed the database.

### Option 1 — GORM Migrator Script

1. Run the migrator:

```sh
go run cmd/migrator/.
```

2. Start the application:

```sh
go run .
```

---

### Option 2 — GORM AutoMigrate + SQL Seed

1. Start the server:

```sh
go run .
```

This runs `AutoMigrate` and builds the database schema automatically.

2. In a separate terminal, execute the seed script:

```sh
psql "$DATABASE_URL" -f seed.sql
```

Since the application is already running, it does not need to be started again.

---

# Concurrency & Idempotency

## Concurrency

All balance-changing operations execute inside database transactions.

To prevent race conditions, a PostgreSQL advisory transaction lock is acquired per user:

```sql
SELECT pg_advisory_xact_lock(hashtext(user_id))
```

Locks are acquired deterministically (per user), ensuring consistent ordering and preventing deadlocks.

This guarantees:

- Only one transaction can modify a user’s data at a time
- No double spending
- No lost updates

Retries handle temporary lock contention.

---

## Idempotency

Each state-changing operation requires a unique UUID `idempotency_key`.

The `idempotency_key` column is enforced as `UNIQUE` in the database.

If the same request is retried:

- Duplicate ledger transactions are rejected
- No duplicate balance updates occur

This provides safe retry behavior and exactly-once execution semantics.

---

# API Documentation

Base URL:

```
http://localhost:<PORT>
```

If `PORT` is not set, it defaults to `8080`.

---

## 1. Top Up

**POST** `/top_up`

Credits a user's account.

### Request Body

```json
{
  "user": "alice",
  "asset": "gold", 
  "amount": 100,
  "idempotency_key": "uuid"
}
```

Assets can be: `"gold"`, `"coins"`, or `"gem"`.

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
  "asset": "gold",
  "amount": 50,
  "idempotency_key": "uuid"
}
```

Assets can be: `"gold"`, `"coins"`, or `"gem"`.

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

- `400` → Missing user parameter
- `404` → User not found
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
- Concurrency is controlled using transactional advisory locks.
- Ledger entries follow double-entry accounting principles.
