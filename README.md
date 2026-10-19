# Go Financial Core

[![CI](https://github.com/didiegovieira/go-financial-core/actions/workflows/ci.yml/badge.svg)](https://github.com/didiegovieira/go-financial-core/actions/workflows/ci.yml)

Go Financial Core is a financial backend MVP built with Go, PostgreSQL, and Kafka. It provides account and wallet operations, asynchronous audit and notification workflows, and integrated logs, metrics, and distributed tracing.

The MVP runs locally with Docker Compose and serves as a practical backend engineering project. Notification delivery is simulated. Production reliability and security features remain planned work; see the current limits below.

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
- `pkg`: shared infrastructure primitives such as API server helpers, logging, errors, telemetry, validation, Kafka utilities, and base use case contracts.
- `docker`: database bootstrap/seeds and observability configuration, data sources, and dashboards.

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
- OpenTelemetry and Alloy for collection; Grafana, Tempo, Loki, and Prometheus for traces, logs, and metrics
- Logrus for structured logging
- Swagger for API documentation
- Google Wire for dependency injection
- GoMock and Testify for tests
- Docker Compose for local infrastructure

## Local Development

See [RUNBOOK.md](RUNBOOK.md) for startup, diagnostics, configuration, and data
handling. Repository guidance for coding agents lives in [AGENTS.md](AGENTS.md).

### Requirements

- Go 1.25.6 (matching `go.mod`)
- Docker and Docker Compose
- Make
- Wire, when regenerating dependency injection code

### Run the full stack

```bash
make up
```

`make up` starts the containers without waiting for app initialization or migrations.
Check app logs before loading optional demo accounts, wallets, and an audit record:

```sh
make seed
```

See [local seed procedures](RUNBOOK.md#local-demo-seeds). SQL seeds do not
publish Kafka events.

The local workflow is `make up`, check initialization/migrations, `make seed`,
then use the application. When finished, shut down with:

```sh
make down
```

`make down` only stops services, preserving containers and data. PostgreSQL and
Kafka use named volumes, so subsequent `make up` runs retain their data.

For a complete local reset, `make drop` removes containers, networks, and volumes,
deleting database, Kafka, and observability data. The next `make up` initializes a fresh environment;
run `make seed` after initialization if you want demo data again.

The main API is exposed on:

```text
http://localhost:8080
```

Swagger documentation is available at:

```text
http://localhost:8080/docs/financial
```

### Observability MVP

The three apps use the official OpenTelemetry Go SDK through `pkg/telemetry`.
Traces and metrics go through Alloy to Tempo and Prometheus. Alloy collects
container stdout/stderr into Loki. Grafana correlates JSON operational logs with
traces, including context propagated through Kafka headers.

Grafana is available at `http://localhost:3000` with anonymous Editor access for
local investigation. Provisioned dashboards include:

| Dashboard | Local URL | Purpose |
| --- | --- | --- |
| Overview | [Open overview](http://localhost:3000/d/financial-overview) | HTTP/Kafka volume, outcomes, latency, and app logs |
| Failures | [Open failures](http://localhost:3000/d/financial-failures) | HTTP failures, worker handler errors, and warning/error logs |
| Investigation | [Open investigation](http://localhost:3000/d/financial-investigation) | Filter logs and search traces by service, span, and duration |

Use the filters at the top and the time picker to narrow results. Expand a JSON
log record and click **View trace**, or copy the API response's `X-Trace-ID` into
[Grafana Explore](http://localhost:3000/explore) with the Tempo source selected.
SQL seeds do not generate events or worker traces; use the API to exercise them.

See [observability operations](RUNBOOK.md#grafana-observability-mvp) and
[dashboard filters](RUNBOOK.md#dashboard-filters) for configuration, queries,
retention, and troubleshooting.

### Current MVP limits

- Local deployment: anonymous Grafana access, unauthenticated telemetry backends,
  filesystem storage, and single instances. This is not a production deployment.
- Notifications use the noop client; no external delivery is performed.
- No outbox/inbox, idempotency, application retry policy, dead-letter handling, or authentication.
  A database write can succeed before event publication fails.
- Kafka readiness is not guaranteed by container startup. Inspect app logs and
  consumer processing before relying on the asynchronous flows.
- Worker metrics measure handler completion, not Kafka lag or offset commits.
  Service graphs, alert rules, SLOs, and production availability remain future work.

### Run the core API locally

```bash
make run-core
```

### Regenerate dependency injection

```bash
make wire
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

### Shared telemetry code

`pkg/telemetry` centralizes OpenTelemetry SDK setup, tracing helpers, operational
metrics, correlated JSON logs, HTTP instrumentation, and context propagation.
Apps use `telemetry.TraceWithError`, `telemetry.TraceHandler`, and
`Settings.Telemetry` for instrumentation and configuration.
Transport-specific Kafka header adaptation remains in `pkg/messaging/kafka`.

## Automated study checks

Push directly to `main` or open an optional PR targeting `main`. The same CI runs
after direct pushes and when a PR is opened, reopened, or updated with new commits.
It can also be started manually in GitHub Actions. It checks formatting, vet, race detection,
builds, and Compose configuration. A failed run reports a problem to investigate;
these files do not require pull requests or configure branch protection.

The **Go vulnerabilities** workflow is manual: open Actions, select it, and choose
**Run workflow** when you want to inspect dependencies. Existing toolchain and
dependency findings remain pending remediation. There are no scheduled scans,
automated dependency PRs, or mandatory contribution templates.
