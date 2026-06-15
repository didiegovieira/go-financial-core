package main

import (
	"context"
	"log"

	"github.com/didiegovieira/go-payments-core/apps/payments/internal/bootstrap"
	"github.com/didiegovieira/go-payments-core/apps/payments/internal/config"
)

func main() {
	settings, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load settings: %v", err)
	}

	app := bootstrap.InitializeApp(settings)
	if err := app.Run(context.Background()); err != nil {
		log.Fatalf("payments api stopped with error: %v", err)
	}
}
