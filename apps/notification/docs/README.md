# Notification Worker

## Overview

The `notification` app is the asynchronous worker responsible for reacting to `core` events and turning them into notification requests. It does not expose an HTTP API and, in its current composition, does not persist state.

The app defines a client port to separate the decision to notify from the provider used to deliver the notification.

## Responsibilities

- Consume account and wallet events from Kafka.
- Translate each topic into its corresponding event name.
- Build a `Notification` entity from the received key and payload.
- Delegate delivery to a `NotificationClient` implementation.
- Shut down the consumer gracefully when the process receives a system signal.

This app does not decide or alter the outcome of the original financial operation.

## Processing flow

```text
Kafka topic
    |
    v
Consumer -> handler -> SendNotification -> NotificationClient
                                             |
                                             +-> noop (current configuration)
                                             +-> HTTP (available implementation)
```

## Consumed topics

The default consumer group is `notification-service`.

| Topic | Handled event |
| --- | --- |
| `core.account.created` | `account.created` |
| `core.account.deleted` | `account.deleted` |
| `core.wallet.credited` | `wallet.credited` |
| `core.wallet.debited` | `wallet.debited` |

An unmapped topic returns an unknown-event error. New `core` events must be added to `internal/infrastructure/messaging/kafka/topics.go` before they can be consumed here.

## Notification model

For each message, the use case creates:

| Field | Current value |
| --- | --- |
| `id` | A new UUID string. |
| `event_name` | Event name derived from the topic. |
| `aggregate_id` | Kafka key, currently the account UUID. |
| `channel` | `external`. |
| `recipient` | The `aggregate_id` itself. |
| `subject` | The event name. |
| `body` | The complete JSON payload converted to text. |
| `status` | `pending`. |
| `created_at` | UTC processing time. |

This mapping is intentionally generic. The worker does not yet extract an email address or phone number from the envelope, select templates by event, or track status changes after delivery.

## Output clients

### Noop, currently active

Wire binds `NotificationClient` to `internal/infrastructure/client/noop`. This implementation writes the notification metadata to the log and returns success; it does not send a real email, SMS, or webhook.

### HTTP, available but not injected

An implementation under `internal/infrastructure/client/http` serializes the entity and sends a `POST` request to a supplied URL. It is not currently selected in `di/inject_client.go`, and the app configuration does not expose a dedicated variable for that URL.

## Configuration and execution

| Variable | Default | Purpose |
| --- | --- | --- |
| `KAFKA_BROKERS` | `kafka:9092` | Broker list. |
| `KAFKA_CONSUMER_GROUP_ID` | `notification-service` | Consumer group. |
| `OTEL_SERVICE_NAME` | module name | Worker name used in telemetry. |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `http://localhost:4317` | OTLP destination. |

Run directly when Kafka is available:

```bash
go run ./apps/notification/cmd/worker
```

Run the complete ecosystem:

```bash
docker compose up --build
```

The HTTP and database settings present in `settings` do not mean that this worker currently opens an API or uses a database. They are not part of the injected runtime flow.

## Current limitations

- The active client only writes logs.
- There are no templates, channel preferences, or real recipient resolution.
- Delivery attempts and status are not persisted.
- There is no explicit deduplication by `event_id`.
- Retries, a dead-letter queue, and poison-message handling are not modeled in this app.

## Where to make changes

- `cmd/worker`: executable entry point.
- `internal/application/usecase`: notification construction and delivery.
- `internal/application/port/client`: output channel contract.
- `internal/domain/entity`: the `Notification` model.
- `internal/infrastructure/messaging/kafka`: topics, handler, and consumer.
- `internal/infrastructure/client`: `noop` and HTTP adapters.
- `internal/runtime`: worker lifecycle and shutdown.
- `di`: dependency selection and composition.

