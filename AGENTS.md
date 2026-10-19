# Repository guidance

## Project map

This is a Go financial backend with three applications in one module:

- `apps/core`: HTTP API for accounts and wallets; entrypoint `cmd/api`.
- `apps/audit`: Kafka audit worker; entrypoint `cmd/worker`.
- `apps/notification`: Kafka notification worker; entrypoint `cmd/worker`.
- `pkg`: shared infrastructure and utilities.
- `docker/postgres/init.sql`: local database bootstrap.
- `docs/RUNBOOK.md`: local operation and troubleshooting procedures.

Read the relevant application and its existing tests before changing behavior.
Use `go.mod`, `Makefile`, and `docker-compose.yml` as the source of truth for
versions, commands, and runtime configuration.

## Architecture and conventions

- Keep business rules in `internal/domain` and orchestration in
  `internal/application`. Define external dependencies through application ports.
- Keep HTTP, PostgreSQL, Kafka, and external clients in infrastructure. Do not
  introduce infrastructure dependencies into domain code.
- Compose dependencies in each app's `di` package with Wire.
- Respect Go `internal` package boundaries. Apps communicate through explicit
  contracts and events, rather than importing another app's internals.
- Put code in `pkg` only when it is genuinely shared; keep app behavior local.
- Follow nearby naming, error handling, logging, and test patterns. Format changed
  Go files with `gofmt`.
- For financial changes, check amount validation, balance rules, transaction
  boundaries, and failure paths. Do not silently change event payloads or
  consumer group semantics.

## Generated code

- Edit Wire providers/injectors, then run `make wire` (all apps), or the relevant
  `wire-core`, `wire-audit`, or `wire-notifications` target. Do not manually edit
  `wire_gen.go`.
- Regenerate mocks with `make mocks-windows` on PowerShell or `make mocks-linux`
  on Linux/macOS after changing mocked ports. Do not manually edit generated mocks.
- Regenerate core API documentation with `make docs-core` when changing Swagger
  annotations. Keep generated changes scoped to the task.

## Validation

- Run the affected package tests while developing; add meaningful regression tests
  for behavior changes.
- For Go code changes, run `go test ./...` and `go vet ./...` before handing off
  when the environment permits. Report failures and checks that could not run.
- For Compose changes, run `docker compose config --quiet`; for Makefile changes,
  inspect the affected targets with `make -n` before executing them.
- Documentation-only changes do not require running the Go suite. Verify commands
  and paths against the repository instead.

## Local environment and data

- Local workflow: `make up`, verify app initialization/migrations in the logs,
  then `make seed` for demo data. After using the application, run `make down`
  to stop services while preserving containers and data. PostgreSQL and Kafka
  use named volumes that survive container recreation. `make drop` removes
  containers, networks, and volumes, deleting local data; run it only when a
  destructive reset is explicitly requested.
- `make up` builds and starts Compose services without waiting for app
  initialization/migrations. `make seed` optionally applies local SQL data
  directly through PostgreSQL; it does not wait for startup or publish events.
- Compose currently reads `.env.example` directly. Per-app `environment` entries
  supply the database name and override shared settings.
- A root `.env` is loaded by apps running on the host; it is ignored by Git and
  excluded from Docker builds. Never commit credentials or log secrets.
- Preserve unrelated working-tree changes. Keep changes focused on the request.
- Database bootstrap includes `DROP DATABASE` statements. Do not rerun it against
  existing data or remove volumes without an explicit request to reset data.
- Treat missing reliability features (outbox, idempotency, retries, authentication)
  as work to implement, not as guarantees the system already provides.

## Handoff

Summarize what changed, validation performed, and any concrete limitations. Update
the runbook when operational commands, configuration, or recovery steps change.
