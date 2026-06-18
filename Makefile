help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  wire    - Generate code using Wire for dependency injection"

wire-payments:
	wire gen ./apps/payments/di

wire-notifications:
	wire gen ./apps/notifications/di

wire-fraud:
	wire gen ./apps/fraud/di

wire-audit:
	wire gen ./apps/audit/di
