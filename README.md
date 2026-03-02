# 🚀 DocNebula — Resilient Asynchronous Document Processing Pipeline

DocNebula is a production-style, fault-tolerant document processing pipeline built in Go.  
It processes large ZIP uploads asynchronously through multiple stages (Unzip → OCR → Vectorize) with strong guarantees around retries, idempotency, and worker crash recovery.

This project demonstrates real-world distributed systems patterns used in modern AI and data platforms.

---

## ✨ Key Features

- 🔄 Fully asynchronous queue-driven pipeline  
- 🛡️ Crash-safe processing using visibility timeout pattern  
- ♻️ Idempotent job creation (retry-safe API)  
- 🔁 Exponential retry with Dead Letter Queue (DLQ)  
- 💓 Worker heartbeat monitoring  
- 📈 Horizontally scalable worker pool  
- 🧠 Per-page OCR checkpointing model  
- 🐳 Docker-first local development  
- 📊 Structured JSON logging (`slog`)  
- ⚡ Memory-safe streaming unzip  

---

## 🏗️ High-Level Architecture

flowchart LR
    A[Client Upload] --> B[API Service]
    B --> C[Orchestrator]

    C --> Q1[Queue: Unzip]
    C --> Q2[Queue: OCR]
    C --> Q3[Queue: Vector]

    Q1 --> W1[Unzip Workers]
    Q2 --> W2[OCR Workers]
    Q3 --> W3[Vector Workers]

    W1 --> DB[(PostgreSQL)]
    W2 --> DB
    W3 --> DB

    W1 --> HB[(Worker Heartbeats)]
    W2 --> HB
    W3 --> HB

    Q1 --> DLQ[Dead Letter Queue]
    Q2 --> DLQ
    Q3 --> DLQ
````

---

## 🔄 Processing Flow

1. Client calls **POST /upload**
2. API generates or accepts **Idempotency-Key**
3. Job is created in PostgreSQL
4. Message is pushed to **unzip_queue**
5. Workers process stages asynchronously
6. Each stage updates database state
7. Failures retry automatically
8. Permanent failures go to DLQ

---

## 🧪 Tech Stack

* **Language:** Go 1.22+
* **Queue:** Redis (visibility-timeout pattern)
* **Database:** PostgreSQL
* **Storage:** Local / MinIO compatible
* **Containerization:** Docker & Docker Compose
* **Logging:** slog (structured JSON)

---

## 📁 Project Structure

docflow/
├── cmd/
│   ├── api/                 # HTTP API service
│   ├── unzip-worker/        # Unzip stage workers
│   ├── ocr-worker/          # OCR stage workers
│   └── vector-worker/       # Vector stage workers
│
├── internal/
│   ├── config/              # Environment configuration
│   ├── db/                  # Database connection & migrations
│   │   └── migrations/
│   ├── models/              # Domain models
│   ├── repository/          # Data access layer
│   ├── queue/               # Redis producer/consumer
│   ├── workers/
│   │   └── unzip/           # Streaming unzip logic
│   ├── storage/             # Object storage (MinIO/S3 ready)
│   ├── heartbeat/           # Worker heartbeat tracking
│   ├── http/                # Health & readiness endpoints
│   └── utils/               # Idempotency utilities
│   |___ orchestrator/
|
├── deployments/             # Docker & infra configs
├── scripts/                 # Helper scripts
├── api/                     # OpenAPI spec
├── go.mod
└── README.md

---

## ⚡ Quick Start (Local)

### 1️⃣ Start infrastructure

```bash
docker compose up -d
```

---

### 2️⃣ Run services (separate terminals)

```bash
go run cmd/api/main.go
go run cmd/unzip-worker/main.go
go run cmd/ocr-worker/main.go
go run cmd/vector-worker/main.go
```

---

### 3️⃣ Create a job

```bash
curl -X POST http://localhost:8080/upload \
  -H "Idempotency-Key: test-123"
```

Retrying with the same key will return the **same job**.

---

## 🛡️ Reliability Guarantees

### ✅ Idempotent API

Safe client retries:

```
same Idempotency-Key → same job
```

Prevents duplicate processing.

---

### ✅ Visibility Timeout Queue

Uses Redis `BRPOPLPUSH` pattern:

* worker crash → message not lost
* message stays in processing queue
* safe retry behavior

---

### ✅ Dead Letter Queue (DLQ)

Messages exceeding retry limit are moved to DLQ to prevent infinite retry loops.

---

### ✅ Per-Page Checkpointing

OCR progress is tracked at page level, enabling:

* partial recovery
* crash-safe resume
* large document support

---

### ✅ Memory-Safe Unzip

Streaming unzip prevents loading entire archives into memory, enabling support for very large ZIP files.

---

## 🔍 Health Endpoints

| Endpoint   | Purpose              |
| ---------- | -------------------- |
| `/healthz` | liveness probe       |
| `/readyz`  | dependency readiness |

---

## 🚀 Future Improvements

* [ ] Direct-to-object-storage uploads
* [ ] Prometheus metrics
* [ ] Worker autoscaling
* [ ] Kubernetes deployment
* [ ] File-level fan-out

---

## 🎯 Design Goals

DocNebula was built to demonstrate:

* event-driven architecture
* resilient background processing
* idempotent API design
* distributed worker coordination
* production-style failure handling

---

## 👨‍💻 Author

Built as a distributed systems practice project in Go.

---

## ⭐ If This Helped You

Consider giving the repo a star ⭐
