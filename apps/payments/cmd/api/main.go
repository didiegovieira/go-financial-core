package main

import (
	"log"

	"github.com/didiegovieira/go-payments-core/apps/payments/di"
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/settings"
)

func main() {
	settings.Load()

	api, cleanup, err := di.InitializeApi()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup()

	api.Start()
}
