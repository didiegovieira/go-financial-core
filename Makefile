.PHONY: help up wire wire-core wire-notifications wire-audit docs-core mocks-linux mocks-windows run-core

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  up      - Build and start all services with Docker Compose"
	@echo "  wire    - Generate dependency injection code for all apps"
	@echo "  mocks-linux    - Generate mocks on Linux/macOS"
	@echo "  mocks-windows  - Generate mocks on Windows PowerShell"

up:
	docker compose up --build -d

wire: wire-core wire-audit wire-notifications

wire-core:
	wire gen ./apps/core/di

wire-notifications:
	wire gen ./apps/notification/di

wire-audit:
	wire gen ./apps/audit/di

docs-core:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init \
		-g main.go \
		-d apps/core/cmd/api \
		-o apps/core/docs \
		--parseInternal

mocks-linux:
	bash scripts/mocks-linux.sh

mocks-windows:
	powershell -ExecutionPolicy Bypass -File scripts/mocks-windows.ps1

run-core:
	go run ./apps/core/cmd/api
