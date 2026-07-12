help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  wire    - Generate code using Wire for dependency injection"

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

run-core:
	go run ./apps/core/cmd/api
