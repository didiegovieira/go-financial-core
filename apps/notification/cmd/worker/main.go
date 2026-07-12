package main

import (
	"context"
	"log"

	"github.com/didiegovieira/go-financial-core/apps/notification/di"
	"github.com/didiegovieira/go-financial-core/apps/notification/internal/settings"
	"github.com/signalfx/splunk-otel-go/distro"
)

func main() {
	settings.Load()

	sdk, err := distro.Run()
	if err != nil {
		log.Fatalf("Failed to run distro: %v", err)
	}
	defer func() {
		if err := sdk.Shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shutdown SDK: %v", err)
		}
	}()

	worker, cleanup, err := di.InitializeWorker()
	if err != nil {
		log.Fatalf("Failed to initialize worker: %v", err)
	}
	defer cleanup()

	worker.Start()
}
