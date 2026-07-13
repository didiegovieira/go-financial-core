# Core API

## Overview

The `core` app is the system's transactional application. It exposes the HTTP API for accounts and wallets, persists financial state in PostgreSQL, and publishes domain events to Kafka for the other workers.

This document covers only the `core` context. The repository-wide view remains in the root `README.md`, while detailed HTTP contracts are available in the generated Swagger files in this directory.

## Responsibilities

- Create, retrieve, and delete accounts.
- Automatically create a zero-balance wallet for every new account.
- Credit and debit wallets.
- Apply the local risk policy before a debit.
- Persist accounts and wallets in `financial_core_db`.
- Publish account and wallet events to Kafka.

Audit storage and notification delivery do not belong to this app. Those effects are asynchronous and belong to the `audit` and `notification` apps, respectively.

## Application flow

```text
HTTP client
    |
    v
Gin handler -> use case -> domain rule -> PostgreSQL repository
                         |
                         +-> Kafka publisher -> audit / notification
```

Handlers convert HTTP input into DTOs. Use cases coordinate validation, normalization, repositories, and event publication. Concrete dependencies are assembled with Wire under `apps/core/di`.

## HTTP API

The API base path is `/v1/financial`.

| Method | Route | Purpose |
| --- | --- | --- |
| `GET` | `/health` | Check whether the HTTP process is responding. |
| `POST` | `/accounts` | Create an account and its initial wallet. |
| `GET` | `/accounts?phone=...` | Find an account by phone number. |
| `DELETE` | `/accounts?phone=...` | Delete the wallet and account associated with a phone number. |
| `POST` | `/wallets/credits` | Add funds to an account wallet. |
| `POST` | `/wallets/debits` | Withdraw funds after risk evaluation. |

Swagger is served at `http://localhost:8080/docs/financial`. The `docs.go`, `swagger.json`, and `swagger.yaml` files are generated; run `make docs-core` from the repository root to update them.

## Relevant rules

### Accounts

- Name, email, password, phone number, and document data are required for account creation.
- Phone numbers are normalized to digits only before storage and lookup.
- Email, phone number, and document number have uniqueness constraints in PostgreSQL.
- Account creation generates UUIDs for both the account and wallet; the wallet starts as `active` with a zero balance.

### Wallets and monetary values

- `amount` and `balance` are integer values expressed in cents. For example, USD 10.00 is represented as `1000`.
- Credits require an amount greater than zero and an active wallet.
- Debits require an amount greater than zero, an active wallet, and sufficient balance.
- The default policy approves amounts below `50000` cents, flags amounts from `50000` cents for review, and rejects amounts from `500000` cents. In the current implementation, both `review` and `rejected` stop the operation with HTTP 403; there is no manual review queue yet.

## Published events

All events use the account UUID as the Kafka key and carry a versioned envelope (`version: 1`) containing `event_id`, `event_name`, `occurred_at`, aggregate type and ID, and `payload`.

| Event | Topic | Published when |
| --- | --- | --- |
| `account.created` | `core.account.created` | After persisting the account and wallet. |
| `account.deleted` | `core.account.deleted` | After deleting the wallet and account. |
| `wallet.credited` | `core.wallet.credited` | After updating the balance. |
| `wallet.debited` | `core.wallet.debited` | After updating the balance. |

In the current design, PostgreSQL persistence and Kafka publication are not part of the same transaction. Publication may fail after a state change has been persisted; an outbox pattern remains a future improvement.

## Persistence

The migrations under `internal/infrastructure/database/postgres/migrations` maintain:

- `accounts`: account identity, contact, and document data.
- `wallets`: wallet ownership, balance in cents, and status.

The app contains connection implementations for other databases as extension material, but the active repositories and migrations currently use PostgreSQL.

## Configuration and execution

Main environment variables:

| Variable | Default | Purpose |
| --- | --- | --- |
| `HTTP_SERVER_PORT` | `:8080` | API port. |
| `POSTGRES_HOST` | `localhost` | PostgreSQL host outside Compose. |
| `POSTGRES_PORT` | `5432` | PostgreSQL port. |
| `POSTGRES_DATABASE` | `financial_core_db` | Transactional database. |
| `KAFKA_BROKERS` | `kafka:9092` | Broker list. |
| `OTEL_SERVICE_NAME` | module name | Service name used in telemetry. |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `http://localhost:4317` | OTLP destination. |

From the repository root:

```bash
make run-core
```

To run the app with PostgreSQL, Kafka, Jaeger, and the workers:

```bash
docker compose up --build
```

## Where to make changes

- `cmd/api`: executable entry point.
- `internal/domain`: entities, events, errors, and risk policy.
- `internal/application`: use cases, DTOs, mappers, and ports.
- `internal/infrastructure/api`: routes, handlers, and middleware.
- `internal/infrastructure/database/postgres`: migrations and repositories.
- `internal/infrastructure/messaging/kafka`: event mapping and publication.
- `di`: dependency composition with Wire.

