package main

import (
	"log"

	"github.com/didiegovieira/go-payments-core/apps/audit/di"
	"github.com/didiegovieira/go-payments-core/apps/audit/internal/settings"
)

func main() {
	settings.Load()

	worker, cleanup, err := di.InitializeWorker()
	if err != nil {
		log.Fatalf("Failed to initialize worker: %v", err)
	}
	defer cleanup()

	worker.Start()
}
