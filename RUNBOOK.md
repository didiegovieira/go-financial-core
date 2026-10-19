# Local development runbook

Run commands from the repository root. This guide describes the current local
Compose stack, not a production deployment.

## Requirements

- Docker Engine/Desktop running, with Docker Compose v2 (`docker compose`).
- Make for the convenience targets; equivalent Docker commands work without it.
- Go matching `go.mod` (currently 1.25.6) for host builds and tests.
- Wire on `PATH` for dependency generation:

```sh
go install github.com/google/wire/cmd/wire@v0.7.0
```

## Start and verify

The local workflow is:

1. Run `make up` to build and start the stack.
2. Check app logs for successful initialization and completed migrations.
3. Run `make seed` to load optional demo data, then use the application.
4. Run `make down` when finished to stop services, preserving containers and data.

```sh
docker compose config --quiet
make up
docker compose ps -a
docker compose logs --tail=100 financial-core audit notifications
```

`make up` runs `docker compose up --build -d`. It confirms container startup,
not successful app initialization or completed migrations. Check core and audit
logs for initialization/migration errors before running `make seed`. The original
PostgreSQL health check only checks database readiness. Kafka dependencies wait
for container startup, not broker readiness.

| Component | Host access | Database |
| --- | --- | --- |
| Core API | `http://localhost:8080` | `financial_core_db` |
| Audit worker | No published HTTP port | `audit_db` |
| Notification worker | No published HTTP port | `notification_db` configured |
| PostgreSQL | `localhost:5432` | Three databases, one instance |
| Kafka | Port 9092 published; advertises `kafka:9092` | — |
| Grafana | `http://localhost:3000` (anonymous Editor) | Logs, metrics, traces |
| Tempo | `http://localhost:3200/ready` | Local trace storage |
| Loki | `http://localhost:3100/ready` | Local log storage |
| Prometheus | `http://localhost:9090` | Local metric storage |
| Alloy | `http://localhost:12345` | Collector |

Check the API with PowerShell:

```powershell
Invoke-RestMethod http://localhost:8080/v1/financial/health
```

Or with curl:

```sh
curl --fail http://localhost:8080/v1/financial/health
```

Expect HTTP 200 and an `ok: true` payload. This endpoint is a liveness check;
it does not verify database connectivity, Kafka, or worker processing.
Swagger UI is served at `http://localhost:8080/docs/financial/index.html`.
Workers have no HTTP health endpoints. Inspect logs and verify the expected
outcome of an event to assess worker operation.

## Configuration

### Local demo seeds

Seeds are optional and never executed by `make up`. After checking app logs and
confirming initialization and migrations completed, run:

```sh
make seed
```

This executes `psql` directly inside PostgreSQL with the read-only SQL directory
mounted at `/seeds`. There is no seed container or shell runner. Each database
uses `ON_ERROR_STOP` and its own transaction. A failed statement makes Make fail;
an audit failure does not roll back a completed core transaction. Running this
before migrations complete fails immediately rather than waiting.

The demo provides Alice (`5511999990001`) with 1000.00 and Bruno
(`5511999990002`) with 250.00 in active wallets. All identities and documents are
fictional. Audit includes one explicitly synthetic `account.created` record.
Notifications currently has no persistence migrations, so its database is left
unchanged. SQL inserts do not publish Kafka events.

Seeds use fixed UUIDs and skip conflicts; reapplying does not reset changed
balances or existing accounts. Conflicting emails/documents may cause a demo
account to be skipped. This is sample data, not a reset or synchronization tool.
Seed files remain local Docker infrastructure; no seed logic is added to apps.

All three apps use `env_file: .env.example`. Host, credentials, Kafka brokers,
and exporter settings are shared; each service's `environment` block sets its
database name, HTTP port, and tracing identity. Those entries take precedence
over `env_file` values. PostgreSQL credentials are also set separately in the
`postgres` service; keep them aligned when changing local credentials.

Copying `.env.example` to `.env` alone does not change the Compose `env_file`.
Apps run directly on the host load the root `.env` if present.
Do not put real secrets in the tracked example file.

Inside Compose, use service names (`postgres`, `kafka`, `alloy`) as hosts.
Host execution needs reachable host addresses. PostgreSQL and Alloy can use
their published localhost ports, but Kafka currently advertises `kafka:9092`.
Setting only `KAFKA_BROKERS=localhost:9092` is insufficient for host clients that
cannot resolve/reach that advertised address. Configure an external Kafka
listener before relying on host-run apps. `make run-core` also needs port 8080
free, so stop the Compose API first if it is running.

## Logs, rebuilds, and shutdown

```sh
docker compose logs -f --tail=100 financial-core audit notifications
docker compose logs --tail=100 postgres kafka tempo grafana alloy loki prometheus
docker compose up --build -d financial-core
docker compose restart audit notifications
make down
```

Use `up --build` after code changes. `restart` does not rebuild an image or apply
changed environment configuration; use `docker compose up -d <service>` to
recreate a service after configuration changes.

`make down` runs `docker compose stop`, retaining containers, networks, and data.
Use `make up` to start them again. PostgreSQL and Kafka use the named volumes
`postgres-data` and `kafka-data` (prefixed by the Compose project name), which
also preserve data when containers are recreated. Tempo and Grafana use `tempo-data` and `grafana-data`; traces have a 24-hour retention policy.

`make drop` runs `docker compose down --volumes`: it removes containers, networks,
and volumes, deleting local PostgreSQL, Kafka, Tempo, and Grafana data. Use it only for an
intentional reset, after backing up anything needed. It does not delete the SQL
files mounted from this repository. On the next `make up`, PostgreSQL runs its
bootstrap against the empty volume; seeds remain a separate `make seed` step.

When first adopting named volumes, existing data in old anonymous volumes or
container filesystems is not migrated automatically. Back up and restore it if
needed before recreating those containers.

## Database inspection and bootstrap

```sh
docker compose exec postgres pg_isready -U postgres
docker compose exec postgres psql -U postgres -c "SELECT datname FROM pg_database;"
docker compose exec postgres psql -U postgres -d financial_core_db -c "\dt"
docker compose exec postgres psql -U postgres -d audit_db -c "\dt"
```

`docker/postgres/init.sql` creates the three databases on initial PostgreSQL
initialization. It is not rerun when restarting an initialized instance.
Core and audit run their embedded SQL migrations during database connection
initialization. Their runner executes the embedded SQL on each startup; it does
not maintain a migration-version ledger. Keep scripts safe to rerun.

Do not manually rerun `init.sql` to fix a missing table: it drops databases.
Before any deliberate reset, identify the affected data and back it up. This
runbook's `make drop` command is a deliberate destructive reset, not a repair step.

## Troubleshooting

| Symptom | Checks and action |
| --- | --- |
| Build or startup fails | Run `docker compose ps -a` and inspect the failing service's logs. Check Docker is running and the image/toolchain can be downloaded. |
| Port already allocated | Check host ports 8080, 5432, 9092, 3000, 3100, 3200, 9090, 12345, 4317, and 4318. Stop the conflicting process or change the published port. |
| App cannot connect to PostgreSQL | Check `pg_isready`, app database name, and credentials in both Compose and `env_file`. Inspect PostgreSQL logs before changing data. |
| Database does not exist | List databases and inspect initialization logs. Bootstrap only runs on an empty data directory; do not rerun the destructive script against existing data. |
| Worker fails during startup | Check Kafka and worker logs. Compose waits for Kafka container startup, not broker readiness. Once Kafka is ready, restart the affected app and inspect its logs again. |
| Events are not processed | Inspect producer and worker logs, topic names, and consumer group configuration. A successful API health check does not prove event delivery. |
| No traces in Grafana | Inspect Tempo/app logs and `http://localhost:3200/ready`. Use `http://alloy:4317` with OTLP gRPC inside Compose. Select Tempo in Explore and expand the time range. |
| Host-run app cannot reach Kafka | Check the advertised listener configuration described above. Prefer the full Compose stack until host access is configured. |

## Development checks and generation

```sh
make wire
go test ./...
go vet ./...
```

Use `make wire-core`, `make wire-audit`, or `make wire-notifications` to regenerate
one app. Use `make mocks-windows` on PowerShell or `make mocks-linux` on Linux/macOS
for mocks, and `make docs-core` for Swagger artifacts. Do not edit generated files
by hand. After generation, review the diff before rebuilding services.

## Grafana observability MVP

The local observability MVP uses Alloy, Tempo, Loki, Prometheus, and Grafana
with provisioned data sources and three dashboards.
All three apps initialize the official OpenTelemetry Go SDK through
`pkg/telemetry` after loading settings. No vendor distribution or token is needed.
The shared endpoint in `.env.example` is now `http://alloy:4317`; per-app service
names and the existing spans/Gin middleware remain unchanged. Apps running on the
host must instead use `http://localhost:4317`. Recreate all three app containers
when applying the environment change; `restart` alone is insufficient.

Kafka now carries W3C trace context and baggage in headers. Each worker extracts
the message context before invoking its handler. Event payloads, consumer groups,
and commit/error behavior are unchanged. Old events without headers are still
accepted but cannot be linked back to an originating request. Avoid putting
sensitive information in baggage.

The shared Gin server also enables request-context fallback: existing handlers
pass `*gin.Context` to use cases, so forwarding its context methods to the HTTP
request is necessary to preserve the trace and request cancellation.

To try it:

1. Run `make up`, then inspect the three app logs for completed initialization.
2. Open `http://localhost:3000/explore` (anonymous Editor access, no login; permits Explore).
3. Request `http://localhost:8080/v1/financial/health` to generate an API trace.
4. Select Tempo, use Search or TraceQL `{ resource.service.name = "financial-core-api" }`,
   and select the last 15 minutes. Allow a few seconds for exporter batching.
5. Create an account using Swagger to publish real events. In that request's
   trace, inspect spans for `financial-core-api`, `audit-worker`, and
   `notification-worker`. Database seeds do not publish events. Health checks
   also do not exercise workers. Worker spans require successful event processing.
6. Run `make down` after use. It stops services and preserves volumes.

On a fresh broker, wait for Kafka and consumer groups to initialize before
creating an account. If publication reports `Unknown Topic Or Partition`, inspect
the topics and app logs; account persistence may already have succeeded, so do
not blindly repeat the same operation. If a worker does not resume consumption after Kafka becomes ready, inspect its
logs and consumer group, then restart the affected worker (for example,
`docker compose restart notifications`).
The MVP does not include broker-readiness checks, retries, or transactional event
publication.

To find a specific request without searching by service, copy its `X-Trace-ID`
response header into the Tempo Trace ID query in Grafana Explore. For example:

```powershell
$response = Invoke-WebRequest -UseBasicParsing http://localhost:8080/v1/financial/health
$response.Headers['X-Trace-ID']
```

If an old Jaeger container exists after the migration, inspect it with
`docker ps -a --filter ancestor=jaegertracing/all-in-one:1.76.0` and stop that
specific container with `docker stop <container-name>`; no volume reset is needed.
Old Jaeger traces are not migrated.

The MVP includes logs, metrics, and traces. Open
`http://localhost:3000/d/financial-overview` for the provisioned dashboard.
The failures dashboard is at `http://localhost:3000/d/financial-failures`.
Grafana sources are Tempo, Loki, and Prometheus; Explore also works for each.

### Collection and correlation

- Traces: all apps export OTLP to Alloy, which batches and forwards to Tempo.
- Metrics: `OTEL_METRICS_EXPORTER=otlp`, a 5-second export interval, and cumulative
  temporality send application metrics to Alloy. It copies only `service.name`
  into the `service_name` metric label, converts to Prometheus format, and remote
  writes to Prometheus. Runtime metrics from the official OpenTelemetry runtime instrumentation also appear.
- Logs: Alloy discovers only this Compose project's three application services
  and tails their stdout/stderr through Docker's API, forwarding them to Loki.
  `OTEL_LOGS_EXPORTER=none` is intentional: log collection does not use OTLP.
  Existing textual logs remain available alongside new JSON operational records.
- JSON completion records carry `trace_id`/`span_id`. Expand a JSON record in
  Loki and click **View trace**. Tempo's related logs query filters the three
  app streams by the trace ID. Historical textual logs do not gain trace IDs.

The dashboard includes HTTP totals, 4xx/5xx outcomes, throughput, p95 latency,
Kafka handler successes/errors, handler duration, per-topic totals, and separate
log panels for the three apps. Counters reset on app restart. Rates and p95 need
at least two exports and traffic; empty panels on a new instance are expected.
Metrics do not measure Kafka lag or successful offset commits: handler errors
are counted before the existing error path stops consumption. Commit failures
are outside the handler metric. HTTP 4xx count client/business failures, not 5xx.
No customer IDs, phones, arbitrary URLs, payloads, or trace IDs are metric labels.
New completion logs exclude request bodies, query strings, keys, and raw errors.
Notification delivery remains simulated by the existing noop client.

Useful queries:

```promql
sum by (route, status) (financial_http_requests_total)
sum by (service_name, topic, result) (financial_kafka_messages_total)
```

```logql
{service=~"financial-core|audit|notifications"}
{service="audit"} | json | trace_id="<X-Trace-ID>"
```

### Operations and limits

Validate configuration and startup with:

```sh
docker compose config --quiet
docker compose logs --tail=100 alloy loki prometheus
docker compose run --rm --no-deps alloy validate /etc/alloy/config.alloy
```

Prometheus: `http://localhost:9090`; Loki readiness: `http://localhost:3100/ready`;
Alloy: `http://localhost:12345`. All observability host ports bind to loopback.
After provisioning changes recreate Grafana (`docker compose up -d grafana`)
or restart it if only the contents of an existing bind mount changed. Rebuild
and recreate all apps for instrumentation/environment changes (`make up`).
Allow about a minute for collector discovery, backend readiness, and exports.

Named volumes also include `loki-data`, `prometheus-data`, and `alloy-data`
(collector positions and remote-write WAL). Logs/metrics retain 7 days; traces
retain 24 hours. `make down` preserves them. `make drop` deletes all observability
volumes as well as application data; do not use it for configuration updates.

Docker Desktop must expose its Linux daemon socket to the Alloy container.
The read-only socket mount still gives access to Docker's API: it is a local
development privilege, not a production isolation boundary. Alloy filters the project
via `COMPOSE_PROJECT_NAME` (default `go-financial-core`); when using a custom
project name set that environment variable rather than only `docker compose -p`.
Check discovery and Alloy logs if application streams are missing from Loki.

Anonymous Grafana access, unauthenticated backends, local filesystem storage,
and root containers for local volume initialization are development settings.
The MVP does not include service graphs, alert rules, outbox, retries, idempotency,
or production availability. No volumes need to be removed for this upgrade.

### Official OpenTelemetry SDK

`pkg/telemetry.Start` builds trace and metric providers with a shared resource,
OTLP exporters, and W3C tracecontext/baggage propagation. All entrypoints use
10-second startup/shutdown contexts. The batch trace processor and periodic
metric reader flush pending data on shutdown when the app reaches its cleanup.
Forced termination or existing fatal exits can still bypass deferred cleanup.

Supported exporter settings are `otlp` or `none` for traces and metrics. Transport
is `grpc` by default; `http/protobuf` is also supported, including signal-specific
`OTEL_EXPORTER_OTLP_TRACES_PROTOCOL` and `OTEL_EXPORTER_OTLP_METRICS_PROTOCOL`.
Official exporters read endpoint, headers, timeout, and TLS environment options.
For HTTP, use `http://alloy:4318` as the shared endpoint; the exporters append
signal paths. For host execution use the corresponding published localhost port.
`OTEL_SERVICE_NAME` and `OTEL_RESOURCE_ATTRIBUTES` set resource identity.
The trace SDK honors explicit `OTEL_TRACES_SAMPLER`/`OTEL_TRACES_SAMPLER_ARG`;
without a sampler setting, the previous always-sample behavior is preserved.
Metric export interval/temporality settings remain in `.env.example`.

`OTEL_PROPAGATORS` supports `tracecontext`, `baggage`, or `none`; legacy vendor
propagators are not loaded. `OTEL_LOGS_EXPORTER` must remain `none` because Alloy
collects container logs. Unsupported exporter/protocol/propagator settings fail
startup with an explicit configuration error. The SDK does not initialize a log
provider; existing Logrus and JSON operational records are unchanged.

### Telemetry package layout

Shared instrumentation lives in `pkg/telemetry`:

- `sdk.go`: official SDK providers, exporters, environment configuration, shutdown.
- `tracing.go`: spans, errors, trace/span IDs, tracer creation.
- `metrics.go`: HTTP/Kafka counters and duration histograms.
- `logs.go`: JSON completion records and trace correlation.
- `http.go`: Gin request tracing and HTTP observations.
- `propagation.go`: context injection/extraction through generic carriers.

Apps expose settings as `Settings.Telemetry` and use this shared package. Kafka header adaptation stays with the Kafka transport; generic
OpenTelemetry propagation is delegated to `pkg/telemetry`. Keep metric names/labels, span names, log fields, OTEL environment names, and
dashboard queries aligned when extending instrumentation. Wire providers obtain tracers from `telemetry.NewTracer`.

### Dashboard filters

Open `http://localhost:3000/d/financial-investigation` for a combined log view
and searchable traces. Use the time picker plus these filters:

- **Servico**: API, Audit, Notifications, multiple selections, or All. Applies
  to logs, traces, and metrics. A trace search matches a span in the chosen
  service; its root can still be the API.
- **HTTP: rota/metodo/status**: HTTP metrics and the API log panel only.
- **Kafka: topico/resultado**: Kafka metrics and the respective worker log panels.
- **Logs: nivel/contem texto/Trace ID**: log panels only. Text and trace ID are
  literal substring searches; blank fields disable those restrictions.
- **Traces: nome do span (regex)/duracao minima**: trace search only, on the
  investigation dashboard. These do not change the logs. To open an exact trace,
  filter its ID in logs and follow **View trace** in the record details.

HTTP and Kafka dropdown values are populated from observed Prometheus labels;
select All to clear a restriction. Explicit log levels or HTTP/Kafka fields can
exclude older textual logs without matching JSON fields. The failure dashboard
always restricts its logs to WARN/ERROR and its failure metrics to errors,
intersecting that restriction with the selected filters. Its combined failure
log view uses service/level/text/trace ID; HTTP/Kafka filters there affect metrics.
Text/trace ID filters never change aggregate metrics. Counter panels still show
process totals; the time picker controls log/trace search and metric history.

The general/failures dashboards link to investigation and preserve shared filter
values and the time range. Use All and clear text fields when a selection yields
no data. The service selector maps the app names to both Compose log service
labels and OpenTelemetry service identities. Refresh the dashboard after changes;
for new dashboard files, restart Grafana if provisioning has not picked them up.
