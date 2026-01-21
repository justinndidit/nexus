# Nexus Account Service (Java)

The Account Service provides the "User View" of the system. It maintains customer profiles and cached balances by listening to events emitted by the Ledger.

## 🛠 Features

- **Kafka Consumer:** Asynchronously listens to `ledger.transactions.v1` to update account states.
- **Idempotency Guard:** Uses a `processed_events` table within a JPA transaction to ensure exact-once processing of ledger events.
- **Optimistic Locking:** Utilizes Hibernate's `@Version` to handle concurrent API updates safely.

## 🏗 Tech Stack

- **Framework:** Spring Boot 4.0 (Java 25)
- **Data Access:** Spring Data JPA / Hibernate
- **Messaging:** Spring Kafka

## 📡 API Endpoints

- `GET /account/{id}`: Retrieves account details and current balance.
- `GET /account/transaction/{id}`: Retrieves transaction details for a specific movement.

## 🚀 Development

### Prerequisites

- Java 25
- Gradle 9.x

### Running Locally

```bash
./gradlew :services:account:bootRun