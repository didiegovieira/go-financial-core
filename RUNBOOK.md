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

```sh
docker compose config --quiet
make up
docker compose ps -a
docker compose logs --tail=100 financial-core audit notifications
```

`make up` runs `docker compose up --build -d`. It starts PostgreSQL, Kafka,
Jaeger, the core API, and both workers. A container being started does not mean
the application or Kafka is ready; inspect logs and the API response.

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
For worker health, inspect logs and verify the expected outcome of an event;
there is currently no Compose health check for either worker.

## Configuration

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
docker compose stop
```

Use `up --build` after code changes. `restart` does not rebuild an image or apply
changed environment configuration; use `docker compose up -d <service>` to
recreate a service after configuration changes.

`docker compose stop` retains containers for later use with `docker compose start`.
`docker compose down` removes containers and networks. The current Compose file
does not declare named data volumes for PostgreSQL or Kafka. Do not assume data
will survive container removal/recreation; PostgreSQL image-created anonymous
volumes are not automatically reattached by a later `up`. Use explicit persistent
volumes and backups if local data must be retained.

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
runbook intentionally has no automatic destructive reset command.

## Troubleshooting

| Symptom | Checks and action |
| --- | --- |
| Build or startup fails | Run `docker compose ps -a` and inspect the failing service's logs. Check Docker is running and the image/toolchain can be downloaded. |
| Port already allocated | Check host ports 8080, 5432, 9092, 16686, 4317, and 4318. Stop the conflicting process or change the published port. |
| App cannot connect to PostgreSQL | Check `pg_isready`, app database name, and credentials in both Compose and `env_file`. Inspect PostgreSQL logs before changing data. |
| Database does not exist | List databases and inspect initialization logs. Bootstrap only runs on an empty data directory; do not rerun the destructive script against existing data. |
| Worker fails during startup | Check Kafka logs. Compose waits for Kafka to start, not become ready. Once Kafka is ready, restart the affected app and inspect logs again. |
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
