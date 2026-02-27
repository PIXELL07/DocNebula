🚀 DocNebula — Resilient Asynchronous Document Processing Pipeline

DocFlow is a production-grade, fault-tolerant document intelligence pipeline designed to process massive ZIP uploads reliably at scale. The system automatically extracts, OCRs, summarizes, and vectorizes documents using an event-driven architecture built for resilience and horizontal scalability.

Unlike traditional synchronous processors, DocFlow ensures that long-running jobs can recover gracefully from worker crashes using checkpointed progress tracking and idempotent task execution.

✨ Key Features

🔄 Fully asynchronous, queue-driven pipeline
🛡️ Crash-safe recovery with per-page checkpointing
♻️ Idempotent worker design
📦 Automatic ZIP extraction
🔍 OCR for scanned documents
🧠 AI summarization pipeline
🔢 Vector embedding generation
🚨 Dead Letter Queue (DLQ) for corrupted files
💓 Worker heartbeat monitoring
📈 Horizontally scalable architecture
🐳 Docker-first local development


🏗️ High-Level Architecture
User Upload
   ↓
Object Storage (S3/MinIO)
   ↓
Orchestrator (Go)
   ↓
Queue A → Unzip Workers
Queue B → OCR Workers
Queue C → Vector Workers
   ↓
PostgreSQL (state machine + metadata)