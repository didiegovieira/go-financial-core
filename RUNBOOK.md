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
| Jaeger | `http://localhost:16686` | — |

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

### Local demo seeds (POC)

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

Inside Compose, use service names (`postgres`, `kafka`, `jaeger`) as hosts.
Host execution needs reachable host addresses. PostgreSQL and Jaeger can use
their published localhost ports, but Kafka currently advertises `kafka:9092`.
Setting only `KAFKA_BROKERS=localhost:9092` is insufficient for host clients that
cannot resolve/reach that advertised address. Configure an external Kafka
listener before relying on host-run apps. `make run-core` also needs port 8080
free, so stop the Compose API first if it is running.

## Logs, rebuilds, and shutdown

```sh
docker compose logs -f --tail=100 financial-core audit notifications
docker compose logs --tail=100 postgres kafka jaeger
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
also preserve data when containers are recreated. Jaeger traces remain ephemeral.

`make drop` runs `docker compose down --volumes`: it removes containers, networks,
and volumes, deleting local PostgreSQL and Kafka data. Use it only for an
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
| Port already allocated | Check host ports 8080, 5432, 9092, 16686, 4317, and 4318. Stop the conflicting process or change the published port. |
| App cannot connect to PostgreSQL | Check `pg_isready`, app database name, and credentials in both Compose and `env_file`. Inspect PostgreSQL logs before changing data. |
| Database does not exist | List databases and inspect initialization logs. Bootstrap only runs on an empty data directory; do not rerun the destructive script against existing data. |
| Worker fails during startup | Check Kafka and worker logs. Compose waits for Kafka container startup, not broker readiness. Once Kafka is ready, restart the affected app and inspect its logs again. |
| Events are not processed | Inspect producer and worker logs, topic names, and consumer group configuration. A successful API health check does not prove event delivery. |
| No traces in Jaeger | Inspect Jaeger/app logs, exporter endpoint, protocol, and service name. Use `jaeger:4317` inside the Compose network. |
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
