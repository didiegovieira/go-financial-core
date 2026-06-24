package main

import (
	"context"
	"log"

	"github.com/didiegovieira/go-payments-core/apps/payments/di"
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/settings"
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

	api, cleanup, err := di.InitializeApi()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup()

	api.Start()
}
