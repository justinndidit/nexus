# Nexus Ledger Service (Go)

The Ledger is the foundational service of the Nexus ecosystem. It is responsible for the immutable recording of all financial movements and ensuring the integrity of account balances.

## 🛠 Features

- **Atomic Transfers:** Uses PostgreSQL row-level locking and atomic SQL updates to ensure balance integrity.
- **Outbox Relay:** Implements a polling publisher that ships events from Postgres to Kafka using `FOR UPDATE SKIP LOCKED` for horizontal scalability.
- **Idempotent Operations:** Prevents duplicate transactions using unique idempotency keys at the database level.

## 🏗 Key Components

- `/internal/ledger`: Core domain logic and repository interfaces.
- `/internal/worker`: The Outbox Relay worker that polls the database for pending events.
- `/db/migrations`: SQL migration files (managed via Tern).

## 🗄 Database Schema

The ledger manages two primary tables:

1. `ledger_entries`: The immutable audit log of every credit/debit.
2. `outbox_events`: Pending notifications to be sent to the rest of the system.

## 🚀 Development

### Prerequisites

- Go 1.23+
- PostgreSQL 16+

### Running Locally

```bash
go run cmd/main.go