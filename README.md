# Go Financial Core

Go Financial Core is an event-driven financial backend built as a senior-level engineering portfolio project and as a practical laboratory for learning production-grade software engineering with Go.

The goal is not to be a small CRUD sample. The project is designed to look and evolve like a real company system: multiple applications, explicit domain boundaries, asynchronous communication, transactional persistence, observability, tests, generated dependencies, and infrastructure that can be run locally.

## Purpose

This repository exists to practice and demonstrate the engineering skills required to build reliable backend systems:

- Idiomatic Go beyond syntax: package design, interfaces, dependency injection, tests, error handling, and simple composition.
- Event-driven architecture with Kafka for decoupled workflows between services.
- Transactional consistency and relational modeling with PostgreSQL.
- Distributed systems fundamentals such as service boundaries, eventual consistency, async processing, retries, and operational visibility.
- Observability through structured logs, metrics, and distributed tracing.
- Clean Architecture with lightweight DDD: domain entities and policies at the center, use cases around them, and infrastructure kept behind ports.

## Domain

The system models a financial core around accounts and wallets.

Current capabilities include:

- Account creation, lookup, and deletion by phone.
- Wallet credit and debit operations.
- Domain events for core financial activity.
- Audit log recording through an asynchronous worker.
- Notification processing through an asynchronous worker.

The domain is intentionally narrow so the architecture can be deep: the focus is on correctness, boundaries, observability, and production-style design rather than on adding many unrelated features.

## Architecture

The repository is organized as a small distributed system:

- `apps/core`: HTTP API for account and wallet operations.
- `apps/audit`: Kafka worker responsible for recording audit events.
- `apps/notification`: Kafka worker responsible for notification workflows.
- `pkg`: shared infrastructure primitives such as API server helpers, logging, errors, metrics, validation, Kafka utilities, and base use case contracts.
- `docker`: local database bootstrap scripts.

Each application follows the same general layering:

- `domain`: entities, domain errors, events, and business policies.
- `application`: use cases, DTOs, mappers, and ports.
- `infrastructure`: HTTP handlers, database repositories, Kafka producers/consumers, migrations, and external clients.
- `di`: Wire-based dependency injection.
- `cmd`: executable entrypoints.

This keeps business rules independent from frameworks and external systems while still allowing the project to use real infrastructure locally.

## Tech Stack

- Go
- Gin for HTTP APIs
- PostgreSQL for transactional persistence
- Kafka for event-driven messaging
- OpenTelemetry and Jaeger for tracing
- Logrus for structured logging
- Swagger for API documentation
- Google Wire for dependency injection
- GoMock and Testify for tests
- Docker Compose for local infrastructure

## Local Development

### Requirements

- Go 1.25+
- Docker and Docker Compose
- Make
- Wire, when regenerating dependency injection code

### Run the full stack

```bash
docker compose up --build
```

The main API is exposed on:

```text
http://localhost:8080
```

Swagger documentation is available at:

```text
http://localhost:8080/docs/financial
```

Jaeger UI is available at:

```text
http://localhost:16686
```

### Run the core API locally

```bash
make run-core
```

### Regenerate dependency injection

```bash
make wire-core
make wire-audit
make wire-notifications
```

### Regenerate mocks

Linux/macOS:

```bash
make mocks-linux
```

Windows PowerShell:

```powershell
make mocks-windows
```

### Run tests

```bash
go test ./...
```

## API Surface

The core API is mounted under:

```text
/v1/financial
```

Available groups:

- `GET /health`
- `POST /accounts`
- `GET /accounts`
- `DELETE /accounts`
- `POST /wallets/credits`
- `POST /wallets/debits`

For request and response details, use the generated Swagger documentation.

## Engineering Direction

This project is meant to keep evolving toward a more realistic production system. Important areas of study and improvement include:

- Stronger transactional guarantees around wallet operations.
- Outbox/inbox patterns for reliable event publishing and consumption.
- Idempotency keys for financial commands.
- Better retry, dead-letter, and poison-message handling in workers.
- Richer metrics and service-level indicators.
- Integration tests with real PostgreSQL and Kafka.
- Security hardening around authentication, authorization, and secrets.
- CI pipelines, linting, and automated quality gates.

## What This Project Demonstrates

Go Financial Core is a deliberate practice project for backend engineering depth. It shows how to structure a Go codebase around real architectural concerns instead of framework-first code, and how financial workflows can be modeled with a mix of synchronous APIs, transactional data, asynchronous events, and operational observability.
