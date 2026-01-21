# Nexus: High-Performance Distributed Ledger System

Nexus is a distributed, multi-language financial system designed for high-concurrency money movement. It demonstrates ACID-compliant transactions, Eventual Consistency, and the Transactional Outbox Pattern across a microservices architecture.

## 🏗 System Architecture

The system is built as a polyglot microservices monorepo:

  - Ledger Service (Go): The core "Engine of Truth." Handles atomic balance updates and maintains the immutable transaction history using PostgreSQL.

  - Account Service (Java/Spring Boot): The "Read Model." Consumes events from Kafka to maintain user-facing account views and profiles.

  - Message Broker (Kafka): The asynchronous bridge ensuring reliable delivery of events between services.

  - Relay Worker (Go): Implements the Transactional Outbox pattern to guarantee "At-Least-Once" delivery of financial events.

## 🚀 Key Architectural Features

1. 1. Atomic Financial Integrity

Instead of application-level math, Nexus performs balance updates directly in the database layer:

```
SQL

UPDATE accounts SET balance = balance + amount WHERE id = @id;

```
This prevents "lost updates" and race conditions during high-concurrency transfers.

2. Transactional Outbox Pattern

To ensure the Ledger and Kafka are always in sync, we use a local outbox_events table within the Ledger's ACID boundary. This guarantees that a notification is never sent if the transaction fails, and a transaction is never "lost" if the broker is down.

3. Distributed Idempotency

Every transaction carries a unique idempotency_key. The Java Account Service uses a processed_events table to ensure that even if Kafka delivers a message multiple times, the user's balance is only updated once.

4. Deadlock Prevention

Nexus implements a strict resource-ordering strategy. By sorting Account IDs before acquisition, we eliminate circular wait conditions in the database, ensuring system stability under heavy load.

---

## 📂 Project Structure

```
├── api-gateway
├── build
├── build.gradle
├── go.work
├── go.work.sum
├── gradle
├── gradle.properties
├── gradlew
├── gradlew.bat
├── infra
├── Makefile
├── proto
├── README.md
├── services
|   ├── account           #Java
|   ├── ledger            #Go
|   └── notification      #Python
|
└── settings.gradle
```
<!--
## 🛠 Tech Stack
________________________________________________________________
Service       |    Language  |   Framework     | Data Store     |
______________|______________|_________________|________________|
Ledger        |      Go      | Standard Library| pgx,PostgreSQL |
______________|______________|_________________|________________|
Account       |     Java     |  Spring Boot 4.0|  PostgreSQL    |
______________|______________|_________________|________________|
Messaging     | Apache Kafka |     -           |     -          |
________________________________________________________________| -->

## 🚦 Getting Started (Dev Environment)

### Prerequisites

  - Docker & Docker Compose

  - Go 1.23+

  - Java 25

### Running the system

Work in progress!

## 📝 Blog Posts & Deep Dives

Work in progress