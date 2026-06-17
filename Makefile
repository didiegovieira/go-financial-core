help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  wire    - Generate code using Wire for dependency injection"

wire:
	wire gen .\apps\payments\internal\bootstrap
