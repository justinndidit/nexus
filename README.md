# Nexus: Distributed Ledger System

Nexus is a polyglot financial system exploring how to move money correctly across service
boundaries: ACID transfers in the write path, eventual consistency in the read path, and
the transactional outbox pattern to bridge the two.

> ### ⚠️ Status: pre-alpha — architecture in place, no working end-to-end path
>
> This is a learning/portfolio project under active construction. The service layering,
> the schema design and the container tooling exist; the wiring between them does not yet.
> **No transfer has ever completed end to end.**
>
> - **What is verified working, what is broken, and how each was proven:**
>   [`docs/ASSESSMENT.md`](docs/ASSESSMENT.md)
> - **What happens next, in order:** [`docs/BUILD_PLAN.md`](docs/BUILD_PLAN.md)
>
> Sections below marked **planned** describe the intended design, not current behaviour.

---

## 📊 Component status

| Component | Language | Builds | Starts | Functional | Tests |
|---|---|:---:|:---:|:---:|:---:|
| [Ledger](services/ledger) | Go 1.25 | ✅ | ❌ stub `main()`, exits immediately | ❌ | 3 unit |
| [Account](services/account) | Java 25 / Spring Boot 4 | ✅ | ✅ | ⚠️ consumes events, misreads amounts by 100× | ❌ |
| [Notification](services/notification) | Python 3.13 / FastAPI | ✅ | ✅ | ❌ never subscribes to Kafka | ❌ |
| [User](services/user) | Go 1.25 | ⚠️ only inside `go.work` | ❌ stub `main()` | ❌ | 1 unit |
| Payment | — | — | — | — | empty directory |
| API Gateway | — | — | — | — | empty directory |

| Infrastructure | Status |
|---|---|
| Dev images (`docker compose build`) | ✅ all three build |
| Dev stack (`docker compose up`) | ⚠️ starts; app containers never reach `healthy` |
| Ledger database schema | ❌ migration is not valid SQL; has never been applied |
| Account database schema | ✅ Flyway migration applies |
| `infra/docker-compose.yml` (prod) | ❌ does not parse; a stub containing only Redis |
| CI | ❌ none |

---

## 🏗 Architecture

```
                   ┌─────────────┐
   HTTP ──────────▶│    User     │  identity, auth, KYC        (Go)
                   └──────┬──────┘
                          │ user.registered.v1
                          ▼
                   ┌─────────────┐
   HTTP/gRPC ─────▶│   Ledger    │  system of record          (Go)
                   │             │  ├ accounts
                   │             │  ├ transactions
                   │             │  ├ ledger_entries
                   │             │  └ outbox_events
                   └──────┬──────┘
                          │ relay worker (transactional outbox)
                          ▼
                   ┌─────────────┐
                   │    Kafka    │  ledger.transactions.v1
                   └──┬───────┬──┘
                      │       │
          ┌───────────▼─┐   ┌─▼──────────────┐
          │   Account   │   │  Notification  │
          │ read model  │   │  email/sms/push│
          │   (Java)    │   │    (Python)    │
          └─────────────┘   └────────────────┘
```

Each service owns its own PostgreSQL database. The ledger is the only authoritative source
of balances; the account service holds a read model derived from the event stream.

---

## 🚀 Key architectural ideas

These are the techniques the project exists to demonstrate. Each is marked with what is
actually true today.

### 1. Atomic balance updates — *partially implemented*

Balance changes happen in the database, not in application memory:

```sql
UPDATE accounts SET available_balance = available_balance + @amount WHERE id = @id;
```

This avoids read-modify-write races. The statement is written
(`services/ledger/internal/ledger/repo_postgres.go`), but the surrounding transfer has never
executed — the `accounts` table does not exist yet (`docs/ASSESSMENT.md` §4.2).

### 2. Transactional outbox — *partially implemented*

Events are written to an `outbox_events` table inside the same transaction as the balance
change, then relayed to Kafka by a polling worker using `FOR UPDATE SKIP LOCKED`. The
outbox write and the relay worker exist; the Kafka publish call is currently a stub that
returns success without publishing (`docs/ASSESSMENT.md` §4.7).

### 3. Deadlock prevention by lock ordering — *implemented, untested against a database*

Account IDs are sorted before rows are locked, eliminating circular wait. See
`utils.SortAccount` and `LedgerService.baseTransfer`. Unit-tested at the service layer;
no concurrency test against a real database yet.

### 4. Distributed idempotency — *planned*

Every transfer carries an `idempotency_key`. The intended behaviour is that a replayed key
returns the original transaction rather than performing a second transfer. The key is
currently stored but never checked before the balance mutation; the account service's
`processed_events` guard does deduplicate redelivered Kafka events
(`docs/ASSESSMENT.md` §4.8, §5.3).

### 5. Money as integer minor units — *inconsistent across services*

Amounts are `int64` in the currency's smallest denomination (kobo, cents) with an ISO-4217
currency code — never floating point. The Go services follow this; the Java service still
uses `BigDecimal` against `DECIMAL(19,4)` columns, so amounts crossing that boundary are
wrong by a factor of 100. This is the highest-priority open defect
(`docs/ASSESSMENT.md` §5.1).

---

## 📂 Project structure

```
├── api-gateway/           # empty — not yet designed
├── docs/
│   ├── ASSESSMENT.md      # verified current state
│   └── BUILD_PLAN.md      # phased plan to get to a working system
├── infra/
│   ├── docker-compose.dev.yml
│   ├── docker-compose.yml       # stub; does not parse
│   └── CONTAINERIZATION.md      # container analysis (§2 partly superseded)
├── proto/                 # protobuf contracts — currently invalid and unused
├── scripts/               # database init scripts
├── services/
│   ├── account/           # Java   — read model
│   ├── ledger/            # Go     — system of record
│   ├── notification/      # Python — side effects
│   ├── payment/           # empty
│   └── user/              # Go     — identity
├── build.gradle           # Gradle root (Java services)
├── go.work                # Go workspace (ledger, user)
└── Makefile
```

---

## 🛠 Tech stack

| Service | Language | Framework | Data | Migrations |
|---|---|---|---|---|
| Ledger | Go 1.25 | stdlib + pgx | PostgreSQL 16 | tern |
| Account | Java 25 | Spring Boot 4.0 | PostgreSQL 16 | Flyway |
| Notification | Python 3.13 | FastAPI + aiokafka | — | — |
| User | Go 1.25 | stdlib + pgx | PostgreSQL 16 (schema not yet written) | tern |

Shared: Apache Kafka (KRaft), Redis, Docker Compose.

---

## 🚦 Getting started

### Prerequisites

- Docker & Docker Compose v2
- Go 1.25+ · Java 25 · Python 3.13 (only needed to run a service outside a container)

### Running the dev stack

```bash
make dev-compose-build   # build all images
make dev-compose-up      # start the stack
```

**What to expect today**, so you can tell defects from your own setup:

- Postgres, Redis, Kafka and kafka-ui come up.
- On a **fresh volume**, `dev-postgres-ledger` exits with code 3 — `scripts/init_ledger_db.sql`
  collides with `POSTGRES_DB` (`docs/ASSESSMENT.md` §8.2). Fixing this is Phase 0.
- `dev-ledger-service` builds, prints two lines and exits 0.
- `dev-notification-service` serves `GET /` on :3002 but stays `unhealthy` — the
  healthcheck probes `/health`, which does not exist.
- `dev-accounts-service` boots Spring Boot, runs Flyway and serves :3000.
- Both Postgres containers publish fixed host ports (5432, 5433); if either is already
  bound on your machine, that container will fail to start.

### Running services individually

```bash
# Ledger (Go)         — currently a stub
cd services/ledger && go run ./cmd/server

# Account (Java)      — needs services/account/.env and a reachable Postgres
./gradlew :services:account:bootRun

# Notification (Python)
cd services/notification && uv sync && uv run uvicorn app.main:app --reload --port 3002

# User (Go)           — only builds inside the workspace
cd services/user && go build ./...
```

### Running tests

```bash
cd services/ledger && go test ./...
cd services/user   && go test ./...
./gradlew :services:account:test     # no tests defined yet
```

---

## 🤝 Contributing

Start with [`docs/BUILD_PLAN.md`](docs/BUILD_PLAN.md) — it lists the work in dependency
order, with the decisions (D1–D10) that shouldn't be re-litigated per-PR.

Four working agreements, each of which exists because of a real defect in this repo:

1. **A README claim requires a passing test.** No test, no claim — mark it *planned*.
2. **No stub returns success.** Unimplemented functions return an explicit error.
3. **Every interface with a production implementation gets `var _ Iface = (*Impl)(nil)`.**
4. **Every migration runs in CI against a real Postgres, up and down.**

---

## 📝 Blog posts & deep dives

Planned once there is a working system to write about.
