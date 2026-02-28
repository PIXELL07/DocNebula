# 🚀 DocFlow — Resilient Asynchronous Document Processing Pipeline

DocFlow is a production-style, fault-tolerant document intelligence pipeline built in Go.  
It processes massive ZIP uploads asynchronously and safely through multiple stages:

- 📦 Unzip  
- 🔍 OCR  
- 🧠 Summarization (extensible)  
- 🔢 Vectorization  
- 💾 Metadata tracking  

The system is designed to **recover gracefully from worker crashes** using idempotent workers, retries, and state tracking in PostgreSQL.

---

## ✨ Key Features

- 🔄 Fully asynchronous queue-driven pipeline  
- 🛡️ Crash-safe recovery  
- ♻️ Idempotent job creation  
- 🔁 Exponential retry with DLQ  
- 💓 Worker heartbeat monitoring  
- 📈 Horizontally scalable workers  
- 🐳 Docker-first local development  
- 🧩 Clean Go project structure  

---

## 🏗️ High-Level Architecture

```text
User Upload
   ↓
Object Storage (S3 / MinIO)
   ↓
Orchestrator (Go)
   ↓
Queue A → Unzip Workers
   ↓
Queue B → OCR Workers
   ↓
Queue C → Vector Workers
   ↓
PostgreSQL (state machine + metadata)

---

## 🔄 Processing Flow

1. User uploads ZIP file

2. API creates idempotent job

3. Orchestrator enqueues unzip task

4. Workers process stages asynchronously

5. Each stage updates PostgreSQL

6. Failures retry automatically

7. Permanent failures go to DLQ


📁 Project Structure

docflow/
├── cmd/
│   ├── api/
│   ├── orchestrator/
│   ├── unzip-worker/
│   ├── ocr-worker/
│   └── vector-worker/
│
├── internal/
│   ├── config/
│   ├── db/
│   ├── models/
│   ├── repository/
│   ├── queue/
│   ├── orchestrator/
│   ├── workers/
│   ├── storage/
│   ├── heartbeat/
│   └── utils/
│
├── deployments/
├── scripts/
├── api/
└── README.md

🧪 Tech Stack

Language: Go 1.22+

Queue: Redis (local) / SQS (cloud-ready)

Database: PostgreSQL

Object Storage: MinIO (S3 compatible)

Containerization: Docker & Docker Compose

⚡ Quick Start (Local)

1️⃣ Start infrastructure
bash scripts/setup_local.sh

2️⃣ Start services (in separate terminals)
go run cmd/api/main.go
go run cmd/unzip-worker/main.go
go run cmd/ocr-worker/main.go
go run cmd/vector-worker/main.go

3️⃣ Submit a test job
bash scripts/create_test_job.sh


🔁 Failure Handling

DocFlow is built with production reliability patterns:

Automatic Retries

Exponential backoff

Max retry limit

At-least-once delivery