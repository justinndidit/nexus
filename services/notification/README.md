# Nexus Notification Service (Python)

A lightweight, asynchronous service responsible for user communications (SMS, Email, Push) triggered by financial events.

## 🛠 Features

- **Async Consumer:** Built with `FastAPI` and `aiokafka` for high-throughput, non-blocking event processing.
- **Side-Effect Management:** Decouples notification logic from core financial transactions.
- **Scalability:** Designed to scale horizontally to handle notification bursts during peak transaction windows.

## 🏗 Tech Stack

- **Language:** Python 3.12+
- **Framework:** FastAPI
- **Library:** AIOKafka (Async Kafka Client)

## 📡 Functionality

1. Listens to `ledger.transactions.v1`.
2. Identifies the notification type based on event metadata.
3. Dispatches to external providers (Mocked for local development).

## 🚀 Development

### Prerequisites

- Python 3.12+
- Virtualenv or Conda

### Setup

```bash
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
uvicorn main:app --reload