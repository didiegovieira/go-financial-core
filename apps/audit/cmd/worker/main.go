package main

import (
	"context"
	"log"
	"time"

	"github.com/didiegovieira/go-financial-core/apps/audit/di"
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/settings"
	"github.com/didiegovieira/go-financial-core/pkg/telemetry"
)

func main() {
	settings.Load()

	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 10*time.Second)
	sdk, err := telemetry.Start(startupCtx)
	cancelStartup()
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := sdk.Shutdown(shutdownCtx); err != nil {
			log.Printf("Failed to shutdown telemetry: %v", err)
		}
	}()

	worker, cleanup, err := di.InitializeWorker()
	if err != nil {
		log.Fatalf("Failed to initialize worker: %v", err)
	}
	defer cleanup()

	worker.Start()
}
